package repouser

import (
	"clousel/core/client"
	"clousel/infra/repo/driver"
	rec "clousel/infra/repo/errors"
	"clousel/infra/repo/types"
	"clousel/lib/fault"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	TableUser     = "user"
	TableBalance  = "balance"
	TableCheckout = "checkout"
)

/*
 */
type UserColumn struct {
	UserId      types.Named[types.UUIDString]
	CompanyName types.Named[string]
	Username    types.Named[string]
	Email       types.Named[string]
	Password    types.Named[string]
	// Balance          types.Named[int]
	RegistrationTime types.Named[types.TimeString]
}

func UserColumnDefault() UserColumn {
	prof := UserColumn{
		UserId:      types.NamedCreateDefault[types.UUIDString]("Id"),
		CompanyName: types.NamedCreateDefault[string]("Companyname"),
		Username:    types.NamedCreateDefault[string]("Username"),
		Email:       types.NamedCreateDefault[string]("Email"),
		Password:    types.NamedCreateDefault[string]("Password"),
		// Balance:          types.NamedCreateDefault[int]("Balance"),
		RegistrationTime: types.NamedCreateDefault[types.TimeString]("Time"),
	}
	return prof
}

func (c *UserColumn) toEntry() *client.ClientEntry {
	entry := &client.ClientEntry{
		UserId:           c.UserId.ValuePtr().Uuid(),
		Username:         c.Username.Value(),
		Email:            c.Email.Value(),
		Password:         c.Password.Value(),
		CompanyName:      c.CompanyName.Value(),
		RegistrationTime: c.RegistrationTime.ValuePtr().Time(),
	}
	return entry
}

/*
 */
type CheckoutColumn struct {
	EventId   types.Named[types.UUIDString]
	SessionId types.Named[string]
	UserId    types.Named[types.UUIDString]
	Price     types.Named[int]
	Tickets   types.Named[int]
	Status    types.Named[string]
	Time      types.Named[types.TimeString]
}

func CheckoutColumnDefault() CheckoutColumn {
	prof := CheckoutColumn{
		EventId:   types.NamedCreateDefault[types.UUIDString]("EventId"),
		SessionId: types.NamedCreateDefault[string]("SessionId"),
		UserId:    types.NamedCreateDefault[types.UUIDString]("UserId"),
		Price:     types.NamedCreateDefault[int]("Price"),
		Tickets:   types.NamedCreateDefault[int]("Tickets"),
		Status:    types.NamedCreateDefault[string]("Status"),
		Time:      types.NamedCreateDefault[types.TimeString]("Time"),
	}
	return prof
}

func (c *CheckoutColumn) toEntry() *client.CheckoutEntry {
	entry := &client.CheckoutEntry{
		EventId:     c.EventId.ValuePtr().Uuid(),
		SessionId:   c.SessionId.Value(),
		UserId:      c.UserId.ValuePtr().Uuid(),
		Price:       c.Price.Value(),
		Tickets:     c.Tickets.Value(),
		Status:      c.Status.Value(),
		PaymentTime: c.Time.ValuePtr().Time(),
	}
	return entry
}

/*
 */
type TicketBalanceColumn struct {
	EventId types.Named[types.UUIDString]
	UserId  types.Named[types.UUIDString]
	Change  types.Named[int]
}

func TicketBalanceColumnDefault() TicketBalanceColumn {
	prof := TicketBalanceColumn{
		EventId: types.NamedCreateDefault[types.UUIDString]("EventId"),
		UserId:  types.NamedCreateDefault[types.UUIDString]("UserId"),
		Change:  types.NamedCreateDefault[int]("Change"),
	}
	return prof
}

func (c *TicketBalanceColumn) toEntry() *client.TicketsBalanceEntry {
	entry := &client.TicketsBalanceEntry{
		EventId: c.EventId.ValuePtr().Uuid(),
		UserId:  c.UserId.ValuePtr().Uuid(),
		Change:  c.Change.Value(),
	}
	return entry
}

/*
Repository User
*/

type RepositoryUser struct {
	drv driver.IDBDriver
	log *zerolog.Logger
}

func RepositoryUserCreate(drv driver.IDBDriver, log *zerolog.Logger) *RepositoryUser {
	return &RepositoryUser{drv: drv, log: log}
}

