package machine

import (
	"clousel/lib/fault"
	"time"
)

type IMachineIpcController interface {
	IpcNotify(msg IIpcMessageReplyable)
}

type IIpcMessage interface {
	Data() []byte
	Subject() string
}

type IIpcMessageReplyable interface {
	IIpcMessage
	Respond(data []byte) fault.IError
}

type IMachineIpcAdapter interface {
	Subscribe(subj string, listener IMachineIpcController) fault.IError
	QueueSubscribe(subj, queue string, listener IMachineIpcController) fault.IError
	Publish(subj string, data []byte) fault.IError
	Request(subj string, data []byte, timeout time.Duration) (IIpcMessage, fault.IError)
}
