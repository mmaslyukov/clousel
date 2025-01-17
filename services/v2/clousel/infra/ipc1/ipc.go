package ipc

import (
	"clousel/lib/fault"

	"github.com/nats-io/nats.go"
)

type INatsConfig interface {
	GetNatsUrl() string
}

type Ipc struct {
	cfg INatsConfig
}

func IpcCreate(cfg INatsConfig) *Ipc {
	return &Ipc{cfg: cfg}
}
func (i *Ipc) Connect() (nc *nats.Conn, err fault.IError) {
	const fn = "Ipc.Connect"
	if c, e := nats.Connect(i.cfg.GetNatsUrl()); e == nil {
		nc = c
	} else {
		err = fault.New(EIpcConnect).Msgf("%s: Fail to connect to '%s', err:%s", fn, i.cfg.GetNatsUrl(), e.Error())
	}
	return nc, err
}
