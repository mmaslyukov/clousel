package carousel

import (
	"github.com/google/uuid"
)

type RideMinimal struct {
	CarouselId uuid.UUID
	Rides      uint
}

type ResponseRefillRides struct {
	RideMinimal
	Success bool
	Error   string
}
