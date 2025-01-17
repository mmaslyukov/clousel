package ipc

import (
	"clousel/core/machine"
	"clousel/lib/fault"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type IpcMessage struct {
	msg *nats.Msg
}

func (m *IpcMessage) Data() []byte {
	return m.msg.Data
}
func (m *IpcMessage) Subject() string {
	return m.msg.Subject
}

func (m *IpcMessage) Respond(data []byte) (err fault.IError) {
	if e := m.msg.Respond(data); e != nil {
		err = fault.New(EIpcNats).Err(e)
	}
	return err
}

type INatsConfig interface {
	GetNatsUrl() string
}

type Ipc struct {
	opts           []nats.Option
	conn           *nats.Conn
	log            *zerolog.Logger
	cfg            INatsConfig
	listeners      map[string][]machine.IMachineIpcController
	listenersQueue map[string][]machine.IMachineIpcController
}

func IpcCreate(cfg INatsConfig, log *zerolog.Logger) *Ipc {
	const fn = "infra.ipc.IpcCreate"
	opts := []nats.Option{nats.Name("NATS Gateway")}
	opts = setupConnOptions(opts, log)
	ipc := &Ipc{
		opts:      opts,
		log:       log,
		cfg:       cfg,
		listeners: make(map[string][]machine.IMachineIpcController),
	}
	ipc.Connect()
	return ipc
}

func setupConnOptions(opts []nats.Option, log *zerolog.Logger) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		const fn = "infra.ipc.DisconnectErrHandler"
		log.Info().Msgf("%s: Disconnected due to:%s, will attempt reconnects for %.0fm", fn, err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		const fn = "infra.ipc.ReconnectHandler"
		log.Info().Msgf("%s: Reconnected [%s]", fn, nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		const fn = "infra.ipc.ClosedHandler"
		log.Error().Msgf("%s: Exiting: %v", fn, nc.LastError())

	}))
	return opts
}

func (i *Ipc) Connect() (err fault.IError) {
	const fn = "infra.ipc.Connect"
	var e error
	if i.conn, e = nats.Connect(i.cfg.GetNatsUrl(), i.opts...); e != nil {
		log.Err(e).Msgf("%s: Fail to connect", fn)
		err = fault.New(EIpcConnection).Err(e)
	}
	return err
}

func (i *Ipc) Subscribe(subj string, listener machine.IMachineIpcController) (err fault.IError) {
	const fn = "infra.ipc.Subscribe"
	if i.conn != nil {
		i.listeners[subj] = append(i.listeners[subj], listener)
		if _, e := i.conn.Subscribe(subj, func(msg *nats.Msg) {
			for _, v := range i.listeners {
				for _, l := range v {
					l.IpcNotify(&IpcMessage{msg: msg})
				}
			}
		}); e != nil {
			err = fault.New(EIpcNats).Err(e).Msgf("Fail to subscribe to subject %s", subj)
		}
	}
	return err
}

func (i *Ipc) QueueSubscribe(subj, queue string, listener machine.IMachineIpcController) (err fault.IError) {
	const fn = "infra.ipc.QueueSubscribe"

	if i.conn != nil {
		i.listenersQueue[subj] = append(i.listenersQueue[subj], listener)
		if _, e := i.conn.Subscribe(subj, func(msg *nats.Msg) {
			for _, v := range i.listenersQueue {
				for _, l := range v {
					l.IpcNotify(&IpcMessage{msg: msg})
				}
			}
		}); e != nil {
			err = fault.New(EIpcNats).Err(e).Msgf("Fail to subscribe to subject %s", subj)
		}
	}
	return err
}

func (i *Ipc) Publish(subj string, data []byte) (err fault.IError) {
	// if i.conn.IsClosed() || !i.conn.IsConnected() {
	// 	i.Connect()
	// }
	if e := i.conn.Publish(subj, data); e != nil {
		err = fault.New(EIpcNats).Err(e)
	}
	return err
}

func (i *Ipc) Request(subj string, data []byte, timeout time.Duration) (resp machine.IIpcMessage, err fault.IError) {
	// if i.conn.IsClosed() || !i.conn.IsConnected() {
	// 	i.Connect()
	// }
	if r, e := i.conn.Request(subj, data, timeout); e != nil {
		err = fault.New(EIpcNats).Err(e)
	} else {
		resp = &IpcMessage{msg: r}
	}
	return resp, err
}
