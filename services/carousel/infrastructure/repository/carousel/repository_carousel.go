package carousel

import (
	"carousel/core/manager"
	"carousel/core/operator"
	"carousel/infrastructure/repository/driver"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	table_carousel = "carousel-record"
	table_event    = "carousel-event"
)

type RepositoryCarousel struct {
	drv driver.IDBDriver
	log *zerolog.Logger
}

func New(drv driver.IDBDriver, log *zerolog.Logger) *RepositoryCarousel {
	return &RepositoryCarousel{drv: drv, log: log}
}

func (r *RepositoryCarousel) AddCarousel(c manager.Carousel) error {
	var prompt string
	prompt = fmt.Sprintf("insert into '%s' (CarouselId, OwnerId) values ('%s', '%s')", table_carousel, c.Cid, c.Oid)
	return r.drv.Session(func(db *sql.DB) error {
		var err error
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Send()
		} else {
			r.log.Err(err).Str("SQL", prompt).Str("CarouselId", c.Cid).Msg("Repository.Caorusel.AddCarousel: Fail to Register")
		}
		return err
	})
}
func (r *RepositoryCarousel) AddEventWithStatusNew(c manager.Carousel) error {
	var prompt string
	prompt = fmt.Sprintf("insert into '%s' (CarouselId, EventId, Rounds, Status) values ('%s', '%s', %d, '%s')", table_event, c.Cid, uuid.New().String(), 0, operator.CarouselStatusNameNew)
	return r.drv.Session(func(db *sql.DB) error {
		var err error
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Send()
		} else {
			r.log.Err(err).Str("SQL", prompt).Str("CarouselId", c.Cid).Msgf("Repository.Caorusel.AddEventWithStatusNew: Fail to Add New Status to '%s' table", table_event)
		}
		return err
	})
}
func (r *RepositoryCarousel) IsCarouselExistsInEvents(c manager.Carousel) (bool, error) {
	var err error
	exists := false
	prompt := fmt.Sprintf("select exists(select 1 from '%s' where CarouselId='%s' limit 1)", table_event, c.Cid)
	err = r.drv.Session(func(db *sql.DB) error {
		var err error
		if err = db.QueryRow(prompt).Scan(&exists); err != nil {
			r.log.Err(err).Str("SQL", prompt).Str("CarouselId", c.Cid).Msg("Repository.Caorusel.IsCarouselExistsInEvents: Fail to Query Carousel")
		}
		return err
	})
	return exists, err
}

// func (r *RepositoryCarousel) Register(c manager.Carousel) error {
// 	var err error
// 	for ok := true; ok; ok = false {
// 		if err = r.addCarousel(c); err != nil {
// 			break
// 		}
// 		if exists, e := r.isCarouselExistsInEvents(c); e != nil || exists {
// 			err = e
// 			break
// 		}
// 		if err = r.addEventWithStatusNew(c); err != nil {
// 			break
// 		}
// 	}
// 	return err
// }

func (r *RepositoryCarousel) Remove(c manager.Carousel) error {
	var prompt string
	if len(c.Cid) != 0 {
		prompt = fmt.Sprintf("delete from '%s' where CarouselId='%s'", table_carousel, c.Cid)
	} else if len(c.Oid) != 0 {
		prompt = fmt.Sprintf("delete from '%s' where OwnerId='%s'", table_carousel, c.Oid)
	}
	err := r.drv.Session(func(db *sql.DB) error {
		var err error
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Caorusel: Remove")
		} else {
			r.log.Err(err).Str("SQL", prompt).Str("CarouselId", c.Cid).Msg("Repository.Caorusel: Fail to Remove Carousel")
		}
		return err
	})
	return err
}

func (r *RepositoryCarousel) ReadOwned(ownerId string) ([]manager.Carousel, error) {
	var err error
	var recordArray []manager.Carousel
	prompt := fmt.Sprintf("select * from '%s' where OwnerId='%s'", table_carousel, ownerId)
	if err = r.drv.Session(func(db *sql.DB) error {
		var err error
		var rows *sql.Rows
		if rows, err = db.Query(prompt); err == nil {
			defer rows.Close()
			for rows.Next() {
				var c manager.Carousel
				if err := rows.Scan(&c.Cid, &c.Oid); err == nil {
					recordArray = append(recordArray, c)
				} else {
					r.log.Err(err).Str("OwnerId", ownerId).Msgf("Repository.Caorusel: Scan of '%s' failed", table_carousel)
				}
			}
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Caorusel: Read")
		}
		return err
	}); err != nil {
		r.log.Err(err).Str("SQL", prompt).Str("OwnerId", ownerId).Msgf("Repository.Caorusel: Fail to Read from '%s' table", table_carousel)
	}
	if recordArray == nil {
		err = fmt.Errorf("Cannot find record by OwnerId='%s'", ownerId)
	}
	return recordArray, err
}

func (r *RepositoryCarousel) IsExists(c manager.Carousel) (bool, error) {
	var err error
	exists := false
	prompt := fmt.Sprintf("select exists(select 1 from '%s' where CarouselId='%s' limit 1)", table_carousel, c.Cid)
	err = r.drv.Session(func(db *sql.DB) error {
		var err error
		if err = db.QueryRow(prompt).Scan(&exists); err != nil {
			r.log.Err(err).Str("CarouselId", c.Cid).Msg("Repository.Caorusel.IsCarouselExists: Fail to Query Carousel")
		}
		return err
	})
	return exists, err
}
