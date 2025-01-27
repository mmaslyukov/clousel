package profile

import (
	"accountant/core/owner"
	erro "accountant/core/owner/error"
	"accountant/infra/repo/driver"
	"accountant/infra/repo/types"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	TableProfile = "profile"
)

type ProfileColumn struct {
	OwnerId    types.Named[string]
	Email      types.Named[string]
	Password   types.Named[string]
	SecretKey  types.NamedOpt[string]
	PublishKey types.NamedOpt[string]
	WebhookId  types.NamedOpt[string]
	WebhookKey types.NamedOpt[string]
	Role       types.Named[owner.UserRole]
	Time       types.Named[string]
}

func ProfileColumDefault() ProfileColumn {
	prof := ProfileColumn{
		OwnerId:    types.NamedCreateDefault[string]("OwnerId"),
		Email:      types.NamedCreateDefault[string]("Email"),
		Password:   types.NamedCreateDefault[string]("Password"),
		SecretKey:  types.NamedOptCreateDefault[string]("SecretKey"),
		PublishKey: types.NamedOptCreateDefault[string]("PublishKey"),
		WebhookId:  types.NamedOptCreateDefault[string]("WebhookId"),
		WebhookKey: types.NamedOptCreateDefault[string]("WebhookKey"),
		Role:       types.NamedCreateDefault[owner.UserRole]("Role"),
		Time:       types.NamedCreateDefault[string]("Time"),
	}
	prof.Role.Value = owner.UserRoleRegular
	return prof
}

type RepositoryProfile struct {
	drv driver.IDBDriver
	log *zerolog.Logger
}

func RepositoryProfileCreate(drv driver.IDBDriver, log *zerolog.Logger) *RepositoryProfile {
	return &RepositoryProfile{drv: drv, log: log}
}

func (r *RepositoryProfile) OwnerRegister(email string, password string, role owner.UserRole) erro.IError {
	var prompt string
	var ierr erro.IError
	c := ProfileColumDefault()
	c.OwnerId.Value = uuid.New().String()
	c.Email.Value = email
	c.Password.Value = password
	c.Role.Value = role
	prompt = fmt.Sprintf("insert into '%s' (%s, %s, %s, %s) values ('%s', '%s', '%s', %d)", TableProfile,
		c.OwnerId.Name(), c.Email.Name(), c.Password.Name(), c.Role.Name(),
		c.OwnerId.Value, c.Email.Value, c.Password.Value, c.Role.Value)
	if e := r.drv.Session(func(db *sql.DB) error {
		var err error
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Send()
		} else {
			r.log.Err(err).Str("SQL", prompt).Msg("Repository.Profile.OwnerRegister: Fail to Register")
		}
		return err
	}); e != nil {
		ierr = erro.New(erro.ECRepoExecPrompt).Msg(e.Error())
	}
	return ierr

}

func (r *RepositoryProfile) readEntryByStr(condition string) (owner.OwnerEntry, erro.IError) {
	var ierr erro.IError
	var e error
	c := ProfileColumDefault()
	prompt := fmt.Sprintf("select * from '%s' where %s", TableProfile, condition)
	if e = r.drv.Session(func(db *sql.DB) error {
		if e = db.QueryRow(prompt).Scan(
			&c.OwnerId.Value,
			&c.Email.Value,
			&c.Password.Value,
			&c.SecretKey.Value,
			&c.PublishKey.Value,
			&c.WebhookId.Value,
			&c.WebhookKey.Value,
			&c.Role.Value,
			&c.Time.Value); e == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Profile.OwnerRead: Success")
		}
		return e
	}); e != nil {
		ierr = erro.New(erro.ECRepoExecPrompt).Msg(e.Error())
		r.log.Err(e).Str("SQL", prompt).Msgf("Repository.Profile.OwnerRead: Fail to Read from '%s' table", TableProfile)
	}
	oid, e := uuid.Parse(c.OwnerId.Value)
	if e != nil {
		ierr = erro.New(erro.ECGeneralFailure).Msg(e.Error())
	}
	return owner.OwnerEntry{
		OwnerId:    oid,
		Email:      c.Email.Value,
		Password:   c.Password.Value,
		SecretKey:  c.SecretKey.Value,
		PublishKey: c.PublishKey.Value,
		WebhookId:  c.WebhookId.Value,
		WebhookKey: c.WebhookKey.Value,
		Role:       c.Role.Value,
	}, ierr
}

