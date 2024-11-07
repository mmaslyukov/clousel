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
	repo    IPortOperatorAdapterRepository
	broker  IPortOperatorAdapterMqtt
	config  IPortOperatorAdapterConfig
	log     *zerolog.Logger
	retries map[string]int
}

func New(
	repo IPortOperatorAdapterRepository,
	broker IPortOperatorAdapterMqtt,
	config IPortOperatorAdapterConfig,
	log *zerolog.Logger) *Operator {

	op := &Operator{repo: repo, broker: broker, config: config, log: log, retries: make(map[string]int)}
	return op
}

func (o *Operator) Refill(c Carousel, rounds int) error {
	var err error
	var exists bool

	for ok := true; ok; ok = false {
		if rounds < 1 {
			err = fmt.Errorf("Operator.Refill: Invalid Rounds value: %d", rounds)
			break
		}
		if exists, err = o.repo.IsExists(c); err != nil {
			break
		}
		if !exists {
			err = fmt.Errorf("Operator.Refill: Doesn't exists")
			break
		}
		rd := RoundsData{CarId: c.CarId, Rounds: rounds, EvtId: uuid.New()}
		o.log.Info().Str("CarouselId", rd.CarId).Str("EventId", rd.EvtId.String()).Int("Rounds", rd.Rounds).Msg("Operator.Refill: About to write an event")
		err = o.repo.Refill(&rd)
	}
	return err
}

func (o *Operator) Play(c Carousel) error {
	var err error
	var exists bool

	for ok := true; ok; ok = false {
		if exists, err = o.repo.IsExists(c); err != nil {
			break
		}
		if !exists {
			err = fmt.Errorf("Operator.Play: Doesn't exists")
			break
		}
		// TODO add read shapshot table
		var s SnapshotData
		if s, err = o.repo.ReadAsSnapshot(&c); err != nil {
			break
		}
		// TODO remove false
		if false && s.Status != CarouselStatusNameOnline {
			err = fmt.Errorf("Operator.Play: Carousel Status is '%s'", s.Status)
			break
		}
		if s.Rounds == 0 {
			err = fmt.Errorf("Operator.Play: Carousel has no Rounds (%d Rounds)", s.Rounds)
			break
		}

		pd := PlayData{CarId: c.CarId, EvtId: uuid.New()}
		o.log.Info().Str("CarouselId", pd.CarId).Str("EventId", pd.EvtId.String()).Msg("Operator.Play: About to write an event")
		if err = o.repo.Play(&pd); err != nil {
			break
		}
		o.log.Info().Str("CarouselId", pd.CarId).Str("EventId", pd.EvtId.String()).Str("Type", MsgTypeRequestPlay).Msg("Operator.Play: About to publish the Play event")
		err = o.publish(&RequestPlay{MessageGeneric: MessageGeneric{MsgType: MsgTypeRequestPlay, CarId: pd.CarId}, EvtId: pd.EvtId.String()})
	}
	return err
}

func (o *Operator) Read(c Carousel) (SnapshotData, error) {
	var err error
	var exists bool
	var sd SnapshotData
	for ok := true; ok; ok = false {
		if exists, err = o.repo.IsExists(c); err != nil {
			break
		}
		if !exists {
			err = fmt.Errorf("Operator.Read: Doesn't exists")
			break
		}
		// TODO add read shapshot table
		sd, err = o.repo.ReadAsSnapshot(&c)
	}
	if err != nil {
		o.log.Err(err).Str("CarouselId", c.CarId).Msg("Operator.Read")
	}
	return sd, err
}

func (o *Operator) ReadWStatus(status string) ([]SnapshotData, error) {
	return o.repo.ReadWStatus(status)

}
func (o *Operator) ReadPending() ([]CompositeData, error) {
	return o.repo.ReadPending()
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
					o.log.Err(err).Msg("Operator.Tick: Execution monitorExpired failed")
				}
			}
			monitorPending += time.Since(ts)
			if monitorPending > tmMonitorPendingPeriod {
				monitorPending = 0
				o.log.Debug().Msg("Monitor Pending")
				if err := o.monitorPending(); err != nil {
					o.log.Err(err).Msg("Operator.Tick: Execution monitorPending failed")
				}
			}
			monitorSnapshot += time.Since(ts)
			if monitorSnapshot > tmMonitorSnapshotPeriod {
				monitorSnapshot = 0
				o.log.Debug().Msg("Monitor Snaphsot")
			}
			ts = time.Now()
		}
	}
}

