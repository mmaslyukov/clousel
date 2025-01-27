package accountment_consumer

import (
	"accountant_service/domain/accountment/accountment_aggregate"

	"github.com/google/uuid"
)

type IPortRepositorySales interface {
	// ReadAvailableRides(id CarouselId) (int, error)
	// LoadReceipts(carouselId uuid.UUID) ([]Receipt, error)
	SaveReceipt(receipt accountment_aggregate.Receipt) error
	LoadPriceTags(carouselId uuid.UUID) (accountment_aggregate.PriceTags, error)
	WritePriceTags(priceTags accountment_aggregate.PriceTags) error
}
