package snapshot

import (
	"carousel/core/operator"
	"carousel/infrastructure/repository/driver"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog"
)

const (
	table_snapshot = "carousel-snapshot"
)

type RepositorySnapshot struct {
	drv driver.IDBDriver
	log *zerolog.Logger
}

func New(drv driver.IDBDriver, log *zerolog.Logger) *RepositorySnapshot {
	return &RepositorySnapshot{drv: drv, log: log}
}

// "CarouselId" string UNIQUE NOT NULL,
// "Status" int NOT NULL,
// "Rounds" int NOT NULL,
// "Extra" string,

func (r *RepositorySnapshot) OperatorLoadSnapshot(carId string) (*operator.SnapshotData, error) {
	var err error
	var snapshot operator.SnapshotData
	prompt := fmt.Sprintf("select * from '%s' where CarouselId='%s'", table_snapshot, carId)
	if err = r.drv.Session(func(db *sql.DB) error {
		if err = db.QueryRow(prompt).Scan(&snapshot.CarId, &snapshot.Status, &snapshot.Rounds, &snapshot.Extra); err == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Snapshot.OperatorLoadSnapshot: Success")
		}
		return err
	}); err != nil {
		r.log.Err(err).Str("SQL", prompt).Str("CarouselId", carId).Msgf("Repository.Snapshot.OperatorLoadSnapshot: Fail to Read from '%s' table", table_snapshot)
	}
	return &snapshot, err
}

func (r *RepositorySnapshot) OperatorStoreSnapshot(sd *operator.SnapshotData) error {
	var prompt string
	var err error
	for ok := true; ok; ok = false {
		prompt = fmt.Sprintf("insert or ignore into '%s' (CarouselId, Status, Rounds) values ('%s', '%s', %d)", table_snapshot, sd.CarId, sd.Status, sd.Rounds)
		if err = r.drv.Session(func(db *sql.DB) error {
			if _, err = db.Exec(prompt); err == nil {
				r.log.Debug().Str("SQL", prompt).Send()
			} else {
				r.log.Err(err).Str("SQL", prompt).Str("CarouselId", sd.CarId).Msg("Repository.Snapshot.OperatorStoreSnapshot: Failure")
			}
			return err
		}); err != nil {
			break
		}
		prompt = fmt.Sprintf("update '%s' set Status='%s', Rounds=%d where CarouselId='%s'", table_snapshot, sd.Status, sd.Rounds, sd.CarId)
		if err = r.drv.Session(func(db *sql.DB) error {
			if _, err = db.Exec(prompt); err == nil {
				r.log.Debug().Str("SQL", prompt).Send()
			} else {
				r.log.Err(err).Str("SQL", prompt).Str("CarouselId", sd.CarId).Msg("Repository.Snapshot.OperatorStoreSnapshot: Failure")
			}
			return err
		}); err != nil {
			break
		}
	}
	return err
}

func (r *RepositorySnapshot) ManagerStoreNewSnapshot(carId string) error {
	return r.OperatorStoreSnapshot(&operator.SnapshotData{CarId: carId, Status: operator.CarouselStatusNameNew, Rounds: 0})
}

func (r *RepositorySnapshot) OperatorDeleteSnapshot(carId string) error {
	var err error
	prompt := fmt.Sprintf("delete from '%s' where CarouselId='%s'", table_snapshot, carId)
	err = r.drv.Session(func(db *sql.DB) error {
		var err error
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Msg("Repository.Snapshot.Delete: Success")
		} else {
			r.log.Err(err).Str("SQL", prompt).Str("CarouselId", carId).Msg("Repository.Snapshot.OperatorDeleteSnapshot: Fail to Remove Snapshot")
		}
		return err
	})
	return err
}
