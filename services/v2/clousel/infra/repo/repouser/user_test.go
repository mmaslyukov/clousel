package repouser_test

import (
	"clousel/core/client"
	"clousel/infra/log"
	"clousel/infra/repo"
	"clousel/infra/repo/driver"
	"clousel/infra/repo/repouser"
	"clousel/lib/fault"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
)

func CreateUserTable(drv driver.IDBDriver) error {
	// 'Balance' int ,
	prompt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS '%s'(
		'Id' string PRIMARY KEY, 
		'Username' string NOT NULL UNIQUE,
		'Email' string UNIQUE NOT NULL,
		'Password' string NOT NULL,
		'Time' datetime)`, repouser.TableUser)
	return drv.Session(func(db *sql.DB) error {
		var err error
		_, err = db.Exec(prompt)
		return err
	})
}

func CreateCheckoutTable(drv driver.IDBDriver) error {
	prompt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS '%s'(
		'EventId' string PRIMARY KEY, 
		'SessionId' string NOT NULL UNIQUE,
		'UserId' string NOT NULL UNIQUE,
		'Price' int NOT NULL,
		'Tickets' int NOT NULL,
		'Status' string NOT NULL,
		'Time' datetime)`, repouser.TableCheckout)
	return drv.Session(func(db *sql.DB) error {
		var err error
		_, err = db.Exec(prompt)
		return err
	})
}

