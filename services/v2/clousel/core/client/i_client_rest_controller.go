package client

import (
	"clousel/lib/fault"
	"time"

	"github.com/google/uuid"
)

type IClientSessionSelector interface {
	UserId() *uuid.UUID
	Status() *PaymentStatus
	TimeFrom() *time.Time
	TimeTill() *time.Time
}

type IClientRestController interface {
	Register(username string, email string, password string, companyName string) fault.IError
	Login(username string, password string) (*ClientEntry, fault.IError)
	BuyTickets(userId uuid.UUID, priceId string, afterSellVisitUrl string) (ISession, fault.IError)
	ReadPriceOptions(userId uuid.UUID) ([]PriceTag, fault.IError)
	ReadBalance(userId uuid.UUID) (int, fault.IError)
	// ReadPriceOptions(companyName string) ([]PriceTag, fault.IError)
	ApplyPaymentResults(sessionId string, status PaymentStatus) fault.IError
	// ReadSessionListBySelector(selector IClientSessionSelector) ([]CheckoutSessionEntry, fault.IError)
}