func (o *Operator) monitorPending() error {
	var err error
	var experiedArray []CompositeData
	if experiedArray, err = o.repo.ReadPending(); err == nil {
		for _, c := range experiedArray {
			if err = o.publish(&RequestPlay{MessageGeneric: MessageGeneric{MsgType: MsgTypeRequestPlay, CarId: c.CarId}, EvtId: c.EvtId.String()}); err == nil {
				o.retries[c.EvtId.String()]++
				o.log.Info().Int("Retries", o.retries[c.EvtId.String()]).Str("CarouselId", c.CarId).Str("EventId", c.EvtId.String()).Msg("Re-Sent request")
				if o.retries[c.EvtId.String()] > maxRetries {
					if err = o.repo.ClearPendingFlag(&PlayData{CarId: c.CarId, EvtId: c.EvtId}); err != nil {
						o.log.Err(err).Str("CarouselId", c.CarId).Str("EventId", c.EvtId.String()).Msg("Fail to clear pering flag")
					}
				}
			} else {
				o.log.Err(err).Int("Retries", o.retries[c.EvtId.String()]).Str("CarouselId", c.CarId).Str("EventId", c.EvtId.String()).Msg("Fail to Re-Send request")
			}
		}
	}
	return err
}

func (o *Operator) monitorExpired() error {
	var err error
	var experiedArray []CompositeData
	if experiedArray, err = o.repo.ReadExpired(tmCarouselOffline); err == nil {
		for _, c := range experiedArray {
			status := CarouselStatusNameOffline
			if c.Status != nil && *c.Status == status {
				continue
			}
			newStatus := StatusData{CarId: c.CarId, EvtId: uuid.New(), Status: &status}
			o.log.Info().Str("CarouselId", newStatus.CarId).Str("Status", *newStatus.Status).Msg("Mark as offline")
			err = o.repo.Mark(&newStatus)
		}
	}
	return err
}

func (o *Operator) Notify(msg IMessageGeneric) {
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
		o.log.Err(err).Str("Type", msg.Name()).Msg("Operator.Notify")
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
	// o.log.Debug().Str(msg.CarId).Str(msg.)
	// var snapshot SnapshotData
	// snapshot, err = o.repo.ReadAsSnapshot(&Carousel{CarId: msg.CarId})
	// if status != snapshot.Status {
	err = o.repo.Mark(&StatusData{CarId: msg.CarId, EvtId: uuid.New(), Status: &status, Error: carouselError})
	// } else if snapshot.Status == CarouselStatusNameOnline {
	// err = o.repo.UpdateTime(&Carousel{CarId: msg.CarId})
	// }
	return err
}

func (o *Operator) completed(msg *EventCompleted) error {
	var err error
	if len(msg.Error) > 0 {
		status := CarouselStatusNameFailure
		o.log.Warn().Str("MsgType", msg.MsgType).Str("CarouselId", msg.CarId).Str("Error", msg.Error).Str("Status", status).Msg("Operator.complete: Scenario has been completed with failure")
		err = o.repo.Mark(&StatusData{CarId: msg.CarId, EvtId: uuid.New(), Status: &status, Error: &msg.Error})
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
			err = o.repo.Mark(&StatusData{CarId: msg.CarId, EvtId: uuid.New(), Status: &status, Error: &msg.Error})
		} else {
			o.log.Info().Str("MsgType", msg.MsgType).Str("CarouselId", msg.CarId).Str("CorrelationId", msg.CorId).Msg("Operator.ack: Confirmed")
			err = o.repo.Confirm(&StatusData{CarId: msg.CarId, EvtId: correlationId})
		}
	}
	return err
}

func (o *Operator) publish(msg IMessageGeneric) error {
	t := topic.New(o.config.RootTopicPub())
	t.Appned(msg.Target())
	return o.broker.Publish(t, msg, o.config.DefaultQOS())
}
