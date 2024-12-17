package operator

import (
	"carousel/infrastructure/broker/topic"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	tmCarouselOffline       = 3 * time.Minute
	tmMonitorOfflinePeriod  = 30 * time.Second
	tmMonitorPendingPeriod  = 10 * time.Second
	tmMonitorSnapshotPeriod = 1 * time.Hour
	maxRetries              = 3
)

type Operator struct {
	evRepo  IPortOperatorAdapterEventRepository
	crRepo  IPortOperatorAdapterCarouselRepository
	snRepo  IPortOperatorAdapterSnapshotRepository
	broker  IPortOperatorAdapterMqtt
	config  IPortOperatorAdapterConfig
	log     *zerolog.Logger
	retries map[string]int
}

func New(
	evRepo IPortOperatorAdapterEventRepository,
	crRepo IPortOperatorAdapterCarouselRepository,
	snRepo IPortOperatorAdapterSnapshotRepository,
	broker IPortOperatorAdapterMqtt,
	config IPortOperatorAdapterConfig,
	log *zerolog.Logger) *Operator {

	op := &Operator{evRepo: evRepo, crRepo: crRepo, snRepo: snRepo, broker: broker, config: config, log: log, retries: make(map[string]int)}
	return op
}

func (o *Operator) Refill(c Carousel, tickets int) error {
	var err error
	var exists bool

	for ok := true; ok; ok = false {
		if tickets < 1 {
			err = fmt.Errorf("Operator.Refill: Invalid Tickets value: %d", tickets)
			break
		}
		if exists, err = o.crRepo.OperatorIsExistsCarousel(c.CarId); err != nil {
			break
		}
		if !exists {
			err = fmt.Errorf("Operator.Refill: Doesn't exists")
			break
		}
		rd := TicketsData{CarId: c.CarId, Tickets: tickets, EvtId: uuid.New()}
		o.log.Info().Str("CarouselId", rd.CarId).Str("EventId", rd.EvtId.String()).Int("Tickets", rd.Tickets).Msg("Operator.Refill: About to write an event")
		err = o.evRepo.OperatorRefill(&rd)
	}
	return err
}

func (o *Operator) Play(c Carousel) error {
	var err error
	var exists bool

	for ok := true; ok; ok = false {
		if exists, err = o.crRepo.OperatorIsExistsCarousel(c.CarId); err != nil {
			break
		}
		if !exists {
			err = fmt.Errorf("Operator.Play: Doesn't exists")
			break
		}
		// TODO add read shapshot table
		var s *SnapshotData
		if s, err = o.readSnapshot(&c); err != nil {
			break
		}
		// TODO remove false
		if false && s.Status != CarouselStatusNameOnline {
			err = fmt.Errorf("Operator.Play: Carousel Status is '%s'", s.Status)
			break
		}
		if s.Tickets == 0 {
			err = fmt.Errorf("Operator.Play: Carousel has no Tickets (%d Tickets)", s.Tickets)
			break
		}

		pd := PlayData{CarId: c.CarId, EvtId: uuid.New()}
		o.log.Info().Str("CarouselId", pd.CarId).Str("EventId", pd.EvtId.String()).Msg("Operator.Play: About to write an event")
		if err = o.evRepo.OperatorPlay(&pd); err != nil {
			break
		}
		o.log.Info().Str("CarouselId", pd.CarId).Str("EventId", pd.EvtId.String()).Str("Type", MsgTypeRequestPlay).Msg("Operator.Play: About to publish the Play event")
		err = o.publish(&RequestPlay{MessageGeneric: MessageGeneric{MsgType: MsgTypeRequestPlay, CarId: pd.CarId}, EvtId: pd.EvtId.String()})
	}
	return err
}

func (o *Operator) Read(c Carousel) (*SnapshotData, error) {
	var err error
	var exists bool
	var sd *SnapshotData
	for ok := true; ok; ok = false {
		if exists, err = o.crRepo.OperatorIsExistsCarousel(c.CarId); err != nil {
			break
		}
		if !exists {
			err = fmt.Errorf("Operator.Read: Doesn't exists")
			break
		}
		sd, err = o.readSnapshot(&c)
	}
	if err != nil {
		o.log.Err(err).Str("CarouselId", c.CarId).Msg("Operator.Read")
	}
	return sd, err
}

func (o *Operator) ReadByStatus(status string) ([]SnapshotData, error) {
	return o.evRepo.OperatorReadByStatus(status)

}
func (o *Operator) ReadPending() ([]CompositeData, error) {
	return o.evRepo.OperatorReadPending()
}