func CreateBalanceTable(drv driver.IDBDriver) error {
	prompt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS '%s'(
		'EventId' string PRIMARY KEY, 
		'UserId' string NOT NULL,
		'Change' int NOT NULL)`, repouser.TableBalance)
	return drv.Session(func(db *sql.DB) error {
		var err error
		_, err = db.Exec(prompt)
		return err
	})
}

// go test -v -run ^TestIClientRepoGeneralAdapterInterface$  .\infra\repo\user\
func TestIClientRepoGeneralAdapterInterface(t *testing.T) {
	const dbPath = "TestIClientRepoGeneralAdapterInterface.db"
	const email = "test@mail.com"
	const username = "JohnDoe"
	const password = "123qweasd"
	const companyName = "default"
	userId := uuid.New()
	log := log.New()
	var err fault.IError
	drv := repo.DriverSQLite.New(dbPath)
	repo := repo.User.New(drv, log)

	for ok := true; ok; ok = false {

		if err := CreateUserTable(drv); err != nil {
			t.Errorf("Fail to create table '%s', err:%s", repouser.TableUser, err)
			break
		}
		if err = repo.SaveNewClientEntry(userId, username, email, password, companyName); err != nil {
			t.Errorf("Fail to create new client entry: %s", err.Error())
			break
		}
		var entry1 *client.ClientEntry
		if entry1, err = repo.ReadClientEntryByName(username); err != nil {
			t.Errorf("Fail to read client entry: %s", err.Error())
			break
		}
		// const newBalance = 11
		// if err = repo.UpdateBalance(userId, newBalance); err != nil {
		// 	t.Errorf("Fail to update user balance: %s", err.Error())
		// 	break
		// }
		if entry1.Email != email ||
			entry1.Password != password ||
			entry1.Username != username {
			t.Errorf("Client data is mismatch")
		}
		var entry2 *client.ClientEntry
		if entry2, err = repo.ReadClientEntryById(entry1.UserId); err != nil {
			t.Errorf("Fail to read client entry: %s", err.Error())
			break
		}
		// log.Printf("%+v\n", entry)
		// log.Printf("%+v\n", entry2)
		if entry1.UserId != entry2.UserId {
			t.Errorf("Entries aren't the same")
			break
		}
		// if entry1.Balance == entry2.Balance || entry2.Balance != newBalance {
		// 	t.Errorf("Unexpected balance value %d vs %d", entry2.Balance, newBalance)
		// 	break

		// }

	}
	os.Remove(dbPath)
}

// go test -v -run ^TestIClientRepoCheckoutSessionAdapterInterface$  .\infra\repo\user\
func TestIClientRepoCheckoutSessionAdapterInterface(t *testing.T) {
	const dbPath = "TestIClientRepoCheckoutSessionAdapterInterface.db"
	const email = "test@mail.com"
	const username = "JohnDoe"
	const password = "123qweasd"
	const companyName = "default"
	userId := uuid.New()
	log := log.New()
	var err fault.IError
	drv := repo.DriverSQLite.New(dbPath)
	repo := repo.User.New(drv, log)

	for ok := true; ok; ok = false {

		if err := CreateUserTable(drv); err != nil {
			t.Errorf("Fail to create table '%s', err:%s", repouser.TableUser, err)
			break
		}
		if err := CreateCheckoutTable(drv); err != nil {
			t.Errorf("Fail to create table '%s', err:%s", repouser.TableCheckout, err)
			break
		}
		if err = repo.SaveNewClientEntry(userId, username, email, password, companyName); err != nil {
			t.Errorf("Fail to create new client entry: %s", err.Error())
			break
		}
		var ce *client.ClientEntry
		if ce, err = repo.ReadClientEntryByName(username); err != nil {
			t.Errorf("Fail to read client entry: %s", err.Error())
			break
		}
		csevt := uuid.New()
		sid := "test_session_id"
		if err = repo.SaveNewCheckoutEntry(csevt, ce.UserId, sid, 3, 6); err != nil {
			t.Errorf("Fail to save new checkout session: %s", err.Error())
			break
		}
		var cses []*client.CheckoutEntry
		if cses, err = repo.ReadCheckoutEntriesByUserId(ce.UserId, nil, nil); err != nil {
			t.Errorf("Fail to save new checkout session: %s", err.Error())
			break
		}
		if len(cses) != 1 {
			t.Errorf("Unexpected number of read records have got %d, but expect 1", len(cses))
			break
		}
		if err = repo.UpdateCheckoutStatus(cses[0].SessionId, client.PaymentStatusPaid); err != nil {
			t.Errorf("Fail to update payment status")
			break
		}
		if cses, err = repo.ReadCheckoutEntriesByUserId(ce.UserId, nil, nil); err != nil {
			t.Errorf("Fail to save new checkout session: %s", err.Error())
			break
		}
		if cses[0].Status != client.PaymentStatusPaid {
			t.Errorf("Unexpected payment status have got '%s', but expect '%s'", cses[0].Status, client.PaymentStatusPaid)
			break
		}
		fmt.Printf("%+v\n", *cses[0])

	}
	os.Remove(dbPath)
}

// go test -v -run ^TestIClientRepoBalanceChangeAdapterInterface$  .\infra\repo\user\
func TestIClientRepoBalanceChangeAdapterInterface(t *testing.T) {
	const dbPath = "TestIClientRepoBalanceChangeAdapterInterface.db"
	const email = "test@mail.com"
	const username = "JohnDoe"
	const password = "123qweasd"
	const companyName = "default"
	userId := uuid.New()
	log := log.New()
	var err fault.IError
	drv := repo.DriverSQLite.New(dbPath)
	repo := repo.User.New(drv, log)

	for ok := true; ok; ok = false {

		if err := CreateUserTable(drv); err != nil {
			t.Errorf("Fail to create table '%s', err:%s", repouser.TableUser, err)
			break
		}
		if err := CreateBalanceTable(drv); err != nil {
			t.Errorf("Fail to create table '%s', err:%s", repouser.TableCheckout, err)
			break
		}
		if err = repo.SaveNewClientEntry(userId, username, email, password, companyName); err != nil {
			t.Errorf("Fail to create new client entry: %s", err.Error())
			break
		}
		var ce *client.ClientEntry
		if ce, err = repo.ReadClientEntryByName(username); err != nil {
			t.Errorf("Fail to read client entry: %s", err.Error())
			break
		}
		if err = repo.SaveNewBalanceChangeEntry(uuid.New(), ce.UserId, 3); err != nil {
			t.Errorf("Fail to save balance change entry: %s", err.Error())
			break
		}
		if err = repo.SaveNewBalanceChangeEntry(uuid.New(), ce.UserId, 2); err != nil {
			t.Errorf("Fail to save balance change entry: %s", err.Error())
			break
		}
		if err = repo.SaveNewBalanceChangeEntry(uuid.New(), ce.UserId, -1); err != nil {
			t.Errorf("Fail to save balance change entry: %s", err.Error())
			break
		}
		var tbe []*client.TicketsBalanceEntry
		if tbe, err = repo.ReadBalanceEntriesByUserId(ce.UserId); err != nil {
			t.Errorf("Fail to read balancee change entry: %s", err.Error())
			break
		}
		if len(tbe) != 3 {
			t.Errorf("Expect to have 3 rows, but have got: %d", len(tbe))
			break
		}
		// var b int
		// if b, err = repo.CalculateBalanceTotal(ce.UserId); err != nil {
		// 	t.Errorf("Fail to calculate balance total %s", err.Error())
		// 	break
		// }
		// expect := 4
		// if b != expect {
		// 	t.Errorf("Expect to have balance as %d tickent, but have got %d tickets", expect, b)
		// 	break
		// }

		// if err = repo.ClearBalanceByUserId(ce.UserId); err != nil {
		// 	t.Errorf("Fail to clear balance history %s", err.Error())
		// 	break
		// }

		if tbe, err = repo.ReadBalanceEntriesByUserId(ce.UserId); err != nil {
			t.Errorf("Fail to read balancee change entry: %s", err.Error())
			break

		} else if len(tbe) != 0 {
			t.Errorf("Fail to read 0 entries but got %d", len(tbe))
			break
		}

	}
	os.Remove(dbPath)
}