/*
IClientRepoGeneralAdapter
*/
func (r *RepositoryUser) SaveNewClientEntry(userId uuid.UUID, username string, email string, password string, companyName string) fault.IError {
	const fn = "Repository.User.SaveNewClientEntry"
	var err fault.IError
	t := UserColumnDefault()
	t.UserId.ValuePtr().SetUuid(userId)
	t.Email.SetValue(email)
	t.Username.SetValue(username)
	t.Password.SetValue(password)
	t.CompanyName.SetValue(companyName)
	t.RegistrationTime.ValuePtr().SetTime(time.Now())
	prompt := fmt.Sprintf("insert into '%s' (%s,%s,%s,%s,%s,%s) values ('%s','%s','%s','%s','%s','%s')", TableUser,
		t.UserId.Name(), t.CompanyName.Name(), t.Username.Name(), t.Email.Name(), t.Password.Name(), t.RegistrationTime.Name(),
		t.UserId.ValuePtr().Str(), t.CompanyName.Value(), t.Username.Value(), t.Email.Value(), t.Password.Value(), t.RegistrationTime.ValuePtr().Str())
	// prompt := fmt.Sprintf("insert into '%s' (%s,%s,%s,%s,%s) values ('%s','%s','%s','%s',%d)", TableUser,
	// 	t.UserId.Name(), t.Username.Name(), t.Email.Name(), t.Password.Name(), t.Balance.Name(),
	// 	t.UserId.ValuePtr().Str(), t.Username.Value(), t.Email.Value(), t.Password.Value(), 0)
	if e := r.drv.Session(func(db *sql.DB) error {
		var err error
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Send()
		} else {
			r.log.Err(err).Str("SQL", prompt).Msgf("%s: Fail to Add", fn)
		}
		return err
	}); e != nil {
		err = fault.New(rec.ERepoExecPrompt).Msg(e.Error())
	}
	return err
}

func (r *RepositoryUser) ReadClientEntryByName(username string) (*client.ClientEntry, fault.IError) {
	t := UserColumnDefault()
	t.Username.SetValue(username)
	filter := fmt.Sprintf("%s='%s'", t.Username.Name(), t.Username.Value())
	entries, err := r.readUserEntriesBy(func() string {
		return filter
	})
	return takeFirst(entries, err)
}

func (r *RepositoryUser) ReadClientEntryById(userId uuid.UUID) (*client.ClientEntry, fault.IError) {
	t := UserColumnDefault()
	t.UserId.ValuePtr().SetUuid(userId)
	filter := fmt.Sprintf("%s='%s'", t.UserId.Name(), t.UserId.ValuePtr().Str())
	entries, err := r.readUserEntriesBy(func() string {
		return filter
	})
	return takeFirst(entries, err)
}

// func (r *RepositoryUser) UpdateBalance(userId uuid.UUID, value int) fault.IError {
// 	const fn = "Repository.User.UpdateBalance"
// 	var err fault.IError
// 	t := UserColumnDefault()
// 	t.UserId.ValuePtr().SetUuid(userId)
// 	t.Balance.SetValue(value)
// 	prompt := fmt.Sprintf("update '%s' set %s=%d where %s='%s'", TableUser,
// 		t.Balance.Name(), t.Balance.Value(),
// 		t.UserId.Name(), t.UserId.ValuePtr().Str())

// 	if e := r.drv.Session(func(db *sql.DB) error {
// 		var e error
// 		if _, e = db.Exec(prompt); e == nil {
// 			r.log.Debug().Str("SQL", prompt).Msgf("%s: Success", fn)
// 		} else {
// 			r.log.Err(e).Str("SQL", prompt).Msgf("%s: Failure", fn)
// 		}
// 		return e
// 	}); e != nil {
// 		err = fault.New(rec.ERepoExecPrompt).Msgf("Fail to update user balance: %s", e.Error())
// 	}
// 	return err
// }

