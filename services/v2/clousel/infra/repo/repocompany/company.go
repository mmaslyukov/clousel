package repocompany

import (
	"clousel/core/business"
	"clousel/infra/repo/driver"
	rec "clousel/infra/repo/errors"
	"clousel/infra/repo/types"
	"clousel/lib/fault"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	TableCompany = "company"
)

type CompanyColumn struct {
	CompanyId   types.Named[types.UUIDString]
	CompanyName types.Named[string]
	Email       types.Named[string]
	Password    types.Named[string]
	ProductId   types.NamedOpt[string]
	SecretKey   types.NamedOpt[string]
	WebhookId   types.NamedOpt[string]
	WebhookKey  types.NamedOpt[string]
	Enabled     types.Named[int]
}

func CompanyColumnDefault() CompanyColumn {
	prof := CompanyColumn{
		CompanyId:   types.NamedCreateDefault[types.UUIDString]("Id"),
		CompanyName: types.NamedCreateDefault[string]("Name"),
		Email:       types.NamedCreateDefault[string]("Email"),
		Password:    types.NamedCreateDefault[string]("Password"),
		ProductId:   types.NamedOptCreateDefault[string]("ProductId"),
		SecretKey:   types.NamedOptCreateDefault[string]("SecKey"),
		WebhookId:   types.NamedOptCreateDefault[string]("WhId"),
		WebhookKey:  types.NamedOptCreateDefault[string]("WhKey"),
		Enabled:     types.NamedCreateDefault[int]("Enabled"),
	}
	return prof
}

func (c *CompanyColumn) toEntry() *business.BusinessEntry {
	entry := &business.BusinessEntry{
		CompanyId:   c.CompanyId.ValuePtr().Uuid(),
		CompanyName: c.CompanyName.Value(),
		Email:       c.Email.Value(),
		Password:    c.Password.Value(),
		ProdId:      c.ProductId.Value(),
		Skey:        c.SecretKey.Value(),
		Whkey:       c.WebhookKey.Value(),
		Whid:        c.WebhookId.Value(),
		Enabled:     c.Enabled.Value() != 0,
	}
	return entry
}

type RepositoryCompany struct {
	drv driver.IDBDriver
	log *zerolog.Logger
}

func RepositoryCompanyCreate(drv driver.IDBDriver, log *zerolog.Logger) *RepositoryCompany {
	return &RepositoryCompany{drv: drv, log: log}
}

func (r *RepositoryCompany) SaveNewBusinessEntry(companyId uuid.UUID, companyName string, email string, password string) fault.IError {
	const fn = "Repository.Company.SaveNewBusinessEntry"
	var err fault.IError
	t := CompanyColumnDefault()
	t.CompanyId.ValuePtr().SetUuid(companyId)
	t.CompanyName.SetValue(companyName)
	t.Email.SetValue(email)
	t.Password.SetValue(password)
	t.Enabled.SetValue(1)
	prompt := fmt.Sprintf("insert into '%s' (%s,%s,%s,%s,%s) values ('%s','%s','%s','%s', %d)", TableCompany,
		t.CompanyId.Name(), t.CompanyName.Name(), t.Email.Name(), t.Password.Name(), t.Enabled.Name(),
		t.CompanyId.ValuePtr().Str(), t.CompanyName.Value(), t.Email.Value(), t.Password.Value(), t.Enabled.Value())
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

func (r *RepositoryCompany) ReadBusinessEntryByName(companyName string) (*business.BusinessEntry, fault.IError) {
	t := CompanyColumnDefault()
	t.CompanyName.SetValue(companyName)
	filter := fmt.Sprintf("%s='%s'", t.CompanyName.Name(), t.CompanyName.Value())
	entries, err := r.readBusinessEntryBy(func() string {
		return filter
	})
	return takeFirst(entries, err)

}

func (r *RepositoryCompany) ReadBusinessEntryById(companyId uuid.UUID) (*business.BusinessEntry, fault.IError) {
	t := CompanyColumnDefault()
	t.CompanyId.ValuePtr().SetUuid(companyId)
	filter := fmt.Sprintf("%s='%s'", t.CompanyId.Name(), t.CompanyId.ValuePtr().Str())
	entries, err := r.readBusinessEntryBy(func() string {
		return filter
	})
	return takeFirst(entries, err)

}

func (r *RepositoryCompany) AssignKeys(companyId uuid.UUID, skey string, prodId string, whid string, whkey string) fault.IError {
	const fn = "Repository.Company.AssignKeys"
	var err fault.IError
	t := CompanyColumnDefault()
	t.CompanyId.ValuePtr().SetUuid(companyId)
	t.SecretKey.SetValue(&skey)
	t.ProductId.SetValue(&prodId)
	t.WebhookId.SetValue(&whid)
	t.WebhookKey.SetValue(&whkey)
	prompt := fmt.Sprintf("update '%s' set %s='%s', %s='%s', %s='%s', %s='%s' where %s='%s'", TableCompany,
		t.ProductId.Name(), *t.ProductId.Value(),
		t.SecretKey.Name(), *t.SecretKey.Value(),
		t.WebhookId.Name(), *t.WebhookId.Value(),
		t.WebhookKey.Name(), *t.WebhookKey.Value(),
		t.CompanyId.Name(), t.CompanyId.ValuePtr().Str())

	if e := r.drv.Session(func(db *sql.DB) error {
		var e error
		if _, e = db.Exec(prompt); e == nil {
			r.log.Debug().Str("SQL", prompt).Msgf("%s: Success", fn)
		} else {
			r.log.Err(e).Str("SQL", prompt).Msgf("%s: Failure", fn)
		}
		return e
	}); e != nil {
		err = fault.New(rec.ERepoExecPrompt).Msgf("Fail to assign company keys %s", e.Error())
	}
	return err
}

/*
Private functions
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
func (r *RepositoryCompany) readBusinessEntryBy(where func() string) ([]*CompanyColumn, fault.IError) {
	const fn = "Repository.Company.readBusinessEntryBy"
	var err fault.IError
	var e error
	// var entries []*business.BusinessEntry
	var entries []*CompanyColumn
	prompt := fmt.Sprintf("select * from '%s' where %s", TableCompany, where())
	if e = r.drv.Session(func(db *sql.DB) error {
		var rows *sql.Rows
		if rows, e = db.Query(prompt); e == nil {
			defer rows.Close()
			for rows.Next() {
				t := CompanyColumnDefault()
				if e := rows.Scan(
					t.CompanyId.ValuePtr().Ptr(),
					t.CompanyName.ValuePtr(),
					t.Email.ValuePtr(),
					t.Password.ValuePtr(),
					t.ProductId.ValuePtr(),
					t.SecretKey.ValuePtr(),
					t.WebhookId.ValuePtr(),
					t.WebhookKey.ValuePtr(),
					t.Enabled.ValuePtr(),
				); e == nil {
					entries = append(entries, &t)
				} else {
					r.log.Err(e).Msgf("%s:", fn)
				}
			}
			r.log.Debug().Str("SQL", prompt).Msgf("%s: Success", fn)
		} else {
			r.log.Err(e).Str("SQL", prompt).Msgf("%s: Fail to Read from '%s' table", fn, TableCompany)
		}
		return e
	}); e != nil {
		err = fault.New(rec.ERepoQueryData).Msg(e.Error())
		r.log.Err(e).Str("SQL", prompt).Msgf("%s: Fail to Read from '%s' table", fn, TableCompany)
	}
	return entries, err
}
