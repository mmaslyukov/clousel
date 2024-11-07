package port_carousel

import (
	. "carousel_service/internal/utils"
)

type CarouselId struct {
	Id string
}

type RoundData struct {
	CarouselId
	RoundsReady Optional[int]
}

type RegisterData struct {
	CarouselId
	OwnerId   string
	RoundTime Optional[int]
}

type RefillData struct {
	CarouselId
	RoundsReady int
}

type PlayData struct {
	CarouselId
}

type ControlData struct {
	CarouselId
}

type AggregationData struct {
	RoundData
	RoundTime Optional[int]
	Status    string
	Time      string
}

type CarouselInterface interface {
	// Rest API threads
	Register(data RegisterData) error
	Delete(data CarouselId) error
	Read(data CarouselId) (Optional[AggregationData], error)
	ReadByOwner(ownerId string) (Optional[[]AggregationData], error)
	Refill(data RefillData) error
	Play(data CarouselId) error
}
