package event

import (
	"carousel/core/operator"
	"carousel/infrastructure/repository/driver"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	table_event = "carousel-event"
)

type RepositoryEvent struct {
	drv    driver.IDBDriver
	log    *zerolog.Logger
	crRepo IRepositoryCarousel
}

func New(drv driver.IDBDriver, crRepo IRepositoryCarousel, log *zerolog.Logger) *RepositoryEvent {
	return &RepositoryEvent{drv: drv, crRepo: crRepo, log: log}
}

func (r *RepositoryEvent) ManagerStoreNewEvent(carId string) error {
	var err error
	prompt := fmt.Sprintf("insert into '%s' (CarouselId, EventId, Status, Tickets) values ('%s', '%s', '%s', %d)", table_event, carId, uuid.New().String(), operator.CarouselStatusNameNew, 0)
	err = r.drv.Session(func(db *sql.DB) error {
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Event.ManagerStoreNewEvent: Sucess")
		} else {
			r.log.Err(err).Str("SQL", prompt).Str("CarouselId", carId).Msg("Repository.Event.ManagerStoreNewEvent: Failure")
		}
		return err
	})
	return err
}

func (r *RepositoryEvent) OperatorRefill(rd *operator.TicketsData) error {
	var err error
	prompt := fmt.Sprintf("insert into '%s' (CarouselId, EventId, Tickets) values ('%s', '%s', %d)", table_event, rd.CarId, rd.EvtId, rd.Tickets)
	err = r.drv.Session(func(db *sql.DB) error {
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Event.OperatorRefill: Sucess")
		} else {
			r.log.Err(err).Str("SQL", prompt).Str("CarouselId", rd.CarId).Str("EventId", rd.EvtId.String()).Int("Tickets", rd.Tickets).Msg("Repository.Event.OperatorRefill: Failure")
		}
		return err
	})
	return err
}

func (r *RepositoryEvent) OperatorPlay(pd *operator.PlayData) error {
	var err error
	prompt := fmt.Sprintf("insert into '%s' (CarouselId, EventId, Tickets, Pending) values ('%s', '%s', %d, %d)", table_event, pd.CarId, pd.EvtId, -1, 1)
	err = r.drv.Session(func(db *sql.DB) error {
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Event.OperatorPlay: Success")
		} else {
			r.log.Err(err).Str("SQL", prompt).Str("CarouselId", pd.CarId).Str("EventId", pd.EvtId.String()).Msg("Repository.Event.OperatorPlay: Failure")
		}
		return err
	})
	return err
}

func (r *RepositoryEvent) OperatorRemove(carId string) error {
	var err error
	prompt := fmt.Sprintf("delete from '%s' where CarouselId='%s'", table_event, carId)
	err = r.drv.Session(func(db *sql.DB) error {
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Event.OperatorRemove: Success")
		} else {
			r.log.Err(err).Str("SQL", prompt).Str("CarouselId", carId).Msgf("Repository.Event.OperatorRemove: Fail to Remove entries from '%s' table", table_event)
		}
		return err
	})
	return err
}

func (r *RepositoryEvent) OperatorRemoveByEvent(evtId uuid.UUID) error {
	var err error
	prompt := fmt.Sprintf("delete from '%s' where EventId='%s'", table_event, evtId.String())
	err = r.drv.Session(func(db *sql.DB) error {
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Event.OperatorRemoveByEvent: Success")
		} else {
			r.log.Err(err).Str("SQL", prompt).Str("EventId", evtId.String()).Msgf("Repository.Event.OperatorRemoveByEvent: Fail to Remove entry from '%s' table", table_event)
		}
		return err
	})
	return err
}