/*
IClientRepoCheckoutAdapter
*/
func (r *RepositoryUser) SaveNewCheckoutEntry(eventId uuid.UUID, userId uuid.UUID, sessionId string, price int, tickets int) fault.IError {
	const fn = "Repository.User.SaveNewCheckoutEntry"
	var err fault.IError
	t := CheckoutColumnDefault()
	t.EventId.ValuePtr().SetUuid(eventId)
	t.SessionId.SetValue(sessionId)
	t.UserId.ValuePtr().SetUuid(userId)
	t.Price.SetValue(price)
	t.Tickets.SetValue(tickets)
	t.Status.SetValue(client.PaymentStatusNew)
	t.Time.ValuePtr().SetTime(time.Now())
	prompt := fmt.Sprintf("insert into '%s' (%s,%s,%s,%s,%s,%s,%s) values ('%s','%s','%s', %d, %d,'%s','%s')", TableCheckout,
		t.EventId.Name(), t.SessionId.Name(), t.UserId.Name(), t.Price.Name(), t.Tickets.Name(), t.Status.Name(), t.Time.Name(),
		t.EventId.ValuePtr().Str(), t.SessionId.Value(), t.UserId.ValuePtr().Str(), t.Price.Value(), t.Tickets.Value(), t.Status.Value(), t.Time.ValuePtr().Str())
	if e := r.drv.Session(func(db *sql.DB) error {
		var err error
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Send()
		} else {
			r.log.Err(err).Str("SQL", prompt).Msgf("%s: Fail to Add", fn)
		}
		return err
	}); e != nil {
		err = fault.New(rec.ERepoExecPrompt).Msg(e.Error())
	}
	return err
	// panic(fmt.Sprintf("%s not implemented", fn))
}

func (r *RepositoryUser) UpdateCheckoutStatus(sessionId string, status client.PaymentStatus) fault.IError {
	const fn = "Repository.User.UpdateCheckoutStatus"
	var err fault.IError
	t := CheckoutColumnDefault()
	t.SessionId.SetValue(sessionId)
	t.Status.SetValue(status)
	t.Time.ValuePtr().SetTime(time.Now())
	prompt := fmt.Sprintf("update '%s' set %s='%s', %s='%s' where %s='%s'", TableCheckout,
		t.Time.Name(), t.Time.ValuePtr().Str(),
		t.Status.Name(), t.Status.Value(),
		t.SessionId.Name(), t.SessionId.Value())

	if e := r.drv.Session(func(db *sql.DB) error {
		var e error
		if _, e = db.Exec(prompt); e == nil {
			r.log.Debug().Str("SQL", prompt).Msgf("%s: Success", fn)
		} else {
			r.log.Err(e).Str("SQL", prompt).Msgf("%s: Failure", fn)
		}
		return e
	}); e != nil {
		err = fault.New(rec.ERepoExecPrompt).Msgf("Fail to update checkout status %s", e.Error())
	}
	return err
	// panic(fmt.Sprintf("%s not implemented", fn))
}

func (r *RepositoryUser) ReadCheckoutEntriesByUserId(userId uuid.UUID, begin *time.Time, end *time.Time) ([]*client.CheckoutEntry, fault.IError) {
	const fn = "Repository.User.ReadCheckoutEntriesByUserId"
	t := TicketBalanceColumnDefault()
	t.UserId.ValuePtr().SetUuid(userId)
	filter := fmt.Sprintf("%s='%s'", t.UserId.Name(), t.UserId.ValuePtr().Str())
	rows, err := r.readCheckoutEntriesBy(func() string {
		return filter
	})
	var entries []*client.CheckoutEntry
	if err == nil {
		for _, r := range rows {
			entries = append(entries, r.toEntry())
		}
	}
	return entries, err
	// panic(fmt.Sprintf("%s not implemented", fn))
}

func (r *RepositoryUser) ReadCheckoutEntriesBySessionId(sessionId string) (entry *client.CheckoutEntry, err fault.IError) {
	const fn = "Repository.User.ReadCheckoutEntriesByUserId"
	var rows []*CheckoutColumn
	t := CheckoutColumnDefault()
	t.SessionId.SetValue(sessionId)
	filter := fmt.Sprintf("%s='%s'", t.SessionId.Name(), t.SessionId.Value())
	rows, err = r.readCheckoutEntriesBy(func() string {
		return filter
	})
	return takeFirst(rows, err)
}

func (r *RepositoryUser) ReadCheckoutEntriesAll(begin *time.Time, end *time.Time) ([]*client.CheckoutEntry, fault.IError) {
	const fn = "Repository.User.ReadCheckoutEntriesAll"
	panic(fmt.Sprintf("%s not implemented", fn))

}

