package driver

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type Lambda func(*sql.DB) error
type IDBDriver interface {
	Session(Lambda) error
}

type DriverSqlite struct {
	path string
}

func (d *DriverSqlite) Session(l Lambda) error {
	var err error
	var db *sql.DB
	if db, err = sql.Open("sqlite", d.path); err == nil {
		defer db.Close()
		err = l(db)
	}
	return err
}

func New(path string) IDBDriver {
	return &DriverSqlite{path: path}
}