func (r *RepositoryProfile) OwnerReadEntry(email string) (owner.OwnerEntry, erro.IError) {
	c := ProfileColumDefault()
	c.Email.Value = email
	return r.readEntryByStr(fmt.Sprintf("%s='%s'", c.Email.Name(), c.Email.Value))
	// var ierr erro.IError
	// var e error
	// c := ProfileColumDefault()
	// c.Email.Value = email
	// prompt := fmt.Sprintf("select * from '%s' where %s='%s'", TableProfile, c.Email.Name(), c.Email.Value)
	// if e = r.drv.Session(func(db *sql.DB) error {
	// 	if e = db.QueryRow(prompt).Scan(
	// 		&c.OwnerId.Value,
	// 		&c.Email.Value,
	// 		&c.Password.Value,
	// 		&c.SecretKey.Value,
	// 		&c.PublishKey.Value,
	// 		&c.Role.Value,
	// 		&c.Time.Value); e == nil {
	// 		r.log.Debug().Str("SQL", prompt).Msg("Repository.Profile.OwnerRead: Success")
	// 	}
	// 	return e
	// }); e != nil {
	// 	ierr = erro.New(erro.ECRepoExecPrompt).Msg(e.Error())
	// 	r.log.Err(e).Str("SQL", prompt).Msgf("Repository.Profile.OwnerRead: Fail to Read from '%s' table", TableProfile)
	// }
	// oid, e := uuid.Parse(c.OwnerId.Value)
	// if e != nil {
	// 	ierr = erro.New(erro.ECGeneralFailure).Msg(e.Error())
	// }
	// return owner.OwnerEntry{
	// 	OwnerId:    oid,
	// 	Email:      c.Email.Value,
	// 	Password:   c.Password.Value,
	// 	SecretKey:  c.SecretKey.Value,
	// 	PublishKey: c.PublishKey.Value,
	// 	Role:       c.Role.Value,
	// }, ierr
}

func (r *RepositoryProfile) OwnerReadEntryByOwner(ownerId string) (owner.OwnerEntry, erro.IError) {
	c := ProfileColumDefault()
	c.OwnerId.Value = ownerId
	return r.readEntryByStr(fmt.Sprintf("%s='%s'", c.OwnerId.Name(), c.OwnerId.Value))
}

func (r *RepositoryProfile) OwnerAssignStripeKeys(ownerId owner.Owner, pk *string, sk *string) erro.IError {
	var ierr erro.IError
	var prompt string
	c := ProfileColumDefault()
	c.OwnerId.Value = ownerId.String()
	c.SecretKey.Value = sk
	c.PublishKey.Value = pk

	if pk != nil && sk == nil {
		prompt = fmt.Sprintf("update '%s' set %s='%s' where %s='%s'", TableProfile,
			c.PublishKey.Name(), *c.PublishKey.Value,
			c.OwnerId.Name(), c.OwnerId.Value)
	} else if sk != nil && pk == nil {
		prompt = fmt.Sprintf("update '%s' set %s='%s' where %s='%s'", TableProfile,
			c.SecretKey.Name(), *c.SecretKey.Value,
			c.OwnerId.Name(), c.OwnerId.Value)
	} else if sk != nil && pk != nil {
		prompt = fmt.Sprintf("update '%s' set %s='%s', %s='%s' where %s='%s'", TableProfile,
			c.SecretKey.Name(), *c.SecretKey.Value, c.PublishKey.Name(), *c.PublishKey.Value,
			c.OwnerId.Name(), c.OwnerId.Value)
	}
	if e := r.drv.Session(func(db *sql.DB) error {
		var e error
		if _, e = db.Exec(prompt); e == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Profile.AssignKeys: Success")
		} else {
			r.log.Err(e).Str("SQL", prompt).Msg("Repository.Profile.AssignKeys: Failure")
		}
		return e
	}); e != nil {
		ierr = erro.New(erro.ECRepoExecPrompt).Msg(e.Error())
	}
	return ierr
}

func (r *RepositoryProfile) OwnerAssignWebhook(ownerId owner.Owner, whid string, whkey string) erro.IError {
	var ierr erro.IError
	var prompt string
	c := ProfileColumDefault()
	c.OwnerId.Value = ownerId.String()
	c.WebhookKey.Value = &whkey
	c.WebhookId.Value = &whid

	prompt = fmt.Sprintf("update '%s' set %s='%s', %s='%s' where %s='%s'", TableProfile,
		c.WebhookKey.Name(), *c.WebhookKey.Value,
		c.WebhookId.Name(), *c.WebhookId.Value,
		c.OwnerId.Name(), c.OwnerId.Value)
	if e := r.drv.Session(func(db *sql.DB) error {
		var e error
		if _, e = db.Exec(prompt); e == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Profile.OwnerAssignWebhook: Success")
		} else {
			r.log.Err(e).Str("SQL", prompt).Msg("Repository.Profile.OwnerAssignWebhook: Failure")
		}
		return e
	}); e != nil {
		ierr = erro.New(erro.ECRepoExecPrompt).Msg(e.Error())
	}
	return ierr
}
