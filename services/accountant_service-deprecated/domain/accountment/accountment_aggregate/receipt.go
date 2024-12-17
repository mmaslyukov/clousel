package accountment_aggregate

import (
	"accountant_service/domain/accountment"
	"accountant_service/domain/accountment/accountment_error"
	"accountant_service/domain/accountment/accountment_event"
	"accountant_service/framework/core"
)

type Receipt struct {
	// TODO: Change to internal data structure and make To/From functions. see ride.go
	receiptDetails accountment.ReceiptDetails
}

func ReceiptCreate(receiptDetails accountment.ReceiptDetails) Receipt {
	return Receipt{
		receiptDetails: receiptDetails,
	}
}
func ReceiptDefaultCreate() Receipt {
	return Receipt{}
}

func (r *Receipt) ReceiptDetails(receiptDetails accountment.ReceiptDetails) accountment.ReceiptDetails {
	return r.receiptDetails
}
func (r *Receipt) ApplyReceipt(receiptDetails accountment.ReceiptDetails) (core.IEvent, error) {
	r.receiptDetails = receiptDetails
	if receiptDetails.Rides > 0 {
		return accountment_event.EventRidesUpdatedCreate(receiptDetails.Rides, receiptDetails.CarouselId), nil
	} else {
		return nil, core.ErrorCreate[accountment_error.ErrorInvalidRidesNumber]().Message(
			"Rides in the ticket supposed to be more than 0")
	}
}
