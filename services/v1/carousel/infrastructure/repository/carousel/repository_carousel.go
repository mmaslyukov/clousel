package carousel

import (
	"carousel/core/manager"
	"carousel/infrastructure/repository/driver"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog"
)

const (
	table_carousel = "carousel-record"
	// table_event    = "carousel-event"
	// table_snapshot = "carousel-snapshot"
)

type RepositoryCarousel struct {
	drv driver.IDBDriver
	log *zerolog.Logger
}

func New(drv driver.IDBDriver, log *zerolog.Logger) *RepositoryCarousel {
	return &RepositoryCarousel{drv: drv, log: log}
}

func (r *RepositoryCarousel) ManagerAddCarousel(c manager.Carousel) error {
	var prompt string
	prompt = fmt.Sprintf("insert into '%s' (CarouselId, OwnerId, Active) values ('%s', '%s', 1)", table_carousel, c.CarId, c.OwnId)
	return r.drv.Session(func(db *sql.DB) error {
		var err error
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Send()
		} else {
			r.log.Err(err).Str("SQL", prompt).Str("CarouselId", c.CarId).Msg("Repository.Caorusel.ManagerAddCarousel: Fail to Register")
		}
		return err
	})
}

// func (r *RepositoryCarousel) ManagerStoreNewSnapshot(carId string) error {
// 	var prompt string
// 	prompt = fmt.Sprintf("insert into '%s' (CarouselId, EventId, Tickets, Status) values ('%s', '%s', %d, '%s')", table_snapshot, carId, uuid.New().String(), 0, operator.CarouselStatusNameNew)
// 	return r.drv.Session(func(db *sql.DB) error {
// 		var err error
// 		if _, err = db.Exec(prompt); err == nil {
// 			r.log.Debug().Str("SQL", prompt).Send()
// 		} else {
// 			r.log.Err(err).Str("SQL", prompt).Str("CarouselId", carId).Msgf("Repository.Caorusel.ManagerStoreSnapshot: Fail to Add New Status to '%s' table", table_snapshot)
// 		}
// 		return err
// 	})
// }

// func (r *RepositoryCarousel) IsCarouselExistsInEvents(c manager.Carousel) (bool, error) {
// 	var err error
// 	exists := false
// 	prompt := fmt.Sprintf("select exists(select 1 from '%s' where CarouselId='%s' limit 1)", table_event, c.CarId)
// 	err = r.drv.Session(func(db *sql.DB) error {
// 		var err error
// 		if err = db.QueryRow(prompt).Scan(&exists); err != nil {
// 			r.log.Err(err).Str("SQL", prompt).Str("CarouselId", c.CarId).Msg("Repository.Caorusel.IsCarouselExistsInEvents: Fail to Query Carousel")
// 		}
// 		return err
// 	})
// 	return exists, err
// }

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

func (r *RepositoryCarousel) ManagerRemoveCarousel(carId string) error {
	var prompt string
	prompt = fmt.Sprintf("delete from '%s' where CarouselId='%s'", table_carousel, carId)
	err := r.drv.Session(func(db *sql.DB) error {
		var err error
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Caorusel.ManagerRemoveCarousel: Success")
		} else {
			r.log.Err(err).Str("SQL", prompt).Str("CarouselId", carId).Msg("Repository.Caorusel.ManagerRemoveCarousel: Fail to Remove Carousel")
		}
		return err
	})
	return err
}

func (r *RepositoryCarousel) ManagerRemoveOwner(ownerId string) error {
	var prompt string
	prompt = fmt.Sprintf("delete from '%s' where OwnerId='%s'", table_carousel, ownerId)
	err := r.drv.Session(func(db *sql.DB) error {
		var err error
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Caorusel.ManagerRemoveOwner: Success")
		} else {
			r.log.Err(err).Str("SQL", prompt).Str("OwnerId", ownerId).Msg("Repository.Caorusel.ManagerRemoveOwner: Fail to Remove Carousel")
		}
		return err
	})
	return err
}

