package operator

type IPortOperatorControllerMqtt interface {
	BrokerNotify(msg IMessageGeneric)
}
