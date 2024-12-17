package accountment_service

import (
	"accountant_service/domain/accountment"
	"accountant_service/domain/accountment/accountment_consumer"
	"accountant_service/framework/core"

	"github.com/google/uuid"
)

type ServiceAnalytics struct {
	repository accountment_consumer.IPortRepositoryAnalytics
	logger     core.ILogger
}

func ServiceAnalyticsCreate(logger core.ILogger, repository accountment_consumer.IPortRepositoryAnalytics) IServiceAnalytics {
	return &ServiceAnalytics{
		logger:     logger,
		repository: repository}
}

func (s *ServiceAnalytics) LoadReceipts(carouselId uuid.UUID) ([]accountment.ReceiptDetails, error) {
	/// No aggregates here necessary
	return s.repository.LoadReceipts(carouselId)
}
