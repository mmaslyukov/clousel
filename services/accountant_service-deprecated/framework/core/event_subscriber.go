package core

import (
	"accountant_service/framework/core/core_error"
)

type EventSubscribable struct {
	subscribers map[string][]IEventListener
}

func EventSubscribableCreate() *EventSubscribable {
	return &EventSubscribable{
		subscribers: make(map[string][]IEventListener),
	}
}

func (e *EventSubscribable) Subscribe(event IEvent, listener IEventListener) {
	slice, ok := e.subscribers[event.Name()]
	if !ok {
		e.subscribers[event.Name()] = make([]IEventListener, 0)
	}
	slice = append(slice, listener)
	e.subscribers[event.Name()] = slice
}

func (e *EventSubscribable) Notify(event IEvent) error {
	var err error
	if event == nil {
		return ErrorCreate[core_error.ErrorNil]().Message("Cannot notify, event is nil")
	}
	for _, listener := range e.subscribers[event.Name()] {
		err = listener.HandleEvent(event)
		if err != nil {
			break
		}
	}
	return err
}
