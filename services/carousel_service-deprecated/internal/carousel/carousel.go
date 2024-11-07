package carousel

import (
	"fmt"
	"time"

	"carousel_service/internal/logger"
	// "carousel_service/internal/persistency"

	// "github.com/google/uuid"

	. "carousel_service/internal/ports"
	pbr "carousel_service/internal/ports/port_broker"
	pcl "carousel_service/internal/ports/port_carousel"
	pdb "carousel_service/internal/ports/port_persistency"
	. "carousel_service/internal/utils"

	"github.com/google/uuid"
)

const (
	heartbeatTimeout  = 2 * time.Minute
	subscriberTimeout = 30 * time.Second
)

type CommandUpdateSubscriber struct {
}
type CommandPlay struct {
	eventId    uuid.UUID
	carouselId string
}

func NewCarouselRunner(broker pbr.BrokerInterface, persistency pdb.PersistencyGatewayInterface) CarouselRunner {
	return CarouselRunner{
		subscribers:      pbr.Subscribers{Carousels: make(map[string]byte)},
		portCmd:          NewPort[any](100),
		broker:           broker,
		persistency:      persistency,
		pendingPlayReply: make(map[string]uuid.UUID),
	}
}

func NewCarouselHandler(client CarouselMasterInterface, persistency pdb.PersistencyGatewayInterface) CarouselHandler {
	return CarouselHandler{
		master:      client,
		persistency: persistency,
	}
	// portBrokerRunner: runner.portCmd,

}

type CarouselHandler struct {
	master      CarouselMasterInterface
	persistency pdb.PersistencyGatewayInterface
	// portBrokerRunner PortInterface[any]
}

func (c *CarouselHandler) Register(data pcl.RegisterData) error {
	s, err := c.persistency.Open()
	if err != nil {
		return err
	}
	defer s.Close()
	cr := pdb.CarouselRecord{
		CarouselId: data.Id,
		OwnerId:    data.OwnerId,
		RoundTime:  data.RoundTime,
	}
	err = s.Carousel().Create(cr)
	if err != nil {
		return err
	}
	sr := pdb.StatusRecord{
		CarouselId:  data.Id,
		Status:      NewOptionalValue[string](pdb.StatusNameoNew),
		RoundsReady: NewOptionalValue[int](0),
	}
	err = s.Status().Create(sr)
	if err != nil {
		return err
	}
	c.master.CommandTransmitter().Send(CommandUpdateSubscriber{})
	return nil
}

func (c *CarouselHandler) Delete(data pcl.CarouselId) error {
	s, err := c.persistency.Open()
	if err != nil {
		return err
	}
	defer s.Close()
	err = s.Carousel().Delete(data.Id)
	if err != nil {
		return err
	}
	err = s.Status().Delete(data.Id)
	if err != nil {
		return err
	}
	c.master.CommandTransmitter().Send(CommandUpdateSubscriber{})
	return nil
}

func (c *CarouselHandler) Read(data pcl.CarouselId) (Optional[pcl.AggregationData], error) {
	var pi pdb.PersistencyInterface
	var err error
	if pi, err = c.persistency.Open(); err != nil {
		return NewOptionalNil[pcl.AggregationData](), err
	}
	defer pi.Close()
	var cr Optional[pdb.CarouselRecord]
	if cr, err = pi.Carousel().Read(data.Id); err != nil {
		return NewOptionalNil[pcl.AggregationData](), err

	}
	if !cr.Valid() {
		return NewOptionalNil[pcl.AggregationData](), fmt.Errorf("Couldn't status carrosel by id:%s", data.Id)
	}
	var sr Optional[pdb.StatusRecord]
	if sr, err = pi.Status().Read(data.Id); err != nil {
		return NewOptionalNil[pcl.AggregationData](), err
	}
	if !sr.Valid() {
		return NewOptionalNil[pcl.AggregationData](), fmt.Errorf("Couldn't find status by id:%s", data.Id)

	}

	cid := pcl.CarouselId{Id: cr.Ptr().CarouselId}
	rd := pcl.RoundData{CarouselId: cid, RoundsReady: sr.Get().RoundsReady}
	ad := pcl.AggregationData{
		RoundData: rd,
		RoundTime: cr.Ptr().RoundTime,
		Status:    sr.Ptr().Status.Get(),
		Time:      sr.Ptr().Time,
	}
	// r, err := c.db.Read(data.Id)
	return NewOptionalValue[pcl.AggregationData](ad), nil
}

