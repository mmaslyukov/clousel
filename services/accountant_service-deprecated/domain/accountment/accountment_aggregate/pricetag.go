package accountment_aggregate

import (
	"accountant_service/domain/accountment"
	"accountant_service/domain/accountment/accountment_error"
	"accountant_service/framework/core"
)

const (
	priceTagsLengthMin = 1
	priceTagsLengthMax = 4
)

type PriceTags struct {
	// TODO: Change to internal data structure and make To/From functions. see ride.go
	tagDetails accountment.PriceTagsDetails
}

func PriceTagsCreate(tagDetails accountment.PriceTagsDetails) PriceTags {
	return PriceTags{tagDetails: tagDetails}
}

func PriceTagsDefaultCreate() PriceTags {
	return PriceTags{tagDetails: accountment.PriceTagsDetails{
		Tags: make([]accountment.PriceTag, 0),
	}}
}

func (p *PriceTags) TagDetails() accountment.PriceTagsDetails {
	return p.tagDetails
}

func (s *PriceTags) ApplyPriceTags(tagDetails accountment.PriceTagsDetails) error {
	var err error
	if len(tagDetails.Tags) < priceTagsLengthMin || len(tagDetails.Tags) > priceTagsLengthMax {
		err = core.ErrorCreate[accountment_error.ErrorInvalidPriceTag]().Message(
			"Price tags length should be more then %d and less than %d",
			priceTagsLengthMin, priceTagsLengthMax)
	}
	return err
}
