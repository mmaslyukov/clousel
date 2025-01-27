package accountment_repository

import (
	"accountant_service/domain/accountment"
	"accountant_service/domain/accountment/accountment_consumer"

	"github.com/google/uuid"
)

type StubAnalyticsRepository struct {
}

func StubAnalyticsRepositoryCreate() accountment_consumer.IPortRepositoryAnalytics {
	return &StubAnalyticsRepository{}
}

func (s *StubAnalyticsRepository) LoadReceipts(carouselId uuid.UUID) ([]accountment.ReceiptDetails, error) {
	receipDetails := accountment.ReceiptDetailsCreate(uuid.New(), 1, 1, "<time>", "<token>")
	receipDetailsArray := make([]accountment.ReceiptDetails, 0)
	receipDetailsArray = append(receipDetailsArray, receipDetails)
	return receipDetailsArray, nil
}
