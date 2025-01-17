package client

import (
	"time"

	"github.com/google/uuid"
)

type ClientEntry struct {
	UserId           uuid.UUID
	Username         string
	CompanyName      string
	Email            string
	Password         string
	Balance          int
	RegistrationTime time.Time
}

func (b *ClientEntry) Id() uuid.UUID {
	return b.UserId
}

type PaymentStatus = string

const (
	PaymentStatusNew       PaymentStatus = "created"
	PaymentStatusPaid      PaymentStatus = "fulfilled"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

type CheckoutEntry struct {
	EventId     uuid.UUID //PK
	SessionId   string
	UserId      uuid.UUID //FK
	Price       int
	Tickets     int
	PaymentTime time.Time
	Status      PaymentStatus
}

type PriceTag struct {
	Amount  int    `json:"Amount"`
	Tickets int    `json:"Tickets"`
	PriceId string `json:"PriceId"`
}

type PaymentResultUrls struct {
	Success string
	Cancel  string
}

type TicketsBalanceEntry struct {
	EventId uuid.UUID // PK. Event ID is one out of two CheckoutSessionEntry.EventId or GameEntry.EventId
	UserId  uuid.UUID // FK
	Change  int
}
