package accountment_provider

import (
	"accountant_service/domain/accountment"

	"github.com/google/uuid"
)

type IPortApiSales interface {
	// ReadAvailableRides(id CarouselId) (int, error)
	WritePriceTags(priceTagsDetails accountment.PriceTagsDetails) error
	ReadPriceTags(carouselId uuid.UUID) (accountment.PriceTagsDetails, error)
	ApplyAndSaveReceipt(receiptDetails accountment.ReceiptDetails) error
	// ApplyAndSaveReceipt(id CarouselId, receipt Receipt) (core.IEvent, error)
}