func (r *RepositoryEvent) OperatorMark(s *operator.StatusData) error {
	var prompt string
	var err error
	for ok := true; ok; ok = false {
		if s == nil || s.Status == nil {
			err = fmt.Errorf("Invalid argumet s=%p, s.Status=%p", s, s.Status)
			break
		}
		var cd operator.CompositeData
		if cd, err = r.readLastEntry(&operator.Carousel{CarId: s.CarId}); err != nil {
			break
		}
		if s.Status == nil {
			err = fmt.Errorf("Provided Status is nil")
			break
		}

		log := r.log.Debug().Str("CarouselId", s.CarId)
		if s.Status != nil {
			log = log.Str("New Status", *s.Status)
		}
		if s.Error != nil {
			log = log.Str("New Error", *s.Error)
		}
		if cd.Status != nil {
			log = log.Str("Status", *cd.Status)
		}
		if cd.Error != nil {
			log = log.Str("Error", *cd.Status)
		}
		log.Msg("Repository.Event.OperatorMark: About to write")
		var evtId string

		if cd.Status != nil && *cd.Status == *s.Status {
			evtId = cd.EvtId.String()
			if *s.Status == operator.CarouselStatusNameOnline {
				prompt = fmt.Sprintf("update '%s' set Time=CURRENT_TIMESTAMP where CarouselId='%s' and EventId='%s'", table_event, cd.CarId, evtId)
			} else {
				r.log.Info().Str("CarouselId", s.CarId).Str("Status", *s.Status).Msg("Repository.Event.OperatorMark: Skip Status is abnormal, will skip time update")
				break
			}
		} else if s.Error != nil {
			evtId = s.EvtId.String()
			prompt = fmt.Sprintf("insert into '%s' (CarouselId, EventId, Tickets, Status, Error) values ('%s', '%s', %d, '%s', '%s')", table_event, s.CarId, evtId, 0, *s.Status, *s.Error)
		} else {
			evtId = s.EvtId.String()
			prompt = fmt.Sprintf("insert into '%s' (CarouselId, EventId, Tickets, Status) values ('%s', '%s', %d, '%s')", table_event, s.CarId, evtId, 0, *s.Status)
		}
		if len(prompt) == 0 {
			err = fmt.Errorf("Promt is empty")
			r.log.Err(err).Str("SQL", prompt).Str("CarouselId", s.CarId).Str("EventId", evtId).Str("Status", *s.Status).Msg("Repository.Event.OperatorMark: Failure")
			break
		}
		err = r.drv.Session(func(db *sql.DB) error {
			if _, err = db.Exec(prompt); err == nil {
				r.log.Debug().Str("SQL", prompt).Str("CarouselId", s.CarId).Str("EventId", evtId).Str("Status", *s.Status).Msg("Repository.Event.OperatorMark: Success")
			} else {
				r.log.Err(err).Str("SQL", prompt).Str("CarouselId", s.CarId).Str("EventId", evtId).Str("Status", *s.Status).Msg("Repository.Event.OperatorMark: Failure")
			}
			return err
		})
	}
	return err
}

func (r *RepositoryEvent) OperatorConfirm(s *operator.StatusData) error {
	var err error
	var prompt string
	if s.Error != nil {
		prompt = fmt.Sprintf("update '%s' set Status='%s', Error='%s', Pending=%d where CarouselId='%s' and EventId='%s", table_event, *s.Status, *s.Error, 0, s.CarId, s.EvtId)
	} else {
		prompt = fmt.Sprintf("update '%s' set Pending=%d where CarouselId='%s' and EventId='%s'", table_event, 0, s.CarId, s.EvtId)
	}
	err = r.drv.Session(func(db *sql.DB) error {
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Event.OperatorConfirm: Success")
		} else {
			r.log.Err(err).Str("SQL", prompt).Str("CarouselId", s.CarId).Str("EventId", s.EvtId.String()).Str("Status", *s.Status).Msg("Repository.Event.OperatorConfirm: Failure")
		}
		return err
	})
	return err
}

func (r *RepositoryEvent) OperatorUpdateTime(c *operator.Carousel) error {
	var err error
	var cd operator.CompositeData
	var prompt string
	if cd, err = r.readLastEntry(c); err == nil {
		prompt = fmt.Sprintf("update '%s' set Time=CURRENT_TIMESTAMP where CarouselId='%s' and EventId='%s'", table_event, c.CarId, cd.EvtId)
		err = r.drv.Session(func(db *sql.DB) error {
			if _, err = db.Exec(prompt); err == nil {
				r.log.Debug().Str("SQL", prompt).Msg("Repository.Event.OperatorUpdaeTime: Success")
			} else {
				r.log.Err(err).Str("SQL", prompt).Str("CarouselId", c.CarId).Str("EventId", cd.EvtId.String()).Msg("Repository.Event.OperatorUpdaeTime: Failure")
			}
			return err
		})
	}
	return err
}

