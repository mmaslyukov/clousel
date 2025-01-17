package dispatcher

import (
	"bytes"
	"encoding/json"
	"gateway/infra/broker/topic"
	"gateway/lib/fault"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	natsSubjectPlay   = "Gw.Play"
	natsSubjectStatus = "Gw.Status"
)

type GameMetadata struct {
	game  *GameEvent
	reply IIpcMessageReplyable
}
type MachineAvailabilityStatus struct {
	machId     uuid.UUID
	time       time.Time
	statusCurr MachineStatus
	statusPrev MachineStatus
}

func MachineAvailabilityStatusCreateDefault(machId uuid.UUID) *MachineAvailabilityStatus {
	return &MachineAvailabilityStatus{
		machId:     machId,
		time:       time.Now(),
		statusCurr: MachineStatusInvalid,
		statusPrev: MachineStatusInvalid,
	}
}
func (m *MachineAvailabilityStatus) updateStatus(newStatus MachineStatus) *MachineAvailabilityStatus {
	m.statusPrev = m.statusCurr
	m.statusCurr = newStatus
	m.time = time.Now()
	return m
}

func (m *MachineAvailabilityStatus) isChanged() bool {
	return m.statusCurr != m.statusPrev

}

type Dispatcher struct {
	broker IDispatcherMqttAdapter
	cfg    IDispatcherConfigAdapter
	log    *zerolog.Logger
	ipc    IDispatcherIpcAdapter
	meta   map[uuid.UUID]*GameMetadata
	macs   map[uuid.UUID]*MachineAvailabilityStatus
}

func DispatcherCreate(broker IDispatcherMqttAdapter, ipc IDispatcherIpcAdapter, cfg IDispatcherConfigAdapter, log *zerolog.Logger) *Dispatcher {
	const fn = "Core.Dispatcher.DispatcherCreate"
	d := &Dispatcher{
		broker: broker,
		cfg:    cfg,
		log:    log,
		ipc:    ipc,
		meta:   make(map[uuid.UUID]*GameMetadata),
		macs:   make(map[uuid.UUID]*MachineAvailabilityStatus),
	}
	if err := broker.Subscribe(topic.New(cfg.RootTopicSub()), cfg.DefaultQOS(), d); err != nil {
		log.Err(err).Msgf("%s: Fail to subscribe", fn)
	}
	const subject = natsSubjectPlay
	if err := d.ipc.Subscribe(subject, d); err != nil {
		d.log.Error().Msgf("%s: Fail to subscribe on subject:%s, error:%s ", fn, subject, err.Error())
	} else {
		d.log.Info().Msgf("%s: Subscribed on subject: %s", fn, subject)
	}
	return d
}
func EncodeJson[T any](obj *T) (data bytes.Buffer, err fault.IError) {
	e := json.NewEncoder(&data).Encode(obj)
	if e != nil {
		err = fault.New(EMachineDecode).Msg(e.Error())
	}
	return data, err
}
func DecodeJson[T any](data *bytes.Buffer) (obj T, err fault.IError) {
	e := json.NewDecoder(data).Decode(&obj)
	if e != nil {
		err = fault.New(EMachineEncode).Msg(e.Error())
	}
	return obj, err
}

func (d *Dispatcher) IpcNotify(msg IIpcMessageReplyable) {
	const fn = "Core.Dispatcher.IpcNotify"
	if msg.Subject() == natsSubjectPlay {
		d.IpcNotifyPlay(msg)
	}
}

func (d *Dispatcher) IpcNotifyPlay(msg IIpcMessageReplyable) {
	const fn = "Core.Dispatcher.IpcNotifyPlay"
	o, err := DecodeJson[GameEvent](bytes.NewBuffer(msg.Data()))
	if err != nil {
		return
	}
	mm := GameMetadata{game: &o, reply: msg}
	err = d.publish(&RequestPlay{
		MessageGeneric: MessageGeneric{
			MsgType: MsgTypeRequestPlay,
			CarId:   o.MachId,
		},
		EvtId: o.EventId,
	})
	if err != nil {
		d.respondWithGameEventAck(msg, &o, GameEventAckPublishFail)
		return
	}
	d.log.Debug().Msgf("%s Save message metadata %+v", fn, mm.game)
	d.meta[o.EventId] = &mm
}