func (o *Operator) Tick(t *time.Ticker) {
	var ts time.Time
	var monitorSnapshot time.Duration
	var monitorPending time.Duration
	var monitorOffline time.Duration
	for {
		select {
		case <-t.C:
			monitorOffline += time.Since(ts)
			if monitorOffline > tmMonitorOfflinePeriod {
				monitorOffline = 0
				o.log.Debug().Msg("Monitor Offline")
				if err := o.monitorExpired(); err != nil {
					// o.log.Err(err).Msg("Operator.Tick: Execution monitorExpired failed")
				}
			}
			monitorPending += time.Since(ts)
			if monitorPending > tmMonitorPendingPeriod {
				monitorPending = 0
				o.log.Debug().Msg("Monitor Pending")
				if err := o.monitorPending(); err != nil {
					// o.log.Err(err).Msg("Operator.Tick: Execution monitorPending failed")
				}
			}
			monitorSnapshot += time.Since(ts)
			if monitorSnapshot > tmMonitorSnapshotPeriod {
				monitorSnapshot = 0
				go func() {
					o.log.Debug().Msg("Monitor Snaphsot")
					if err := o.monitorSnapshot(); err != nil {
						// o.log.Err(err).Msg("Operator.Tick: Execution monitorSnapshot failed")
					}
				}()
			}
			ts = time.Now()
		}
	}
}

func (o *Operator) readSnapshot(c *Carousel) (*SnapshotData, error) {
	var err error
	var snapshot *SnapshotData
	var snapshotEvent *SnapshotData

	snapshotEvent, err = o.evRepo.OperatorReadAsSnapshot(c.CarId)
	if snapshot, err = o.snRepo.OperatorLoadSnapshot(c.CarId); err == nil && snapshot != nil && snapshotEvent != nil {
		snapshot.Status = snapshotEvent.Status
		snapshot.Tickets += snapshotEvent.Tickets
	} else {
		snapshot = snapshotEvent
	}

	return snapshot, err
}

func (o *Operator) monitorEventsTable() error {
	return nil
}

func (o *Operator) monitorPending() error {
	var err error
	var experiedArray []CompositeData
	if experiedArray, err = o.evRepo.OperatorReadPending(); err == nil {
		for _, c := range experiedArray {
			if err = o.publish(&RequestPlay{MessageGeneric: MessageGeneric{MsgType: MsgTypeRequestPlay, CarId: c.CarId}, EvtId: c.EvtId.String()}); err == nil {
				o.retries[c.EvtId.String()]++
				o.log.Info().Int("Retries", o.retries[c.EvtId.String()]).Str("CarouselId", c.CarId).Str("EventId", c.EvtId.String()).Msg("Operator.monitorPending: Re-Sent request")
				if o.retries[c.EvtId.String()] > maxRetries {
					if err = o.evRepo.OperatorClearPendingFlag(&PlayData{CarId: c.CarId, EvtId: c.EvtId}); err != nil {
						o.log.Err(err).Str("CarouselId", c.CarId).Str("EventId", c.EvtId.String()).Msg("Operator.monitorPending: Fail to clear pering flag")
					}
				}
			} else {
				o.log.Err(err).Int("Retries", o.retries[c.EvtId.String()]).Str("CarouselId", c.CarId).Str("EventId", c.EvtId.String()).Msg("Operator.monitorPending: Fail to Re-Send request")
			}
		}
	}
	return err
}

func (o *Operator) monitorExpired() error {
	var err error
	var experiedArray []CompositeData
	if experiedArray, err = o.evRepo.OperatorReadExpired(tmCarouselOffline); err == nil {
		for _, c := range experiedArray {
			status := CarouselStatusNameOffline
			if c.Status != nil && *c.Status == status {
				continue
			}
			newStatus := StatusData{CarId: c.CarId, EvtId: uuid.New(), Status: &status}
			o.log.Info().Str("CarouselId", newStatus.CarId).Str("Status", *newStatus.Status).Msg("Operator.monitorExpired: Mark as offline")
			err = o.evRepo.OperatorMark(&newStatus)
		}
	}
	return err
}

