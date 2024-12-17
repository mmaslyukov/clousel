package accountment_service

import (
	"accountant_service/domain/accountment"
	"accountant_service/domain/accountment/accountment_aggregate"
	"accountant_service/domain/accountment/accountment_consumer"
	"accountant_service/framework/core"

	"github.com/google/uuid"
)

type ServiceSales struct {
	core.EventSubscribable
	repository accountment_consumer.IPortRepositorySales
	logger     core.ILogger
}

// func NewServiceSales(sales *accountment_context.Sales) IServiceSales {
func ServiceSalesCreate(logger core.ILogger, repository accountment_consumer.IPortRepositorySales) IServiceSales {
	return &ServiceSales{
		EventSubscribable: *core.EventSubscribableCreate(),
		logger:            logger,
		repository:        repository}
}

// func (s *ServiceSales) ReadAvailableRides(id accountment.CarouselId) (int, error) {
// 	return s.repository.ReadAvailableRides(id)
// }

func (s *ServiceSales) WritePriceTags(priceTagsDetails accountment.PriceTagsDetails) error {
	priceTags := accountment_aggregate.PriceTagsDefaultCreate()
	err := priceTags.ApplyPriceTags(priceTagsDetails)
	if err != nil {
		return err
	}
	return s.repository.WritePriceTags(priceTags)
}

func (s *ServiceSales) ReadPriceTags(carouselId uuid.UUID) (accountment.PriceTagsDetails, error) {
	priceTags, err := s.repository.LoadPriceTags(carouselId)
	return priceTags.TagDetails(), err
}

func (s *ServiceSales) ApplyAndSaveReceipt(receiptDetails accountment.ReceiptDetails) error {
	receipt := accountment_aggregate.ReceiptDefaultCreate()
	event, err := receipt.ApplyReceipt(receiptDetails)
	if err != nil {
		return err
	}
	err = s.repository.SaveReceipt(receipt)
	if err != nil {
		return err
	}
	s.Notify(event)
	return nil
}
