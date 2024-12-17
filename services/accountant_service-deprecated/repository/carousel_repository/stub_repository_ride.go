package carousel_repository

import (
	"accountant_service/domain/carousel/carousel_aggregate"
	"accountant_service/domain/carousel/carousel_consumer"

	"github.com/google/uuid"
)

type StubRideRepository struct {
}

func StubRideRepositoryCreate() carousel_consumer.IPortRepositoryRide {
	return &StubRideRepository{}
}

func (s *StubRideRepository) Save(ride carousel_aggregate.Ride) error {
	return nil
}

func (s *StubRideRepository) LoadUndelivered(caroselId uuid.UUID) ([]carousel_aggregate.Ride, error) {
	return make([]carousel_aggregate.Ride, 0), nil
}

func (s *StubRideRepository) LoadDelivered(caroselId uuid.UUID) ([]carousel_aggregate.Ride, error) {
	return make([]carousel_aggregate.Ride, 0), nil

}