func (r *RepositoryEvent) OperatorClearPendingFlag(ed *operator.EventData) error {
	var err error
	var prompt string
	prompt = fmt.Sprintf("update '%s' set Pending=null where CarouselId='%s' and EventId='%s'", table_event, ed.CarId, ed.EvtId)
	err = r.drv.Session(func(db *sql.DB) error {
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Event.OperatorClearPendingFlag: Success")
		} else {
			r.log.Err(err).Str("SQL", prompt).Str("CarouselId", ed.CarId).Str("EventId", ed.EvtId.String()).Msg("Repository.Event.OperatorClearPendingFlag: Failure")
		}
		return err
	})
	return err
}

func (r *RepositoryEvent) OperatorRead(carId string) ([]operator.CompositeData, error) {
	var err error
	var recordArray []operator.CompositeData
	prompt := fmt.Sprintf("select * from '%s' where CarouselId='%s'", table_event, carId)
	if err = r.drv.Session(func(db *sql.DB) error {
		var rows *sql.Rows
		if rows, err = db.Query(prompt); err == nil {
			defer rows.Close()
			for rows.Next() {
				var cd operator.CompositeData
				if err := rows.Scan(&cd.CarId, &cd.EvtId, &cd.Time, &cd.Status, &cd.Tickets, &cd.Pending, &cd.Error, &cd.Extra); err == nil {
					recordArray = append(recordArray, cd)
				} else {
					r.log.Err(err).Str("CarouselId", carId).Msgf("Repository.Event.OperatorRead: Scan of '%s' failed", table_event)
				}
			}
			r.log.Debug().Str("SQL", prompt).Str("CarouselId", carId).Msg("Repository.Event.OperatorRead: Success")
		} else {
			r.log.Err(err).Str("SQL", prompt).Str("CarouselId", carId).Msgf("Repository.Event.OperatorRead: Fail to Read from '%s' table", table_event)
		}
		return err
	}); err != nil {
		r.log.Err(err).Str("SQL", prompt).Str("CarouselId", carId).Msgf("Repository.Event.OperatorRead: Fail to Read from %s table", table_event)
	}
	return recordArray, err
}

func (r *RepositoryEvent) OperatorReadAsSnapshot(carId string) (*operator.SnapshotData, error) {
	var err error
	var snapshot operator.SnapshotData
	var recordArray []operator.CompositeData

	prompt := fmt.Sprintf("select * from '%s' where CarouselId='%s'", table_event, carId)
	if err = r.drv.Session(func(db *sql.DB) error {
		var rows *sql.Rows
		if rows, err = db.Query(prompt); err == nil {
			defer rows.Close()
			for rows.Next() {
				var cd operator.CompositeData
				if err := rows.Scan(&cd.CarId, &cd.EvtId, &cd.Time, &cd.Status, &cd.Tickets, &cd.Pending, &cd.Error, &cd.Extra); err == nil {
					recordArray = append(recordArray, cd)
				} else {
					r.log.Err(err).Str("CarouselId", carId).Msgf("Repository.Event.OperatorReadAsSnapshot: Scan of '%s' failed", table_event)
				}
			}
			r.log.Debug().Str("SQL", prompt).Str("CarouselId", carId).Msgf("Repository.Event.OperatorReadAsSnapshot: Success to Read from '%s' table", table_event)
		} else {
			r.log.Err(err).Str("SQL", prompt).Str("CarouselId", carId).Msgf("Repository.Event.OperatorReadAsSnapshot: Fail to Read from '%s' table", table_event)
		}
		return err
	}); err == nil {
		if len(recordArray) > 0 {
			for _, record := range recordArray {
				if len(snapshot.CarId) == 0 {
					snapshot.CarId = record.CarId
				}
				if record.Status != nil {
					snapshot.Status = *record.Status
				}
				if record.Extra != nil {
					snapshot.Extra = record.Error
				}
				snapshot.Tickets += record.Tickets
				// snapshot.Error = record.Error
				// r.log.Debug().Str("Time", record.Time).Str("CarouselId", c.CarId).Msgf("Status: %v, Error: %v", record.Status, record.Error)

			}
		} else {
			err = fmt.Errorf("Have no entries")
			// r.log.Warn().Str("SQL", prompt).Str("CarouselId", carId).Msgf("Repository.Event.OperatorReadAsSnapshot: Fail to Read from '%s' table", table_event)
		}
	} else {
		r.log.Err(err).Str("SQL", prompt).Str("CarouselId", carId).Msgf("Repository.Event.OperatorReadAsSnapshot: Fail to Read from '%s' table", table_event)
	}
	return &snapshot, err
}

