package carousel_service

import (
	"accountant_service/domain/accountment/accountment_event"
	"accountant_service/domain/carousel"
	"accountant_service/domain/carousel/carousel_aggregate"
	"accountant_service/domain/carousel/carousel_consumer"
	"accountant_service/domain/carousel/carousel_error"
	"accountant_service/framework/core"

	"github.com/google/uuid"
)

type ServiceRide struct {
	gw         carousel_consumer.IPortPublisherGateway
	repository carousel_consumer.IPortRepositoryRide
	logger     core.ILogger
}

func ServiceRideCreate(
	gw carousel_consumer.IPortPublisherGateway,
	repository carousel_consumer.IPortRepositoryRide,
	logger core.ILogger,
) IServiceRide {
	return &ServiceRide{
		gw:         gw,
		repository: repository,
		logger:     logger,
	}
}

func (s *ServiceRide) HandleEvent(event core.IEvent) error {
	var err error
	switch e := event.(type) {
	case *accountment_event.EventRidesUpdated:
		s.logger.Inf().Printf(
			"Event is arrived: %s(%s): carouselId: %s, rides:%d",
			e.Name(), e.Id(), e.CarouselId(), e.Rides())

		ride := carousel_aggregate.RideCreate(e.CarouselId(), e.Rides())
		if err = s.repository.Save(ride); err != nil {
			return err
		}
		var responseRefill carousel.ResponseRefillRides
		if responseRefill, err = s.gw.PublishRefill(ride.ToMinimal()); err != nil {
			return err
		}
		if !responseRefill.Success {
			return core.ErrorCreate[carousel_error.ErrorRefill]().Message(responseRefill.Error)
		}
		ride.SetDelivered()
		s.logger.Inf().Printf("Set Delivered flag and save the ride for CarouselId: %s", ride.CarouselId())
		if err = s.repository.Save(ride); err != nil {
			return err
		}

	}
	return err
}

func (s *ServiceRide) CheckAndPublishUndeliveredRides(
	carouselId uuid.UUID,
) ([]carousel.ResponseRefillRides, error) {
	var err error
	var rides []carousel_aggregate.Ride
	responseRides := make([]carousel.ResponseRefillRides, 0)

	if rides, err = s.repository.LoadUndelivered(carouselId); err != nil {
		return nil, err
	}

	for _, ride := range rides {
		var responseRefill carousel.ResponseRefillRides
		if responseRefill, err = s.gw.PublishRefill(ride.ToMinimal()); err != nil {
			return nil, err
		}
		if responseRefill.Success {
			ride.SetDelivered()
			s.logger.Inf().Printf("Set Delivered flag and save the ride for CarouselId: %s", ride.CarouselId())
			if err = s.repository.Save(ride); err != nil {
				s.logger.Err().Printf("Couldn't save the ride: %s", err)
			}
		}
		responseRides = append(responseRides, responseRefill)
	}

	return responseRides, nil
}

func (s *ServiceRide) ReadUndeliveredRides(
	carouselId uuid.UUID,
) ([]carousel.RideMinimal, error) {
	var err error
	var rides []carousel_aggregate.Ride

	if rides, err = s.repository.LoadUndelivered(carouselId); err != nil {
		return nil, err
	}

	availableRides := make([]carousel.RideMinimal, 0)
	for _, ride := range rides {
		availableRides = append(availableRides, ride.ToMinimal())
	}
	return availableRides, nil
}
