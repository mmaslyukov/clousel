package repo

import (
	"clousel/infra/repo/driver"
	"clousel/infra/repo/repocompany"
	"clousel/infra/repo/repomachine"
	"clousel/infra/repo/repouser"

	"github.com/rs/zerolog"
	_ "modernc.org/sqlite"
)

type exportCompany struct{}

func (e exportCompany) New(drv driver.IDBDriver, log *zerolog.Logger) *repocompany.RepositoryCompany {
	return repocompany.RepositoryCompanyCreate(drv, log)
}

type exportUser struct{}

func (e exportUser) New(drv driver.IDBDriver, log *zerolog.Logger) *repouser.RepositoryUser {
	return repouser.RepositoryUserCreate(drv, log)
}

type exportMachine struct{}

func (e exportMachine) New(drv driver.IDBDriver, log *zerolog.Logger) *repomachine.RepositoryMachine {
	return repomachine.RepositoryMachineCreate(drv, log)
}

type exportDriver struct{}

func (e exportDriver) New(path string) driver.IDBDriver {
	return driver.DriverSqliteCreate(path)
}

var DriverSQLite = exportDriver{}

var Company = exportCompany{}
var User = exportUser{}
var Machine = exportMachine{}
