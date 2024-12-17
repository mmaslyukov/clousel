package core

type IEventListener interface {
	HandleEvent(IEvent) error
}