func (c *CarouselHandler) ReadByOwner(ownerId string) (Optional[[]pcl.AggregationData], error) {
	var pi pdb.PersistencyInterface
	var err error
	if pi, err = c.persistency.Open(); err != nil {
		return NewOptionalNil[[]pcl.AggregationData](), err
	}
	defer pi.Close()
	var crArray Optional[[]pdb.CarouselRecord]
	crArray, err = pi.Carousel().ReadManyBy(func() string {
		return fmt.Sprintf("OwnerId='%s'", ownerId)
	})
	if err != nil {
		return NewOptionalNil[[]pcl.AggregationData](), err

	}
	if !crArray.Valid() {
		return NewOptionalNil[[]pcl.AggregationData](), fmt.Errorf("Couldn't status carrosel by id:%s", ownerId)
	}
	var adArray []pcl.AggregationData
	for _, cr := range crArray.Get() {
		var sr Optional[pdb.StatusRecord]
		if sr, err = pi.Status().Read(cr.CarouselId); err != nil {
			return NewOptionalNil[[]pcl.AggregationData](), err
		}
		if !sr.Valid() {
			return NewOptionalNil[[]pcl.AggregationData](), fmt.Errorf("Couldn't find status by id:%s", ownerId)
		}

		cid := pcl.CarouselId{Id: cr.CarouselId}
		rd := pcl.RoundData{CarouselId: cid, RoundsReady: sr.Get().RoundsReady}
		ad := pcl.AggregationData{
			RoundData: rd,
			RoundTime: cr.RoundTime,
			Status:    sr.Ptr().Status.Get(),
			Time:      sr.Ptr().Time,
		}
		adArray = append(adArray, ad)
	}
	return NewOptionalValue[[]pcl.AggregationData](adArray), nil
}

func (c *CarouselHandler) Refill(data pcl.RefillData) error {
	var pi pdb.PersistencyInterface
	var err error
	if pi, err = c.persistency.Open(); err != nil {
		return err
	}
	defer pi.Close()

	if data.RoundsReady < 1 {
		return fmt.Errorf("Refill of id: %s has failed due to invalid data:%d", data.Id, data.RoundsReady)
	}

	var sr Optional[pdb.StatusRecord]
	if sr, err = pi.Status().Read(data.Id); err != nil {
		return err
	}
	if !sr.Valid() {
		return fmt.Errorf("Couldn't find status by id:%s", data.Id)
	}

	sr.Ptr().RoundsReady.Set(data.RoundsReady + sr.Ptr().RoundsReady.Get())
	if err = pi.Status().Update(sr.Get()); err != nil {
		return err
	}
	return nil
}
func (c *CarouselHandler) Play(data pcl.CarouselId) error {
	// Use event record table for persistent storage event and not loose any
	var pi pdb.PersistencyInterface
	var err error
	if pi, err = c.persistency.Open(); err != nil {
		return err
	}
	defer pi.Close()

	var sr Optional[pdb.StatusRecord]
	if sr, err = pi.Status().Read(data.Id); err != nil {
		return err
	}
	if !sr.Valid() {
		return fmt.Errorf("Couldn't find status by id:%s", data.Id)
	}
	if sr.Ptr().Status.Get() != pdb.StatusNameOnline {
		return fmt.Errorf("Couldn't operate when target status is '%s'", sr.Ptr().Status.Get())
	}
	if sr.Ptr().RoundsReady.Get() < 1 {
		return fmt.Errorf("Not allowed, rounds aren't enough '%d'", sr.Ptr().RoundsReady.Get())
	}
	play := CommandPlay{
		carouselId: data.Id,
		eventId:    uuid.New(),
	}
	c.master.CommandTransmitter().Send(play)
	return nil
}

type CarouselRunnerInterface interface {
	// Separate thread
	Run()
}

type CarouselMasterInterface interface {
	CommandTransmitter() PortInterface[any]
	// ResponseListener() PortInterface[any]
}

type CarouselRunner struct {
	subscribers      pbr.Subscribers
	pendingPlayReply map[string]uuid.UUID
	// Driven Ports
	broker      pbr.BrokerInterface
	persistency pdb.PersistencyGatewayInterface

	// Driver Ports
	portCmd PortInterface[any]
}

func (c *CarouselRunner) CommandTransmitter() PortInterface[any] {
	return c.portCmd
}

