package machine

// Barouness: BAck gROUNd procESS

import (
	"bytes"
	"clousel/lib/fault"
	"encoding/gob"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	natsSubjectPlay   = "Gw.Play"
	natsSubjectStatus = "Gw.Status"
)

type IMachineBarounessController interface {
	// Change Game Status to Completed
	// Add Entry to Balance table

	/* Games */
	ReadPendingGames() ([]*GameEvent, fault.IError)
	ReadExpiredGames(ts time.Time) ([]*GameEvent, fault.IError)
	GameStartConfirm(gameId uuid.UUID) fault.IError
	GameStartFailed(gameId uuid.UUID) fault.IError
	/* Machine */
	MachineUpdateStatus(machId uuid.UUID, status MachineStatus) fault.IError
}

type IMachineBarounessAdapter interface {
	SendGameRequestEvent(event GameEvent)
}

type Barounes struct {
	gameInput chan GameEvent
	machCtrl  IMachineBarounessController
	ipc       IMachineIpcAdapter
	log       *zerolog.Logger
	cfg       IMachineConfigAdapter
}

func BarounessCreatePartial(log *zerolog.Logger, cfg IMachineConfigAdapter, ipc IMachineIpcAdapter) *Barounes {
	return &Barounes{
		gameInput: make(chan GameEvent, 100),
		log:       log,
		cfg:       cfg,
		ipc:       ipc,
		// tickInput: make(chan any),
	}
}

func DecodeBinary[T any](data *bytes.Buffer) (obj T, err fault.IError) {
	e := gob.NewDecoder(data).Decode(&obj)
	if e != nil {
		err = fault.New(EMachineDecode).Msg(e.Error())
	}
	return obj, err
}
func EncodeBinary[T any](obj *T) (data bytes.Buffer, err fault.IError) {
	e := gob.NewEncoder(&data).Encode(obj)
	if e != nil {
		err = fault.New(EMachineEncode).Msg(e.Error())
	}
	return data, err
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

func (b *Barounes) InjectAndFinish(machCtrl IMachineBarounessController) *Barounes {
	b.machCtrl = machCtrl
	return b
}

func (b *Barounes) Init() (err fault.IError) {
	const fn = "Core.Machine.Barouness.Init"
	const subject = natsSubjectStatus
	if err = b.ipc.Subscribe(subject, b); err != nil {
		b.log.Error().Msgf("%s: Fail to subscribe on subject:%s, error:%s ", fn, subject, err.Error())
	} else {
		b.log.Info().Msgf("%s: Subscribed on subject: %s", fn, subject)
	}
	return err
}

func (b *Barounes) IpcNotify(msg IIpcMessageReplyable) {
	const fn = "Core.Machine.Barouness.IpcNotify"
	if msg.Subject() == natsSubjectStatus {
		b.log.Debug().Msgf("%s: '%s' has %d bytes: %s", fn, msg.Subject(), len(msg.Data()), string(msg.Data()))
		if ms, err := DecodeJson[RemoteMachineStatus](bytes.NewBuffer(msg.Data())); err == nil {
			if err = b.machCtrl.MachineUpdateStatus(ms.MachId, ms.Status); err != nil {
				b.log.Err(err).Msgf("%s: Fail to update machine(%s) status(%s)", fn, ms.MachId.String(), ms.Status)
			}
		} else {
			b.log.Err(err).Msgf("%s: Fail to decode Json", fn)
		}
	}
}

func (b *Barounes) SendGameRequestEvent(event GameEvent) {
	b.gameInput <- event
}

func (b *Barounes) handleGameEvent(e *GameEvent) {
	const fn = "Core.Machine.Barouness.handleGameEvent"
	b.log.Debug().Msgf("%s: Got game event %+v", fn, e)
	for ok := true; ok; ok = false {
		var err error
		buf, err := EncodeJson(e)
		var resp IIpcMessage
		if resp, err = b.ipc.Request(natsSubjectPlay, buf.Bytes(), b.cfg.GameStartTimeout()); err != nil {
			b.log.Err(err).Msgf("%s: Unable to request subject: %s", fn, natsSubjectPlay)
			break
		}
		if resp == nil {
			b.log.Error().Msgf("%s: Unable to read response form subject: %s", fn, natsSubjectPlay)
			break
		}
		var r GameEventAck
		if r, err = DecodeJson[GameEventAck](bytes.NewBuffer(resp.Data())); err != nil {
			b.log.Err(err).Msgf("%s: Unable to decode game event ack", fn)
			break
		}
		if r.CorrelationId != e.EventId {
			b.log.Warn().Msgf("%s: Got unexpected correlation id", fn)
			break
		}

		if r.Code == 0 {
			err = b.machCtrl.GameStartConfirm(r.CorrelationId)
		} else {
			err = b.machCtrl.GameStartFailed(r.CorrelationId)
		}
		if err != nil {
			b.log.Err(err).Msgf("%s: Failed to update game status", fn)
		}
	}
}
func (b *Barounes) monitorExpiredGames() {
	const fn = "Core.Machine.Barouness.monitorExpiredGames"
	if games, err := b.machCtrl.ReadExpiredGames(time.Now().Add(-b.cfg.GameStartTimeout())); err == nil {
		for _, g := range games {
			b.log.Warn().Msgf("%s: Mark game id %s as failed", fn, g.EventId.String())
			b.machCtrl.GameStartFailed(g.EventId)
		}
	}
}

func (b *Barounes) Run() {
	const fn = "Core.Machine.Barouness.Run"
	b.log.Info().Msgf("%s: Started", fn)
	tm := time.NewTicker(b.cfg.GameStartMonitor())

	for {
		select {
		case e := <-b.gameInput:
			b.handleGameEvent(&e)
			break
		// case <-b.tickInput:
		// 	break
		case <-tm.C:
			b.monitorExpiredGames()
			// select games with status pending and timestamp earlier than tm
			break
		}
	}
}
