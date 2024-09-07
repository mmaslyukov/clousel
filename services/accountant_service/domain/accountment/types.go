package accountment

import "github.com/google/uuid"

type PriceTag struct {
	Price       uint
	Rides       uint
	PaymentLink string
}
type PriceTagsDetails struct {
	CarouselId uuid.UUID
	Tags       []PriceTag
}

func PriceTagsDetailsCreate(carouselId uuid.UUID, tags []PriceTag) PriceTagsDetails {
	return PriceTagsDetails{
		CarouselId: carouselId, Tags: tags,
	}
}

type ReceiptDetails struct {
	CarouselId uuid.UUID
	Price      uint
	Rides      uint
	Time       string
	Token      string
}

func ReceiptDetailsCreate(carouselId uuid.UUID, price uint, rides uint, time string, token string) ReceiptDetails {
	return ReceiptDetails{
		CarouselId: carouselId,
		Price:      price,
		Rides:      rides,
		Time:       time,
		Token:      token,
	}
}