//	func (c *CarouselRunner) GetEventReceiverPort() PortReceiverInterface[any] {
//		return &c.portInternalEvt
//	}
func (c *CarouselRunner) updateSubscribers() error {
	// logger.Debug.Println("Enter to updateSubscribers")
	var err error
	var pi pdb.PersistencyInterface
	if pi, err = c.persistency.Open(); err != nil {
		return err
	}
	defer pi.Close()
	var srArray Optional[[]pdb.StatusRecord]
	if srArray, err = pi.Status().ReadAll(); err != nil {
		return err
	}
	var updated bool
	for _, sr := range srArray.Get() {
		_, ok := c.subscribers.Carousels[sr.CarouselId]
		if !ok {
			updated = true
			c.subscribers.Carousels[sr.CarouselId] = 1 //QOS
			logger.Debug.Printf("Add to a subscriber: %s", sr.CarouselId)
		}
		// c.subscribers.Carousels[sr.CarouselId] = 1
	}
	if updated {
		c.broker.SetSubscribers(c.subscribers)
	}
	return nil
}
func (c *CarouselRunner) timerHandlerHeartbeat() error {
	// logger.Debug.Println("Enter to timerHandlerHeartbeat")
	var err error
	var pi pdb.PersistencyInterface
	if pi, err = c.persistency.Open(); err != nil {
		return err
	}
	defer pi.Close()
	var srArray Optional[[]pdb.StatusRecord]
	if srArray, err = pi.Status().ReadAll(); err != nil {
		return err
	}

	for _, st := range srArray.Get() {
		t, err := time.Parse(time.RFC3339, st.Time)
		if err != nil {
			logger.Error.Printf("Timer heartbeat handler: %s", err)
			continue
		}
		logger.Debug.Printf("Carouseld:'%s', LastUpdate: '%s', Elapsed: '%s', Status:%s", st.CarouselId, t, time.Since(t), st.Status.Get())
		// if st.Status.Get() != pdb.StatusNameOffline && time.Since(t) > heartbeatTimeout {
		if st.Status.Get() == pdb.StatusNameOnline && time.Since(t) > heartbeatTimeout {
			st.Status.Set(pdb.StatusNameOffline)
			if err := pi.Status().Update(st); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *CarouselRunner) eventHandlerHeartbeat(event pbr.EventHeartbeat) error {
	var err error
	var pi pdb.PersistencyInterface
	if pi, err = c.persistency.Open(); err != nil {
		return err
	}
	defer pi.Close()
	var sr Optional[pdb.StatusRecord]
	if sr, err = pi.Status().Read(event.CarouselId); err != nil {
		return err
	}
	if !sr.Valid() {
		return fmt.Errorf("Can't find status record by CarouselId: %s", event.CarouselId)
	}
	sr.Ptr().Status.Set(pdb.StatusNameOnline)
	return pi.Status().Update(sr.Get())
}

func (c *CarouselRunner) commandHandler(cmd any) {
	switch ct := cmd.(type) {
	case nil:
		break
	case CommandUpdateSubscriber:
		if err := c.updateSubscribers(); err != nil {
			logger.Error.Printf("Fail to execute command - Update Subscribers, due to error: '%s'", err)
		}
	case CommandPlay:
		c.pendingPlayReply[ct.carouselId] = ct.eventId
		play := pbr.NewMessageCommand("Play")
		play.CarouselId = ct.carouselId
		play.EventId = ct.eventId.String()

		c.broker.PublishQueue(*play)
	default:
		logger.Error.Println("Type is unknown") // here v has type interface{}
	}

}
func (c *CarouselRunner) eventHandler(event any) {
	switch e := event.(type) {
	case nil:
		break
	case pbr.EventAck:
		var remoteError error
		if len(e.Error) != 0 {
			remoteError = fmt.Errorf("%s", e.Error)
			logger.Error.Printf("Got failure from CarouselId:%s caused by EventId:%s, with error:'%s'", e.CarouselId, e.CorrelationId, remoteError)
		}
		if v, ok := c.pendingPlayReply[e.CarouselId]; ok && (v.String() == e.CorrelationId) {
			delete(c.pendingPlayReply, e.CarouselId)
			if pi, err := c.persistency.Open(); err == nil {
				defer pi.Close()
				// TODO save the remoteError to the database
				if st, err := pi.Status().Read(e.CarouselId); (err == nil) && st.Valid() && remoteError == nil {
					if st.Ptr().RoundsReady.Get() > 0 {
						*st.Ptr().RoundsReady.Ptr() -= 1
						logger.Info.Printf("New RoundsReady:%d for CarouselId:%s", st.Ptr().RoundsReady.Get(), st.Ptr().CarouselId)
					}
					if err := pi.Status().Update(st.Get()); err != nil {
						logger.Error.Printf("Fail to update the Status record by id: '%s', err: %s", e.CarouselId, err)
					}
				} else {
					logger.Error.Printf("Fail to read the Status record by id: '%s, err: %s'", e.CarouselId, err)
				}
			} else {
				logger.Error.Printf("Fail to open persistency: '%s'", err)
			}
		}
	case pbr.EventHeartbeat:
		if err := c.eventHandlerHeartbeat(e); err != nil {
			logger.Error.Panicf("Fail to handle Heartbeat evet: %s", err)
		}
	default:
		logger.Error.Println("Type is unknown") // here v has type interface{}
	}
}

func (c *CarouselRunner) Run() {
	heartbeat := time.NewTicker(heartbeatTimeout)
	subscriber := time.NewTicker(subscriberTimeout)
	// Run subscribing logic for the first time
	if err := c.updateSubscribers(); err != nil {
		logger.Error.Printf("Fail to update Subscribers: %s", err)
	}
	for {
		select {
		case cmd := <-c.portCmd.Receiver():
			c.commandHandler(cmd)
			break
		case evt := <-c.broker.ListenQueue().Receiver():
			// logger.Debug.Printf("Receive: %+v\n", evt)
			c.eventHandler(evt)
			_ = evt
			break
		case <-subscriber.C:
			if err := c.updateSubscribers(); err != nil {
				logger.Error.Printf("Timer: update subscribers failed: %s", err)
			}
			break
		case <-heartbeat.C:
			if err := c.timerHandlerHeartbeat(); err != nil {
				logger.Error.Printf("Timer: heartbeat handler failed: %s", err)
			}
			break
		}
	}
}
