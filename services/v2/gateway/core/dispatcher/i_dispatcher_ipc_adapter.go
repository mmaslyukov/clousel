package dispatcher

import (
	"gateway/lib/fault"
	"time"
)

type IDispatcherIpcController interface {
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

type IDispatcherIpcAdapter interface {
	Subscribe(subj string, listener IDispatcherIpcController) fault.IError
	QueueSubscribe(subj, queue string, listener IDispatcherIpcController) fault.IError
	Publish(subj string, data []byte) fault.IError
	Request(subj string, data []byte, timeout time.Duration) (IIpcMessage, fault.IError)
}
