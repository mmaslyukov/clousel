package carousel_consumer

import (
	"accountant_service/domain/carousel/carousel_aggregate"

	"github.com/google/uuid"
)

type IPortRepositoryRide interface {
	Save(ride carousel_aggregate.Ride) error
	LoadUndelivered(caroselId uuid.UUID) ([]carousel_aggregate.Ride, error)
	LoadDelivered(caroselId uuid.UUID) ([]carousel_aggregate.Ride, error)
}