func (r *RepositoryEvent) OperatorReadPendingTimeout(dur time.Duration) ([]operator.CompositeData, error) {
	var err error
	var recordArray []operator.CompositeData
	// select (julianday('now') - julianday(last_upd_dt)) * 24 * 60 * 60 as date_diff_seconds from mytable
	return recordArray, err
}

func (r *RepositoryEvent) OperatorReadPending() ([]operator.CompositeData, error) {
	var err error
	var recordArray []operator.CompositeData
	prompt := fmt.Sprintf("select * from '%s' where Pending=%d", table_event, 1)
	if err = r.drv.Session(func(db *sql.DB) error {
		var rows *sql.Rows
		if rows, err = db.Query(prompt); err == nil {
			defer rows.Close()
			for rows.Next() {
				var r operator.CompositeData
				if err := rows.Scan(&r.CarId, &r.EvtId, &r.Time, &r.Status, &r.Tickets, &r.Pending, &r.Error, &r.Extra); err == nil {
					recordArray = append(recordArray, r)
				}
			}
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Event.OperatorReadPendingAll: Success")
		} else {
			r.log.Err(err).Str("SQL", prompt).Msgf("Repository.Event.OperatorReadPendingAll: Fail to Read from '%s' table", table_event)
		}
		return err
	}); err != nil {
		r.log.Err(err).Str("SQL", prompt).Msgf("Repository.Event.OperatorReadPendingAll: Fail to Read from '%s' table", table_event)
	}
	return recordArray, err
}
func (r *RepositoryEvent) OperatorReadByStatus(status string) ([]operator.SnapshotData, error) {
	var err error
	var sdArray []operator.SnapshotData
	var carouselArray []string
	if carouselArray, err = r.crRepo.ReadCarouselsIds(); err == nil {
		for _, c := range carouselArray {
			var snapshot *operator.SnapshotData
			if snapshot, err = r.OperatorReadAsSnapshot(c); err == nil {
				// if snapshot.Status != operator.CarouselStatusNameOnline && snapshot.Status != operator.CarouselStatusNameNew {
				if snapshot.Status == status {
					sdArray = append(sdArray, *snapshot)
				}
			}
		}
	} else {
		r.log.Err(err).Msgf("Repository.Event.OperatorReadByStatus: Fail to Read")
	}
	return sdArray, err
}

// func (r *RepositoryEvent) IsExists(c operator.Carousel) (bool, error) {
// 	var err error
// 	exists := false
// 	prompt := fmt.Sprintf("select exists(select 1 from '%s' where CarouselId='%s' limit 1)", table_carousel, c.CarId)
// 	err = r.drv.Session(func(db *sql.DB) error {
// 		if err = db.QueryRow(prompt).Scan(&exists); err != nil {
// 			r.log.Err(err).Str("CarouselId", c.CarId).Msg("Repository.Caorusel.IsCarouselExists: Fail to Query Carousel")
// 		}
// 		return err
// 	})
// 	return exists, err
// }
//-----------------------------------------------
// func (r *RepositoryEvent) readCarousels() ([]manager.Carousel, error) {
// 	var err error
// 	var carouselArray []manager.Carousel
// 	prompt := fmt.Sprintf("select * from '%s'", table_carousel)
// 	if err = r.drv.Session(func(db *sql.DB) error {
// 		var rows *sql.Rows
// 		if rows, err = db.Query(prompt); err == nil {
// 			defer rows.Close()
// 			for rows.Next() {
// 				var c manager.Carousel
// 				if err := rows.Scan(&c.CarId, &c.OwnId); err == nil {
// 					carouselArray = append(carouselArray, c)
// 				} else {
// 					r.log.Err(err).Msgf("Repository.Event.readCarousels: Scan of '%s' failed", table_carousel)
// 				}
// 			}
// 			r.log.Debug().Str("SQL", prompt).Msg("Repository.Event.readCarousels: Success")
// 		} else {
// 			r.log.Err(err).Str("SQL", prompt).Msgf("Repository.Event.readCarousels: Fail to Read from '%s' table", table_event)
// 		}
// 		return err
// 	}); err != nil {
// 		r.log.Err(err).Str("SQL", prompt).Msgf("Repository.Event.readCarousels: Fail to Read from '%s' table", table_carousel)
// 	}
// 	return carouselArray, err
// }

