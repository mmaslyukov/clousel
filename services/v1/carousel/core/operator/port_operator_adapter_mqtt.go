package operator

type IPortOperatorAdapterMqtt interface {
	Publish(topic Itopic, msg IMessageGeneric, qos byte) error
	Subscribe(topic Itopic, qos byte, listener IPortOperatorControllerMqtt) error
}
