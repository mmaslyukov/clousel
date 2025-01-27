package carousel_consumer

import "accountant_service/domain/carousel"

type IPortPublisherGateway interface {
	PublishRefill(msg carousel.RideMinimal) (carousel.ResponseRefillRides, error)
	// CheckAndPublishIfAvailableRides(carouselId uuid.UUID) (int, error)
}
