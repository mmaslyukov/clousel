package persistency

import (
	"fmt"

	"carousel_service/internal/config"
	"database/sql"

	pp "carousel_service/internal/ports/port_persistency"
	_ "modernc.org/sqlite"
)

func NewPersistency() PersistencyAdapter {
	return PersistencyAdapter{driver: NewDriverWrapperSqlite(config.GetSqlitePath())}

}

type DriverInterface interface {
	Open() (*sql.DB, error)
	Close() error
}

type DriverWrapperSqlite struct {
	path string
	db   *sql.DB
}

func NewDriverWrapperSqlite(path string) DriverInterface {
	return &DriverWrapperSqlite{path: path}
}

func (d *DriverWrapperSqlite) Open() (*sql.DB, error) {
	db, err := sql.Open("sqlite", d.path)
	if err != nil {
		d.db = db
	}
	return db, err
}

func (d *DriverWrapperSqlite) Close() error {
	return d.db.Close()
}

type Session struct {
	db       *sql.DB
	status   pp.RecordInterface[pp.StatusRecord]
	carousel pp.RecordInterface[pp.CarouselRecord]
}

func (s *Session) init() {
	s.carousel = &CarouselRecordAccessor{session: s}
	s.status = &StatusRecordAccessor{session: s}
}
func (s *Session) Close() error {
	return s.db.Close()
}

func (s *Session) Carousel() pp.RecordInterface[pp.CarouselRecord] {
	return s.carousel
}
func (s *Session) Status() pp.RecordInterface[pp.StatusRecord] {
	return s.status
}

type PersistencyAdapter struct {
	driver DriverInterface
}

func (p *PersistencyAdapter) Open() (pp.PersistencyInterface, error) {
	return p.open()
}

func (p *PersistencyAdapter) open() (*Session, error) {
	db, err := p.driver.Open()
	if err != nil {
		fmt.Printf("Error openning driver: %s\n", err)
		return nil, err
	}
	s := Session{db: db}
	s.init()
	return &s, nil
}
