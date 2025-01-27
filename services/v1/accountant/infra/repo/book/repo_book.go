package book

import (
	"accountant/core/store"
	errs "accountant/core/store/error"
	"accountant/infra/repo/driver"
	"accountant/infra/repo/types"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	TableBook = "book"
)

type BookColumn struct {
	Time      types.Named[string]
	SessionId types.Named[string]
	CarId     types.Named[string]
	Amount    types.Named[int]
	Tickets   types.Named[int]
	Status    types.Named[string]
	Error     types.NamedOpt[string]
}

func BookColumnDefault() BookColumn {
	return BookColumn{
		Time:      types.NamedCreateDefault[string]("Time"),
		SessionId: types.NamedCreateDefault[string]("SessionId"),
		CarId:     types.NamedCreateDefault[string]("CarouselId"),
		Amount:    types.NamedCreateDefault[int]("Amount"),
		Tickets:   types.NamedCreateDefault[int]("Tickets"),
		Status:    types.NamedCreateDefault[string]("Status"),
		Error:     types.NamedOptCreateDefault[string]("Error"),
	}
}

func (c *BookColumn) toEntry() store.BookEntry {
	cid, _ := uuid.Parse(c.CarId.Value)
	return store.BookEntry{
		SessionId: c.SessionId.Value,
		Time:      c.Time.Value,
		CarId:     cid,
		Amount:    c.Amount.Value,
		Tickets:   c.Tickets.Value,
		Status:    c.Status.Value,
		Error:     c.Error.Value,
	}
}
func (c *BookColumn) fromEntry(entry *store.BookEntry) {
	c.SessionId.Value = entry.SessionId
	c.CarId.Value = entry.CarId.String()
	c.Amount.Value = entry.Amount
	c.Tickets.Value = entry.Tickets
	// c.Status.Value = store.BookOrderStatusNew
}

type RepositoryBook struct {
	drv driver.IDBDriver
	log *zerolog.Logger
}

func RepositoryBookCreate(drv driver.IDBDriver, log *zerolog.Logger) *RepositoryBook {
	return &RepositoryBook{drv: drv, log: log}
}

func (r *RepositoryBook) StoreAddBookEntry(entry *store.BookEntry) errs.IError {
	var prompt string
	var ierr errs.IError
	c := BookColumnDefault()
	c.fromEntry(entry)
	c.Status.Value = store.BookOrderStatusNew

	prompt = fmt.Sprintf("insert into '%s' (%s,%s,%s,%s,%s) values ('%s','%s',%d,%d,'%s')", TableBook,
		c.SessionId.Name(), c.CarId.Name(), c.Amount.Name(), c.Tickets.Name(), c.Status.Name(),
		c.SessionId.Value, c.CarId.Value, c.Amount.Value, c.Tickets.Value, c.Status.Value)
	if e := r.drv.Session(func(db *sql.DB) error {
		var err error
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Send()
		} else {
			r.log.Err(err).Str("SQL", prompt).Msg("Repository.Book.StoreAddBookEntry: Fail to Add")
		}
		return err
	}); e != nil {
		ierr = errs.New(errs.ECBookRepoInsert).Msg(e.Error())
	}
	return ierr
}

func (r *RepositoryBook) StoreMarkBookEntryBySessionIdWithData(sessionId store.Session, status string, err *string) errs.IError {
	var ierr errs.IError
	var e error
	c := BookColumnDefault()
	c.SessionId.Value = sessionId
	c.Status.Value = status
	c.Error.Value = err
	var prompt string
	if c.Error.Value != nil {
		prompt = fmt.Sprintf("update '%s' set %s='%s', %s='%s' where %s='%s'", TableBook,
			c.Status.Name(), c.Status.Value,
			c.Error.Name(), *c.Error.Value,
			c.SessionId.Name(), c.SessionId.Value)
	} else {
		prompt = fmt.Sprintf("update '%s' set %s='%s' where %s='%s'", TableBook,
			c.Status.Name(), c.Status.Value,
			c.SessionId.Name(), c.SessionId.Value)
	}
	if e = r.drv.Session(func(db *sql.DB) error {
		if _, e = db.Exec(prompt); e == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Book.StoreMarkBookEntryBySessionIdWithData: Success")
		} else {
			r.log.Err(e).Str("SQL", prompt).
				Str(c.SessionId.Name(), c.SessionId.Value).
				Str(c.Status.Name(), c.Status.Value).
				Msg("Repository.Book.StoreMarkBookEntryBySessionIdWithData: Failure")
		}
		return e
	}); e != nil {
		ierr = errs.New(errs.ECBookRepoMark).Msgf("Fail to assign ProductId %s", e.Error())
	}
	return ierr
}

func (r *RepositoryBook) StoreReadBookEntryBySessionId(sessionId store.Session) (store.BookEntry, errs.IError) {
	var ierr errs.IError
	var e error
	c := BookColumnDefault()
	c.SessionId.Value = sessionId
	prompt := fmt.Sprintf("select * from '%s' where %s='%s'", TableBook, c.SessionId.Name(), c.SessionId.Value)
	if e = r.drv.Session(func(db *sql.DB) error {
		if e = db.QueryRow(prompt).Scan(
			&c.Time.Value,
			&c.SessionId.Value,
			&c.CarId.Value,
			&c.Amount.Value,
			&c.Tickets.Value,
			&c.Status.Value,
			&c.Error.Value); e == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Book.StoreReadBookEntryBySessionId: Success")
		}
		return e
	}); e != nil {
		ierr = errs.New(errs.ECBookRepoRead).Msg(e.Error())
		r.log.Err(e).Str("SQL", prompt).Msgf("Repository.Book.StoreReadBookEntryBySessionId: Fail to Read from '%s' table", TableBook)
	}
	return c.toEntry(), ierr
}

func (r *RepositoryBook) StoreReadBookEntriesByCarosuelId(carId store.Carousel) ([]store.BookEntry, errs.IError) {
	var ierr errs.IError
	var e error
	var entries []store.BookEntry
	c := BookColumnDefault()
	c.CarId.Value = carId.String()
	prompt := fmt.Sprintf("select * from '%s' where %s='%s'", TableBook, c.CarId.Name(), c.CarId.Value)
	if e = r.drv.Session(func(db *sql.DB) error {
		var rows *sql.Rows
		if rows, e = db.Query(prompt); e == nil {
			defer rows.Close()
			for rows.Next() {
				if e := rows.Scan(
					&c.Time.Value,
					&c.SessionId.Value,
					&c.CarId.Value,
					&c.Amount.Value,
					&c.Tickets.Value,
					&c.Status.Value,
					&c.Error.Value); e == nil {
					entries = append(entries, c.toEntry())
				} else {
					ierr = errs.New(errs.ECBookRepoRead).Msgf("Fail to Scan entry, error:%s", e.Error())
					r.log.Err(e).Msgf("Repository.Book.StoreReadBookEntriesByCarosuelId: %s", ierr.Error())
				}
			}
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Book.StoreReadBookEntriesByCarosuelId: Success")
		} else {
			r.log.Err(e).Str("SQL", prompt).Msgf("Repository.Book.StoreReadBookEntriesByCarosuelId: Fail to Read from '%s' table", TableBook)
		}
		return e
	}); e != nil {
		ierr = errs.New(errs.ECBookRepoRead).Msg(e.Error())
		r.log.Err(e).Str("SQL", prompt).Msgf("Repository.Book.StoreReadBookEntriesByCarosuelId: Fail to Read from '%s' table", TableBook)
	}
	return entries, ierr
}
