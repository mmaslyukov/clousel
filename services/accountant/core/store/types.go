package store

import "github.com/google/uuid"

const PaymentStatusPaid = "paid"

type Session = string
type Product = string
type Owner = uuid.UUID
type Carousel = uuid.UUID

const (
	BookOrderStatusNew           = "new"
	BookOrderStatusPaid          = "paid"
	BookOrderStatusFailed        = "failed"
	BookOrderStatusPendingRefill = "pending_refill"
	BookOrderStatusRefilled      = "refilled"
	BookOrderStatusCanceled      = "canceled"
)

type BookEntry struct {
	SessionId Session
	Time      string
	CarId     Carousel
	Status    string
	Error     *string
	Amount    int
	Tickets   int
}

func BookEntryCreate(sessionId Session, carId Carousel, amount int, tickets int) BookEntry {
	return BookEntry{
		SessionId: sessionId,
		Time:      "",
		CarId:     carId,
		Amount:    amount,
		Tickets:   tickets,
		Status:    BookOrderStatusNew,
	}
}

type PaymentResltUrls struct {
	Success string
	Cancel  string
}

type PriceTag struct {
	Amount  int
	Tickets int
	PriceId string
}
