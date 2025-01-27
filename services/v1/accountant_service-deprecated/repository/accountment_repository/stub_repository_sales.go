package accountment_repository

import (
	"accountant_service/domain/accountment"
	"accountant_service/domain/accountment/accountment_aggregate"
	"accountant_service/domain/accountment/accountment_consumer"

	"github.com/google/uuid"
)

type StubSalesRepository struct {
}

func StubSalesRepositoryCreate() accountment_consumer.IPortRepositorySales {
	return &StubSalesRepository{}
}

func (s *StubSalesRepository) WritePriceTags(priceTags accountment_aggregate.PriceTags) error {
	return nil
}

func (s *StubSalesRepository) LoadPriceTags(id uuid.UUID) (accountment_aggregate.PriceTags, error) {
	tagsArray := make([]accountment.PriceTag, 0)
	tagsArray = append(tagsArray, accountment.PriceTag{Price: 1, Rides: 1})
	tagsDetails := accountment.PriceTagsDetailsCreate(uuid.New(), tagsArray)
	tags := accountment_aggregate.PriceTagsCreate(tagsDetails)
	return tags, nil
}

func (s *StubSalesRepository) SaveReceipt(receipt accountment_aggregate.Receipt) error {
	return nil
}
