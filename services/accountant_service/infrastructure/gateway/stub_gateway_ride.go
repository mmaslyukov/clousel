package gateway

import (
	"accountant_service/domain/carousel"
	"accountant_service/domain/carousel/carousel_consumer"
)

type StubPublisherGatewat struct {
}

func StubPublisherGatewatCreate() carousel_consumer.IPortPublisherGateway {
	return &StubPublisherGatewat{}
}

func (s *StubPublisherGatewat) PublishRefill(
	msg carousel.RideMinimal,
) (carousel.ResponseRefillRides, error) {
	return carousel.ResponseRefillRides{
		RideMinimal: msg,
		Success:     true,
		Error:       "",
	}, nil
}