func (r *RepositoryCarousel) ReadCarouselsIds() ([]string, error) {
	var err error
	var idArray []string
	prompt := fmt.Sprintf("select * from '%s'", table_carousel)
	if err = r.drv.Session(func(db *sql.DB) error {
		var rows *sql.Rows
		if rows, err = db.Query(prompt); err == nil {
			defer rows.Close()
			for rows.Next() {
				var c manager.Carousel
				if err := rows.Scan(&c.CarId, &c.OwnId, &c.Active); err == nil {
					idArray = append(idArray, c.CarId)
				} else {
					r.log.Err(err).Msgf("Repository.Carousel.ReadCarouselsIds: Scan of '%s' failed", table_carousel)
				}
			}
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Carousel.ReadCarouselsIds: Success")
		} else {
			r.log.Err(err).Str("SQL", prompt).Msgf("Repository.Carousel.ReadCarouselsIds: Fail to Read from '%s' table", table_carousel)
		}
		return err
	}); err != nil {
		r.log.Err(err).Str("SQL", prompt).Msgf("Repository.Carousel.ReadCarouselsIds: Fail to Read from '%s' table", table_carousel)
	}
	return idArray, err
}

func (r *RepositoryCarousel) ManagerReadOwnedCarousel(ownerId string) ([]manager.Carousel, error) {
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
				if err := rows.Scan(&c.CarId, &c.OwnId, &c.Active); err == nil {
					recordArray = append(recordArray, c)
				} else {
					r.log.Err(err).Str("OwnerId", ownerId).Msgf("Repository.Caorusel.ManagerReadOwnedCarousel: Scan of '%s' failed", table_carousel)
				}
			}
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Caorusel: Read")
		}
		return err
	}); err != nil {
		r.log.Err(err).Str("SQL", prompt).Str("OwnerId", ownerId).Msgf("Repository.Caorusel.ManagerReadOwnedCarousel: Fail to Read from '%s' table", table_carousel)
	}
	if recordArray == nil {
		err = fmt.Errorf("Cannot find record by OwnerId='%s'", ownerId)
	}
	return recordArray, err
}

func (r *RepositoryCarousel) isExists(CarId string) (bool, error) {
	var err error
	exists := false
	prompt := fmt.Sprintf("select exists(select 1 from '%s' where CarouselId='%s' limit 1)", table_carousel, CarId)
	err = r.drv.Session(func(db *sql.DB) error {
		var err error
		if err = db.QueryRow(prompt).Scan(&exists); err != nil {
			r.log.Err(err).Str("CarouselId", CarId).Msg("Repository.Caorusel.isExists: Fail to Query Carousel")
		}
		return err
	})
	return exists, err
}

func (r *RepositoryCarousel) ManagerIsExistsCarousel(carId string) (bool, error) {
	return r.isExists(carId)
}

func (r *RepositoryCarousel) OperatorIsExistsCarousel(carId string) (bool, error) {
	return r.isExists(carId)
}
func (r *RepositoryCarousel) OperarotReadAllCarouselIds() ([]string, error) {
	var err error
	var ids []string
	prompt := fmt.Sprintf("select * from '%s'", table_carousel)
	if err = r.drv.Session(func(db *sql.DB) error {
		var err error
		var rows *sql.Rows
		if rows, err = db.Query(prompt); err == nil {
			defer rows.Close()
			for rows.Next() {
				var c manager.Carousel
				if err := rows.Scan(&c.CarId, &c.OwnId, &c.Active); err == nil {
					ids = append(ids, c.CarId)
				} else {
					r.log.Err(err).Msgf("Repository.Caorusel.OperarotReadAllCarouselIds: Scan of '%s' failed", table_carousel)
				}
			}
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Caorusel.OperarotReadAllCarouselIds: Read")
		}
		return err
	}); err != nil {
		r.log.Err(err).Str("SQL", prompt).Msgf("Repository.Caorusel.OperarotReadAllCarouselIds: Fail to Read from '%s' table", table_carousel)
	}
	if ids == nil {
		err = fmt.Errorf("Cannot find any record in '%s' table", table_carousel)
	}
	return ids, err
}
