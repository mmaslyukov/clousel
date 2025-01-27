package repo

import (
	"accountant/infra/repo/book"
	"accountant/infra/repo/driver"
	"accountant/infra/repo/prod"
	"accountant/infra/repo/profile"

	"github.com/rs/zerolog"
	_ "modernc.org/sqlite"
)

type exportProfile struct{}

func (e exportProfile) New(drv driver.IDBDriver, log *zerolog.Logger) *profile.RepositoryProfile {
	return profile.RepositoryProfileCreate(drv, log)
}

type exportProduct struct{}

func (e exportProduct) New(drv driver.IDBDriver, log *zerolog.Logger) *prod.RepositoryProduct {
	return prod.RepositoryProductCreate(drv, log)
}

type exportBook struct{}

func (e exportBook) New(drv driver.IDBDriver, log *zerolog.Logger) *book.RepositoryBook {
	return book.RepositoryBookCreate(drv, log)
}

type exportDriver struct{}

func (e exportDriver) New(path string) driver.IDBDriver {
	return driver.DriverSqliteCreate(path)
}

var DriverSQLite = exportDriver{}

var Profile = exportProfile{}
var Product = exportProduct{}
var Book = exportBook{}
