package port_broker

import (
	. "carousel_service/internal/ports"
)

type MessageMinimal struct {
	Type        string `json:"Type"`
	CarouselId  string `json:"CarouselId"`
	SequenceNum int    `json:"SequenceNum"`
	EventId     string `json:"EventId"`
}

type EventMinimal struct {
	Type        string `json:"Type"`
	CarouselId  string `json:"CarouselId"`
	SequenceNum int    `json:"SequenceNum"`
}

type EventAck struct {
	EventMinimal
	CorrelationId string `json:"CorrelationId"`

	// optional field for showing status of the last command
	Error string `json:"Error"`
}

func NewEventAck() EventAck {
	return EventAck{EventMinimal: EventMinimal{Type: "EventAck"}}
}

type EventHeartbeat struct {
	EventMinimal
	// optional field for collecting active errors
	Error string `json:"Error"`
}

func NewEventHeartbeat() EventHeartbeat {
	return EventHeartbeat{EventMinimal: EventMinimal{Type: "EventHeartbeat"}}
}

type MessageCommand struct {
	MessageMinimal
	Command string `json:"Command"`
}

func NewMessageCommand(cmd string) *MessageCommand {
	return &MessageCommand{
		MessageMinimal: MessageMinimal{
			Type: "MessageCommand",
		},
		Command: cmd,
	}
}

type MessageFeedDelme struct {
	Feed int `json:"feed"`
}

type Subscribers struct {
	Carousels map[string]byte
}

type BrokerInterface interface {
	SetSubscribers(subs Subscribers)
	PublishQueue(data any)
	ListenQueue() PortInterface[any]
}

type InfoCommandInterface interface {
	populateDefault() InfoCommandInterface
}

// func (e *EventPlay) populateDefault() InfoCommandInterface {
// 	e.Command = "Play"
// 	return e
// }

// type BrokerRunnerInterface interface {
// 	Connect() error
// 	Run()
// }

// type BrokerPrototypeInterface interface {
// 	CloneBroker() BrokerInterface
// 	CloneRunner() BrokerRunnerInterface
// }

// var Prototype BrokerPrototypeInterface

// func Prototype() BrokerPrototypeInterface {
// 	return prototype
// }
