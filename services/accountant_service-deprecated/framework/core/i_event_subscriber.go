package core

type IEventSubscribable interface {
	Subscribe(event IEvent, listener IEventListener)
}
