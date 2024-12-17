package accountment_consumer

import (
	"accountant_service/domain/accountment"

	"github.com/google/uuid"
)

type IPortRepositoryAnalytics interface {
	LoadReceipts(carouselId uuid.UUID) ([]accountment.ReceiptDetails, error)
}