/*
IClientRepoBalanceChangeAdapter
*/
func (r *RepositoryUser) SaveNewBalanceChangeEntry(eventId uuid.UUID, userId uuid.UUID, tickets int) fault.IError {
	const fn = "Repository.User.SaveNewBalanceChangeEntry"
	var err fault.IError
	t := TicketBalanceColumnDefault()
	t.EventId.ValuePtr().SetUuid(eventId)
	t.UserId.ValuePtr().SetUuid(userId)
	t.Change.SetValue(tickets)
	prompt := fmt.Sprintf("insert into '%s' (%s,%s,%s) values ('%s','%s',%d)", TableBalance,
		t.EventId.Name(), t.UserId.Name(), t.Change.Name(),
		t.EventId.ValuePtr().Str(), t.UserId.ValuePtr().Str(), t.Change.Value())
	if e := r.drv.Session(func(db *sql.DB) error {
		var err error
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Send()
		} else {
			r.log.Err(err).Str("SQL", prompt).Msgf("%s: Fail to Add", fn)
		}
		return err
	}); e != nil {
		err = fault.New(rec.ERepoExecPrompt).Msg(e.Error())
	}
	return err
}

func (r *RepositoryUser) ReadBalanceEntriesByUserId(userId uuid.UUID) ([]*client.TicketsBalanceEntry, fault.IError) {
	const fn = "Repository.User.ReadBalanceEntriesByUserId"
	t := TicketBalanceColumnDefault()
	t.UserId.ValuePtr().SetUuid(userId)
	filter := fmt.Sprintf("%s='%s'", t.UserId.Name(), t.UserId.ValuePtr().Str())
	rows, err := r.readBalanceEntriesBy(func() string {
		return filter
	})
	var entries []*client.TicketsBalanceEntry
	if err == nil {
		for _, r := range rows {
			entries = append(entries, r.toEntry())
		}
	}

	return entries, err
}

func (r *RepositoryUser) RemoveBalanceByEventId(eventId uuid.UUID) (err fault.IError) {
	const fn = "Repository.User.ClearBalanceByUserId"
	t := TicketBalanceColumnDefault()
	t.EventId.ValuePtr().SetUuid(eventId)
	prompt := fmt.Sprintf("delete from '%s' where %s='%s'", TableBalance,
		t.EventId.Name(), t.EventId.ValuePtr().Str())
	if e := r.drv.Session(func(db *sql.DB) error {
		var err error
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Send()
		} else {
			r.log.Err(err).Str("SQL", prompt).Msgf("%s: Fail to delete entries", fn)
		}
		return err
	}); e != nil {
		err = fault.New(rec.ERepoExecPrompt).Msg(e.Error())
	}
	return err
}

//	func (r *RepositoryUser) CalculateBalanceTotal(userId uuid.UUID) (int, fault.IError) {
//		const fn = "Repository.User.CalculateBalanceTotal"
//		var err fault.IError
//		t := TicketBalanceColumnDefault()
//		t.UserId.ValuePtr().SetUuid(userId)
//		filter := fmt.Sprintf("%s='%s'", t.UserId.Name(), t.UserId.ValuePtr().Str())
//		rows, err := r.readBalanceEntriesBy(func() string {
//			return filter
//		})
//		balance := 0
//		if err == nil {
//			for _, r := range rows {
//				balance += r.Change.Value()
//			}
//		}
//		return balance, err
//	}
// func (r *RepositoryUser) ClearBalanceByUserId(userId uuid.UUID) fault.IError {
// 	const fn = "Repository.User.ClearBalanceByUserId"
// 	var err fault.IError
// 	t := TicketBalanceColumnDefault()
// 	t.UserId.ValuePtr().SetUuid(userId)
// 	prompt := fmt.Sprintf("delete from '%s' where %s='%s'", TableBalance,
// 		t.UserId.Name(), t.UserId.ValuePtr().Str())
// 	if e := r.drv.Session(func(db *sql.DB) error {
// 		var err error
// 		if _, err = db.Exec(prompt); err == nil {
// 			r.log.Debug().Str("SQL", prompt).Send()
// 		} else {
// 			r.log.Err(err).Str("SQL", prompt).Msgf("%s: Fail to delete entries", fn)
// 		}
// 		return err
// 	}); e != nil {
// 		err = fault.New(rec.ERepoExecPrompt).Msg(e.Error())
// 	}
// 	return err
// }

/*
Private functions and interfaces
*/
type ToEntry[K any] interface {
	toEntry() *K
}

func takeFirst[T ToEntry[K], K any](arr []T, err fault.IError) (*K, fault.IError) {
	len := len(arr)
	if err != nil {
		return nil, err
	} else if len == 1 {
		return arr[0].toEntry(), nil
	} else {
		return nil, fault.New(rec.ERepoUnexpectedEntriesCount).Msgf("Entries count by filter is unexpected: %d", len)
	}
}

