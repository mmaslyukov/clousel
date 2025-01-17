package client

import (
	"clousel/lib/fault"
	"time"

	"github.com/google/uuid"
)

type IClientRepoGeneralAdapter interface {
	// TODO validate email
	SaveNewClientEntry(userId uuid.UUID, username string, email string, password string, companyName string) fault.IError
	ReadClientEntryByName(username string) (*ClientEntry, fault.IError)
	ReadClientEntryById(userId uuid.UUID) (*ClientEntry, fault.IError)
	// use balance form Balance table
	// UpdateBalance(userId uuid.UUID, value int) fault.IError
}

type IClientRepoCheckoutSessionAdapter interface {
	SaveNewCheckoutEntry(eventId uuid.UUID, userId uuid.UUID, sessionId string, price int, tickets int) fault.IError
	UpdateCheckoutStatus(sessionId string, status PaymentStatus) fault.IError
	ReadCheckoutEntriesBySessionId(sessionId string) (*CheckoutEntry, fault.IError)
	ReadCheckoutEntriesByUserId(userId uuid.UUID, begin *time.Time, end *time.Time) ([]*CheckoutEntry, fault.IError)
	ReadCheckoutEntriesAll(begin *time.Time, end *time.Time) ([]*CheckoutEntry, fault.IError)
}

type IClientRepoBalanceChangeAdapter interface {
	SaveNewBalanceChangeEntry(eventId uuid.UUID, userId uuid.UUID, tickets int) fault.IError
	ReadBalanceEntriesByUserId(userId uuid.UUID) ([]*TicketsBalanceEntry, fault.IError)
	// Run go routing once per 24h to check balance (using ReadBalanceEntriesByUserId)
	// and remove one by one (using RemoveBalanceByEventId) if total is 0
	RemoveBalanceByEventId(eventId uuid.UUID) fault.IError
	// deprecated
	// CalculateBalanceTotal(userId uuid.UUID) (int, fault.IError)
	// ClearBalanceByUserId(userId uuid.UUID) fault.IError
}
