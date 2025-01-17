package dispatcher

import "gateway/lib/fault"

type IDispatcherMqttController interface {
	BrokerNotify(msg IMessageGeneric)
}

type IDispatcherMqttAdapter interface {
	Publish(topic ITopic, msg IMessageGeneric, qos byte) fault.IError
	Subscribe(topic ITopic, qos byte, listener IDispatcherMqttController) fault.IError
}
