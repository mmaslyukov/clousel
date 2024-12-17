package profile_test

import (
	"accountant/core/owner"
	"accountant/infra/logger"
	"accountant/infra/repo"
	"accountant/infra/repo/driver"
	"accountant/infra/repo/profile"
	"database/sql"
	"fmt"
	"os"
	"testing"
)

func CreateTable(drv driver.IDBDriver) error {
	_ = profile.ProfileColumDefault()
	prompt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS '%s' ('OwnerId' string NOT NULL,'Email' string UNIQUE NOT NULL,'Password' string NOT NULL,'SecretKey' string UNIQUE,'PublishKey' string UNIQUE,'WebhookId' string UNIQUE,'WebhookKey' string UNIQUE,'Role' integer,'Time' datetime DEFAULT CURRENT_TIMESTAMP)", profile.TableProfile)
	return drv.Session(func(db *sql.DB) error {
		var err error
		_, err = db.Exec(prompt)
		return err
	})
}

func TestOwnerRegister(t *testing.T) {
	const dbPath = "test.db"
	const email = "test@mail.com"
	const password = "123qweasd"
	log := logger.New()
	drv := repo.DriverSQLite.New(dbPath)
	repo := repo.Profile.New(drv, &log)

	for ok := true; ok; ok = false {
		if err := CreateTable(drv); err != nil {
			t.Errorf("Fail to create table '%s'", profile.TableProfile)
			break
		}
		err := repo.OwnerRegister(email, password, owner.UserRoleRegular)
		if err != nil {
			t.Errorf("Fail to register Owner: %s", err.Full())
			break
		}
		entry, err := repo.OwnerReadEntry(email)
		if err != nil {
			t.Errorf("Fail to read Owner entry: %s", err.Full())
			break
		}
		if entry.Email != email {
			t.Errorf("Emails are mismatch: %s != %s", entry.Email, email)
		}
		if entry.Password != password {
			t.Errorf("Passwords are mismatch: %s != %s", entry.Password, password)
		}
	}
	os.Remove(dbPath)
}

func TestAssignStripeKeys(t *testing.T) {
	const dbPath = "test.db"
	const email = "test@mail.com"
	const password = "123qweasd"
	log := logger.New()
	drv := repo.DriverSQLite.New(dbPath)
	repo := repo.Profile.New(drv, &log)

	for ok := true; ok; ok = false {
		if err := CreateTable(drv); err != nil {
			t.Errorf("Fail to create table '%s'", profile.TableProfile)
			break
		}
		err := repo.OwnerRegister(email, password, owner.UserRoleRegular)
		if err != nil {
			t.Errorf("Fail to register Owner: %s", err.Full())
			break
		}
		entry, err := repo.OwnerReadEntry(email)
		if err != nil {
			t.Errorf("Fail to read Owner entry: %s", err.Full())
			break
		}
		skey := "skey"
		pkey := "pkey"
		if err = repo.OwnerAssignStripeKeys(entry.OwnerId, &pkey, &skey); err != nil {
			t.Errorf("Fail to read Assign Keys: %s", err.Full())
			break
		}
		const whkey = "whkey"
		const whid = "whid"
		if err = repo.OwnerAssignWebhook(entry.OwnerId, whid, whkey); err != nil {
			t.Errorf("Fail to read Assign Wh Key: %s", err.Full())
			break
		}
		entry, err = repo.OwnerReadEntry(email)
		if err != nil {
			t.Errorf("Fail to read Owner entry: %s", err.Full())
			break
		}
		if entry.PublishKey == nil || entry.SecretKey == nil || entry.WebhookId == nil || entry.WebhookKey == nil {
			t.Errorf("Some keys are nil")
		}

		if *entry.PublishKey != pkey {
			t.Errorf("Publish key is mismatch: %s != %s", *entry.PublishKey, pkey)
		}
		if *entry.SecretKey != skey {
			t.Errorf("Secret key is mismatch: %s != %s", *entry.SecretKey, skey)
		}
		if *entry.WebhookKey != whkey {
			t.Errorf("Webhook key is mismatch: %s != %s", *entry.SecretKey, skey)
		}
	}
	os.Remove(dbPath)
}
