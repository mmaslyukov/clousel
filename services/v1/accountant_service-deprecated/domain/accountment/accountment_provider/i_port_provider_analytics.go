package accountment_provider

import (
	"accountant_service/domain/accountment"

	"github.com/google/uuid"
)

type IPortApiAnalytics interface {
	LoadReceipts(carouselId uuid.UUID) ([]accountment.ReceiptDetails, error)
	// LoadReceiptsOwner(owner uuid.UUID) ([]accountment.ReceiptDetails, error)
}