func (r *RepositoryEvent) readLastEntry(c *operator.Carousel) (operator.CompositeData, error) {
	var err error
	var cd operator.CompositeData
	prompt := fmt.Sprintf("select * from '%s' where CarouselId='%s' order by Time desc limit 1", table_event, c.CarId)
	if err = r.drv.Session(func(db *sql.DB) error {
		if err = db.QueryRow(prompt).Scan(&cd.CarId, &cd.EvtId, &cd.Time, &cd.Status, &cd.Tickets, &cd.Pending, &cd.Error, &cd.Extra); err == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Event.readLastEntry: Success")
		}
		return err
	}); err != nil {
		r.log.Err(err).Str("SQL", prompt).Str("CarouselId", c.CarId).Msgf("Repository.Event.readLastEntry: Fail to Read from '%s' table", table_event)
	}
	return cd, err
}

func (r *RepositoryEvent) OperatorReadExpired(dur time.Duration) ([]operator.CompositeData, error) {
	var err error
	var cdArray []operator.CompositeData
	var carouselArray []string

	if carouselArray, err = r.crRepo.ReadCarouselsIds(); err == nil {
		for _, c := range carouselArray {
			var cdp *operator.CompositeData
			if cdp, err = r.readLastEntryExpired(&operator.Carousel{CarId: c}, dur); err == nil && cdp != nil {
				cdArray = append(cdArray, *cdp)
			}
		}
	} else {
		r.log.Err(err).Msg("Repository.Caorusel.OperatorReadExpired: Fail to Query")
	}
	return cdArray, err
}

// select * from (select * from 'carousel-event' where CarouselId='550e8400-e29b-41d4-a716-446655440000' order by Time desc limit 1) where Time<=datetime('now', '-8 minutes')
func (r *RepositoryEvent) readLastEntryExpired(c *operator.Carousel, dur time.Duration) (*operator.CompositeData, error) {
	var err error
	var cdp *operator.CompositeData
	from := fmt.Sprintf("select * from '%s' where CarouselId='%s' order by Time desc limit 1", table_event, c.CarId)
	prompt := fmt.Sprintf("select * from (%s) where Time<=datetime('now', '-%d seconds') order by Time desc limit 1", from, int(dur.Seconds()))
	if err = r.drv.Session(func(db *sql.DB) error {
		var rows *sql.Rows
		if rows, err = db.Query(prompt); err == nil {
			defer rows.Close()
			var cd operator.CompositeData
			for rows.Next() {
				if err := rows.Scan(&cd.CarId, &cd.EvtId, &cd.Time, &cd.Status, &cd.Tickets, &cd.Pending, &cd.Error, &cd.Extra); err == nil {
					cdp = &cd
					r.log.Debug().Str("SQL", prompt).Msg("Repository.Event.readLastEntryExpired: Success")
				} else {
					r.log.Err(err).Str("SQL", prompt).Msgf("Repository.Event.readLastEntryExpired: Fail to Scan from '%s' table", table_event)
				}
			}
		} else {
			r.log.Err(err).Str("SQL", prompt).Msgf("Repository.Event.readLastEntryExpired: Fail to Read from '%s' table", table_event)
		}
		return err
	}); err != nil {
		r.log.Err(err).Str("SQL", prompt).Str("CarouselId", c.CarId).Msgf("Repository.Event.readLastEntryExpired: Fail to Read from '%s' table", table_event)
	}
	return cdp, err
}
