package accountment_event

import "github.com/google/uuid"

type EventRidesUpdated struct {
	uuid       uuid.UUID
	carouselId uuid.UUID
	rides      uint
}

func EventRidesUpdatedCreate(rides uint, carouselId uuid.UUID) *EventRidesUpdated {
	return &EventRidesUpdated{
		uuid:       uuid.New(),
		carouselId: carouselId,
		rides:      rides,
	}
}

func EventRidesUpdatedCreateEmpty() *EventRidesUpdated {
	return &EventRidesUpdated{}
}

func (e *EventRidesUpdated) Name() string {
	return "accountment.event.balance.updated"
}

func (e *EventRidesUpdated) Id() uuid.UUID {
	return e.uuid
}

func (e *EventRidesUpdated) Rides() uint {
	return e.rides
}

func (e *EventRidesUpdated) CarouselId() uuid.UUID {
	return e.carouselId
}
