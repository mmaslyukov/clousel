package book_test

import (
	"accountant/core/store"
	"accountant/infra/logger"
	"accountant/infra/repo"
	"accountant/infra/repo/book"
	"accountant/infra/repo/driver"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
)

// Time:      types.NamedCreateDefault[string]("Time"),
// SessionId: types.NamedCreateDefault[string]("SessionId"),
// CarId:     types.NamedCreateDefault[string]("CarouselId"),
// Amount:    types.NamedCreateDefault[int]("Amount"),
// Tickets:    types.NamedCreateDefault[int]("Tickets"),
// Status:    types.NamedCreateDefault[string]("Status"),
// Error:     types.NamedOptCreateDefault[string]("Error"),

func CreateTableProduct(drv driver.IDBDriver) error {
	prompt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS '%s' ('Time' datetime DEFAULT CURRENT_TIMESTAMP, 'SessionId' string UNIQUE NOT NULL,'CarouselId' string NOT NULL,'Amount' int, 'Tickets' int, 'Status' string NOT NULL, 'Error' string)", book.TableBook)
	return drv.Session(func(db *sql.DB) error {
		var err error
		_, err = db.Exec(prompt)
		return err
	})
}

func TestBookRepo(t *testing.T) {
	var ierr store.IError
	var entry, rentry, sentry store.BookEntry
	var entries []store.BookEntry
	const dbPath = "test.db"
	log := logger.New()
	var carId = uuid.New()
	sessionId := "s_hfusd7123s"
	ssessionId := "s_2hfusd7123s"
	drv := repo.DriverSQLite.New(dbPath)
	repoBook := repo.Book.New(drv, &log)
	for ok := true; ok; ok = false {
		if err := CreateTableProduct(drv); err != nil {
			t.Errorf("Fail to create table '%s', err:%s", book.TableBook, err)
			break
		}
		entry = store.BookEntryCreate(sessionId, carId, 1, 1)
		if ierr = repoBook.StoreAddBookEntry(&entry); ierr != nil {
			t.Errorf("Fail to create Book Entry: %s", ierr.Error())
			break
		}
		if rentry, ierr = repoBook.StoreReadBookEntryBySessionId(sessionId); ierr != nil {
			t.Errorf("Fail to read Book Entry: %s", ierr.Error())
			break
		}
		if rentry.Status != store.BookOrderStatusNew {
			t.Errorf("Entry status has unexpected value: %s, expected: %s", rentry.Status, store.BookOrderStatusNew)
			break
		}
		if ierr = repoBook.StoreMarkBookEntryBySessionIdWithData(sessionId, store.BookOrderStatusPaid, nil); ierr != nil {
			t.Errorf("Fail to update enttry status to %s", rentry.Status)
			break
		}
		if rentry, ierr = repoBook.StoreReadBookEntryBySessionId(sessionId); ierr != nil {
			t.Errorf("Fail to read Book Entry: %s", ierr.Error())
			break
		}
		if rentry.Status != store.BookOrderStatusPaid {
			t.Errorf("Entry status has unexpected value: %s, expected: %s", rentry.Status, store.BookOrderStatusNew)
			break
		}
		sentry = store.BookEntryCreate(ssessionId, carId, 2, 2)
		if ierr = repoBook.StoreAddBookEntry(&sentry); ierr != nil {
			t.Errorf("Fail to create Book Entry: %s", ierr.Error())
			break
		}
		if entries, ierr = repoBook.StoreReadBookEntriesByCarosuelId(carId); ierr != nil {
			t.Errorf("Fail to read Book Entries: %s", ierr.Error())
			break
		}
		if len(entries) != 2 {
			t.Errorf("Read entries size is unexpected %d", len(entries))
			break
		}
	}
	os.Remove(dbPath)
	// StoreReadBookEntriesByCarosuelId(carId Carousel) ([]BookEntry, IError)
}
