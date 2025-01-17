package repocompany_test

import (
	"clousel/core/business"
	"clousel/infra/log"
	"clousel/infra/repo"
	"clousel/infra/repo/driver"
	"clousel/infra/repo/repocompany"
	"clousel/lib/fault"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
)

func CreateCompanyTable(drv driver.IDBDriver) error {
	prompt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS '%s'(
		'Id' string PRIMARY KEY,
		'Name' string UNIQUE NOT NULL,
		'Email' string UNIQUE NOT NULL,
		'Password' string NOT NULL,
		'ProductId' string UNIQUE,
		'SecKey' string UNIQUE,
		'WhId' string,
		'WhKey' string)`, repocompany.TableCompany)
	return drv.Session(func(db *sql.DB) error {
		var err error
		_, err = db.Exec(prompt)
		return err
	})
}

// go test -v -run ^TestIBusinessRepoAdapterInterface$  .\infra\repo\repomachine\
func TestIBusinessRepoAdapterInterface(t *testing.T) {
	const dbPath = "TestIBusinessRepoAdapterInterface.db"
	const email = "any@mail.com"
	const companyName = "Default"
	const password = "123qweasd"
	companyId := uuid.New()
	log := log.New()
	var err fault.IError
	drv := repo.DriverSQLite.New(dbPath)
	repo := repo.Company.New(drv, log)

	for ok := true; ok; ok = false {

		if err := CreateCompanyTable(drv); err != nil {
			t.Errorf("Fail to create table '%s', err:%s", repocompany.TableCompany, err)
			break
		}

		if err = repo.SaveNewBusinessEntry(companyId, companyName, email, password); err != nil {
			t.Errorf("Fail to create new machine entry: %s", err.Error())
			break
		}
		skey := "skey_1234"
		prodid := "prodid_1234"
		whid := "whid_1234"
		whkey := "whkey_1234"
		if err = repo.AssignKeys(companyId, skey, prodid, whid, whkey); err != nil {
			t.Errorf("Fail to create new machine entry: %s", err.Error())
			break
		}
		var entry *business.BusinessEntry
		if entry, err = repo.ReadBusinessEntryById(companyId); err != nil {
			t.Errorf("Fail to read machine entry: %s", err.Error())
			break
		}

		if entry.CompanyId != companyId {
			t.Errorf("CompanyId mismatch")
		}
		if entry.CompanyName != companyName {
			t.Errorf("CompanyName mismatch")
		}
		if entry.Email != email {
			t.Errorf("Email mismatch")
		}
		if entry.ProdId == nil || *entry.ProdId != prodid {
			t.Errorf("Prodid mismatch")
		}
		if entry.Skey == nil || *entry.Skey != skey {
			t.Errorf("Skey mismatch")
		}
		if entry.Whid == nil || *entry.Whid != whid {
			t.Errorf("Whid mismatch")
		}
		if entry.Whkey == nil || *entry.Whkey != whkey {
			t.Errorf("Whkey mismatch")
		}

		// t.Logf("%+v", *ml[1])

	}
	os.Remove(dbPath)
}