func (r *RepositoryUser) readUserEntriesBy(where func() string) ([]*UserColumn, fault.IError) {
	const fn = "Repository.User.readUserEntryBy"
	var err fault.IError
	var e error
	var entries []*UserColumn
	prompt := fmt.Sprintf("select * from '%s' where %s", TableUser, where())
	if e = r.drv.Session(func(db *sql.DB) error {
		var rows *sql.Rows
		if rows, e = db.Query(prompt); e == nil {
			defer rows.Close()
			for rows.Next() {
				t := UserColumnDefault()
				if e := rows.Scan(
					t.UserId.ValuePtr().Ptr(),
					t.CompanyName.ValuePtr(),
					t.Username.ValuePtr(),
					t.Email.ValuePtr(),
					t.Password.ValuePtr(),
					// t.Balance.ValuePtr(),
					t.RegistrationTime.ValuePtr().Ptr(),
				); e == nil {
					entries = append(entries, &t)
				} else {
					r.log.Err(e).Msgf("%s:", fn)
				}
			}
			r.log.Debug().Str("SQL", prompt).Msgf("%s: Success", fn)
		} else {
			r.log.Err(e).Str("SQL", prompt).Msgf("%s: Fail to Read from '%s' table", fn, TableUser)
		}
		return e
	}); e != nil {
		err = fault.New(rec.ERepoQueryData).Msg(e.Error())
		r.log.Err(e).Str("SQL", prompt).Msgf("%s: Fail to Read from '%s' table", fn, TableUser)
	}
	return entries, err
}

func (r *RepositoryUser) readCheckoutEntriesBy(where func() string) ([]*CheckoutColumn, fault.IError) {
	const fn = "Repository.User.readCheckoutEntriesBy"
	var err fault.IError
	var e error
	var entries []*CheckoutColumn
	prompt := fmt.Sprintf("select * from '%s' where %s", TableCheckout, where())
	if e = r.drv.Session(func(db *sql.DB) error {
		var rows *sql.Rows
		if rows, e = db.Query(prompt); e == nil {
			defer rows.Close()
			for rows.Next() {
				t := CheckoutColumnDefault()
				if e := rows.Scan(
					t.EventId.ValuePtr().Ptr(),
					t.SessionId.ValuePtr(),
					t.UserId.ValuePtr().Ptr(),
					t.Price.ValuePtr(),
					t.Tickets.ValuePtr(),
					t.Status.ValuePtr(),
					t.Time.ValuePtr().Ptr(),
				); e == nil {
					entries = append(entries, &t)
				} else {
					r.log.Err(e).Msgf("%s:", fn)
				}
			}
			r.log.Debug().Str("SQL", prompt).Msgf("%s: Success", fn)
		} else {
			r.log.Err(e).Str("SQL", prompt).Msgf("%s: Fail to Read from '%s' table", fn, TableCheckout)
		}
		return e
	}); e != nil {
		err = fault.New(rec.ERepoQueryData).Msg(e.Error())
		r.log.Err(e).Str("SQL", prompt).Msgf("%s: Fail to Read from '%s' table", fn, TableCheckout)
	}
	return entries, err
}

func (r *RepositoryUser) readBalanceEntriesBy(where func() string) ([]*TicketBalanceColumn, fault.IError) {
	const fn = "Repository.User.readBalanceEntryBy"
	var err fault.IError
	var e error
	var entries []*TicketBalanceColumn
	prompt := fmt.Sprintf("select * from '%s' where %s", TableBalance, where())
	if e = r.drv.Session(func(db *sql.DB) error {
		var rows *sql.Rows
		if rows, e = db.Query(prompt); e == nil {
			defer rows.Close()
			for rows.Next() {
				t := TicketBalanceColumnDefault()
				if e := rows.Scan(
					t.EventId.ValuePtr().Ptr(),
					t.UserId.ValuePtr().Ptr(),
					t.Change.ValuePtr()); e == nil {
					entries = append(entries, &t)
				} else {
					r.log.Err(e).Msgf("%s:", fn)
				}
			}
			r.log.Debug().Str("SQL", prompt).Msgf("%s: Success", fn)
		} else {
			r.log.Err(e).Str("SQL", prompt).Msgf("%s: Fail to Read from '%s' table", fn, TableBalance)
		}
		return e
	}); e != nil {
		err = fault.New(rec.ERepoQueryData).Msg(e.Error())
		r.log.Err(e).Str("SQL", prompt).Msgf("%s: Fail to Read from '%s' table", fn, TableBalance)
	}
	return entries, err
}
