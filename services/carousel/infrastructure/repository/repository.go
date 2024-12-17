package repository

import (
	"carousel/infrastructure/repository/carousel"
	"carousel/infrastructure/repository/driver"
	"carousel/infrastructure/repository/event"
	"carousel/infrastructure/repository/snapshot"

	"github.com/rs/zerolog"
	_ "modernc.org/sqlite"
)

type exportCarousel struct{}

func (e exportCarousel) New(drv driver.IDBDriver, log *zerolog.Logger) *carousel.RepositoryCarousel {
	return carousel.New(drv, log)
}

type exportEvet struct{}

func (e exportEvet) New(drv driver.IDBDriver, crRepo event.IRepositoryCarousel, log *zerolog.Logger) *event.RepositoryEvent {
	return event.New(drv, crRepo, log)
}

type exportSnapshot struct{}

func (e exportSnapshot) New(drv driver.IDBDriver, log *zerolog.Logger) *snapshot.RepositorySnapshot {
	return snapshot.New(drv, log)
}

type exportDriver struct{}

func (e exportDriver) New(path string) driver.IDBDriver {
	return driver.New(path)
}

var DriverSQLite = exportDriver{}
var Carousel = exportCarousel{}
var Event = exportEvet{}
var Snaphsot = exportSnapshot{}