func (d *Dispatcher) respondWithGameEventAck(msg IIpcMessageReplyable, ge *GameEvent, code GameAckCode) fault.IError {
	const fn = "Core.Dispatcher.BrokerNotify"
	var err fault.IError
	var buf bytes.Buffer
	gea := GameEventAck{MachId: ge.MachId, CorrelationId: ge.EventId, Code: code}
	if buf, err = EncodeJson(&gea); err == nil {
		if e := msg.Respond(buf.Bytes()); e != nil {
			err = fault.New(EMachineRespond).Err(e)
			d.log.Err(e)
		} else {
			d.log.Debug().Msgf("%s: Sent reposnse to subject: %s, %s", fn, msg.Subject(), string(buf.Bytes()))
		}
	} else {
		d.log.Err(err)
	}
	return err
}

func (d *Dispatcher) BrokerNotify(msg IMessageGeneric) {
	const fn = "Core.Dispatcher.BrokerNotify"
	switch m := msg.(type) {
	case *ResponseAck:
		d.BrokerNotifyAck(m)
	// case *EventCompleted:
	case *EventHeartbeat:
		d.BrokerNotifyHeartbeat(m)
	}

}

func (d *Dispatcher) BrokerNotifyHeartbeat(m *EventHeartbeat) {
	const fn = "Core.Dispatcher.BrokerNotifyHeartbeat"
	d.log.Debug().Msgf("%s: Got HB message from %s", fn, m.CarId)
	if _, ok := d.macs[m.CarId]; !ok {
		d.macs[m.CarId] = MachineAvailabilityStatusCreateDefault(m.CarId)
	}
	d.macs[m.CarId].updateStatus(MachineStatusOnline)
	d.log.Debug().Msgf("%s: machId:%s, statusPrev:%s, statusCurr:%s, ", fn,
		d.macs[m.CarId].machId, d.macs[m.CarId].statusPrev, d.macs[m.CarId].statusCurr)
}

func (d *Dispatcher) BrokerNotifyAck(m *ResponseAck) {
	const fn = "Core.Dispatcher.BrokerNotifyAck"
	d.log.Debug().Msgf("%s: Ack %+v", fn, m)
	if val, ok := d.meta[m.CorId]; ok {
		var code GameAckCode
		if len(m.Error) > 0 {
			code = GameEventAckRemoteError
		} else {
			code = GameEventAckOk
		}
		if err := d.respondWithGameEventAck(val.reply, val.game, code); err != nil {
			d.log.Err(err).Msgf("%s", fn)
		}
		delete(d.meta, m.CorId)
	}
}

// func (d *Dispatcher) IpcDisconnected(con *nats.Conn, err error) {
// 	const fn = "Core.Dispatcher.IpcDisconnected"
// 	if err != nil {
// 		d.log.Err(err).Msgf("%s: on disconnect callback", fn)
// 	}

// 	if d.con, err = d.nats.Connect(); err != nil {
// 		d.log.Err(err).Msgf("%s: at reconnect", fn)
// 	}
// }

func (d *Dispatcher) publish(msg IMessageGeneric) fault.IError {
	t := topic.New(d.cfg.RootTopicPub())
	t.Appned(msg.Target().String())
	return d.broker.Publish(t, msg, d.cfg.DefaultQOS())
}
func (d *Dispatcher) monitorMachineForOffline() {
	for _, v := range d.macs {
		if time.Since(v.time) >= d.cfg.GetMachineExpiredTime() &&
			v.statusCurr != MachineStatusOffline {
			v.updateStatus(MachineStatusOffline)
		}
	}
}
func (d *Dispatcher) monitorAndUpdateMachineStatus() {
	const fn = "Core.Dispatcher.monitorAndUpdateMachineStatus"
	for _, v := range d.macs {
		if v.isChanged() {
			v.updateStatus(v.statusCurr)
			rms := RemoteMachineStatus{MachId: v.machId, Status: v.statusCurr}
			if buf, err := EncodeJson(&rms); err == nil {
				d.log.Debug().Msgf("%s: About to publish to '%s' %s", fn, natsSubjectStatus, buf.String())

				if e := d.ipc.Publish(natsSubjectStatus, buf.Bytes()); e != nil {
					d.log.Err(e).Msgf("%s: Fail to publish machine(carousel) status", fn)
				}
			} else {
				d.log.Err(err).Msgf("%s: Fail to encode machine(carousel) status", fn)
			}
		}
	}
}
func (d *Dispatcher) Run() {
	const fn = "Core.Dispatcher.Run"
	tm := time.NewTicker(d.cfg.GetMachineMonitorTm())

	d.log.Info().Msgf("%s: Started", fn)
	for {
		select {
		case <-tm.C:
			d.log.Debug().Msgf("%s: Tick", fn)
			d.monitorMachineForOffline()
			d.monitorAndUpdateMachineStatus()
			break
		}
	}
}
