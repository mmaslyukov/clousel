package accountment_aggregate

// import (
// 	"accountant_service/domain/accountment"
// 	"accountant_service/domain/accountment/accountment_error"

// 	"github.com/google/uuid"
// )

// const (
// 	priceTagsLengthMin = 1
// 	priceTagsLengthMax = 4
// )

// // type PriceTag struct {
// // 	price uint
// // 	rides uint
// // }

// // func (p *PriceTag) Price() uint {
// // 	return p.price
// // }
// // func (p *PriceTag) Rides() uint {
// // 	return p.rides
// // }

// // type Receipt struct {
// // 	carouselId uuid.UUID
// // 	price      uint
// // 	rides      uint
// // 	time       string
// // 	token      string
// // }

// //	func (r *Receipt) CarouselId() uuid.UUID {
// //		return r.carouselId
// //	}
// //
// //	func (r *Receipt) Price() uint {
// //		return r.price
// //	}
// //
// //	func (r *Receipt) Rides() uint {
// //		return r.rides
// //	}
// //
// //	func (r *Receipt) Time() string {
// //		return r.time
// //	}
// //
// //	func (r *Receipt) Token() string {
// //		return r.token
// //	}
// //
// // type Receipt = accountment.Receipt
// type PriceTags = accountment.PriceTags

// type Sales struct {
// 	carouselId uuid.UUID
// 	// receipts   []Receipt
// 	priceTags PriceTags
// 	// logger core.ILogger
// 	// priceTags  accountment.PriceTag
// }

// func NewSales(id uuid.UUID) Sales {
// 	return Sales{
// 		carouselId: id,
// 		// receipts:   make([]Receipt, 0),
// 		// priceTags:
// 	}
// }

// // func (s *Sales) WritePriceTags(id accountment.CarouselId, tags []accountment.PriceTag) error {
// // 	return s.repository.WritePriceTags(id, tags)
// // }

// //	func (s *Sales) ReadAvailableRides(id accountment.CarouselId) (int, error) {
// //		s.logger.Dbg().Printf("Reading available rides")
// //		return s.repository.ReadAvailableRides(id)
// //	}
// //
// //	func (s *Sales) ReadPriceTags(id accountment.CarouselId) ([]accountment.PriceTag, error) {
// //		return s.repository.ReadPriceTags(id)
// //	}
// func (s *Sales) ApplyPriceTags(tags accountment.PriceTags) error {
// 	var err error
// 	if len(tags) < priceTagsLengthMin || len(tags) > priceTagsLengthMax {
// 		err = accountment_error.NewErrorInvalidPriceTag().Message(
// 			"Price tags length should be more then %d and less than %d",
// 			priceTagsLengthMin, priceTagsLengthMax)
// 	}
// 	s.priceTags = tags
// 	// s.priceTags = make([]PriceTag, 0)
// 	// for _, t := range tags {
// 	// 	s.priceTags = append(s.priceTags, PriceTag{
// 	// 		price: t.Price(),
// 	// 		rides: t.Rides(),
// 	// 	})
// 	// }
// 	return err
// }

// // func (s *Sales) ApplyReceipt(receipt accountment.Receipt) (core.IEvent, error) {
// // 	s.receipts = append(s.receipts, receipt)
// // 	// s.receipt = utils.NewOptionalValue[Receipt](Receipt{
// // 	// 	carouselId: receipt.CarouselId(),
// // 	// 	time:       receipt.Time(),
// // 	// 	token:      receipt.Token(),
// // 	// 	price:      receipt.Price(),
// // 	// 	rides:      receipt.Rides(),
// // 	// })
// // 	if receipt.Rides > 0 {
// // 		return accountment_event.NewEventRidesUpdated(receipt.Rides), nil
// // 	} else {
// // 		return nil, accountment_error.NewErrorInvalidRidesNumber().Message(
// // 			"Rides in the ticket supposed to be more than 0")
// // 	}
// // }
