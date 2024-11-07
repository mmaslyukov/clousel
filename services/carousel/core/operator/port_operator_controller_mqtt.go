package operator

type IPortOperatorControllerMqtt interface {
	Notify(msg IMessageGeneric)
}
