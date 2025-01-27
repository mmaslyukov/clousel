package carousel_provider

import (
	"accountant_service/domain/carousel"

	"github.com/google/uuid"
)

type IPortApiRide interface {
	CheckAndPublishUndeliveredRides(carouselId uuid.UUID) ([]carousel.ResponseRefillRides, error)
	ReadUndeliveredRides(carouselId uuid.UUID) ([]carousel.RideMinimal, error)
}