func (o *Operator) monitorSnapshot() error {
	var err error
	var ids []string

	for ok := true; ok; ok = false {
		if ids, err = o.crRepo.OperarotReadAllCarouselIds(); err != nil {
			o.log.Err(err).Msg("Operator.monitorSnapshot: Can't find any carousels")
			break
		}
		var cdArray []CompositeData
		var snapshotEvent *SnapshotData
		var snapshot *SnapshotData
		for _, id := range ids {
			if cdArray, err = o.evRepo.OperatorRead(id); err != nil {
				o.log.Err(err).Str("CarouselId", id).Msg("Operator.monitorSnapshot: Can't read a carousels")
				break
			}
			if snapshotEvent, err = o.evRepo.OperatorReadAsSnapshot(id); err != nil {
				o.log.Err(err).Str("CarouselId", id).Msg("Operator.monitorSnapshot: Can't read as a snapshot")
				break
			}
			if snapshot, err = o.snRepo.OperatorLoadSnapshot(id); err != nil || snapshot == nil {
				snapshot = &SnapshotData{CarId: snapshotEvent.CarId, Status: "", Tickets: 0}
			}

			cdArrayLen := len(cdArray) - 1 // skip latest event, so len-1
			if cdArrayLen < CarouselEventMax {
				continue
			}

			for index, cd := range cdArray {
				if index == cdArrayLen {
					break
				}
				if err = o.evRepo.OperatorRemoveByEvent(cd.EvtId); err != nil {
					o.log.Error().Str("CarouselId", id).Str("EventId", cd.EvtId.String()).Msg("Operator.monitorSnapshot: Can't remove a record")
				}
			}
			snapshot.Tickets += snapshotEvent.Tickets
			snapshot.Status = snapshotEvent.Status

			if err = o.snRepo.OperatorStoreSnapshot(snapshot); err != nil {
				o.log.Error().Str("CarouselId", id).Msg("Operator.monitorSnapshot: Can't save a snapshot")
			}
		}

	}

	// var experiedArray []CompositeData
	// if experiedArray, err = o.evRepo.Read(&Carousel{}); err == nil {
	// 	for _, c := range experiedArray {
	// 		status := CarouselStatusNameOffline
	// 		if c.Status != nil && *c.Status == status {
	// 			continue
	// 		}
	// 		newStatus := StatusData{CarId: c.CarId, EvtId: uuid.New(), Status: &status}
	// 		o.log.Info().Str("CarouselId", newStatus.CarId).Str("Status", *newStatus.Status).Msg("Mark as offline")
	// 		err = o.evRepo.Mark(&newStatus)
	// 	}
	// }
	return err
}

func (o *Operator) BrokerNotify(msg IMessageGeneric) {
	var err error
	switch m := msg.(type) {
	case *ResponseAck:
		err = o.ack(m)
	case *EventCompleted:
		err = o.completed(m)
	case *EventHeartbeat:
		err = o.heartbeat(m)
	}
	if err != nil {
		o.log.Err(err).Str("Type", msg.Name()).Msg("Operator.BrokerNotify")
	}
}

func (o *Operator) heartbeat(msg *EventHeartbeat) error {
	var err error
	var status string
	var carouselError *string
	if len(msg.Error) > 0 {
		status = CarouselStatusNameFailure
		carouselError = &msg.Error
	} else {
		carouselError = nil
		status = CarouselStatusNameOnline
	}
	err = o.evRepo.OperatorMark(&StatusData{CarId: msg.CarId, EvtId: uuid.New(), Status: &status, Error: carouselError})
	return err
}

func (o *Operator) completed(msg *EventCompleted) error {
	var err error
	if len(msg.Error) > 0 {
		status := CarouselStatusNameFailure
		o.log.Warn().Str("MsgType", msg.MsgType).Str("CarouselId", msg.CarId).Str("Error", msg.Error).Str("Status", status).Msg("Operator.complete: Scenario has been completed with failure")
		err = o.evRepo.OperatorMark(&StatusData{CarId: msg.CarId, EvtId: uuid.New(), Status: &status, Error: &msg.Error})
	} else {
		o.log.Info().Str("MsgType", msg.MsgType).Str("CarouselId", msg.CarId).Msg("Operator.complete: Scenario has been completed successfully")
	}
	return err
}

func (o *Operator) ack(msg *ResponseAck) error {
	var err error
	var correlationId uuid.UUID
	if correlationId, err = uuid.Parse(msg.CorId); err == nil {
		if len(msg.Error) > 0 {
			status := CarouselStatusNameFailure
			o.log.Warn().Str("MsgType", msg.MsgType).Str("CarouselId", msg.CarId).Str("CorrelationId", msg.CorId).Str("Error", msg.Error).Str("Status", status).Msg("Operator.ack: Carousel is not operable")
			err = o.evRepo.OperatorMark(&StatusData{CarId: msg.CarId, EvtId: uuid.New(), Status: &status, Error: &msg.Error})
		} else {
			o.log.Info().Str("MsgType", msg.MsgType).Str("CarouselId", msg.CarId).Str("CorrelationId", msg.CorId).Msg("Operator.ack: Confirmed")
			err = o.evRepo.OperatorConfirm(&StatusData{CarId: msg.CarId, EvtId: correlationId})
		}
	}
	return err
}

func (o *Operator) publish(msg IMessageGeneric) error {
	t := topic.New(o.config.RootTopicPub())
	t.Appned(msg.Target())
	return o.broker.Publish(t, msg, o.config.DefaultQOS())
}
