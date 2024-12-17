package snapshot_test

import (
	"carousel/core/operator"
	"carousel/infrastructure/logger"
	"carousel/infrastructure/repository"
	"carousel/infrastructure/repository/snapshot"
	"os"
	"os/exec"
	"testing"
)

var dbPath = "test.db"

func pre(t *testing.T) {
	cmd := exec.Command("sqlite3", dbPath, ".read ../../../scripts/sqlite/carousel-tables-creat.sql")
	_, err := cmd.Output()
	if err != nil {
		t.Fatalf("Fail to create database: %s", err)
	}

}
func post(t *testing.T) {
	os.Remove(dbPath)
}

func TestTopicAppend(t *testing.T) {
	var err error
	var sdl *operator.SnapshotData
	log := logger.New()
	drv := repository.DriverSQLite.New(dbPath)
	sr := snapshot.New(drv, &log)
	sd := operator.SnapshotData{CarId: "Test", Status: "New", Tickets: 0}
	pre(t)
	_, err = sr.OperatorLoadSnapshot(sd.CarId)
	if err == nil {
		t.Errorf("Fail to unexisted snapshot")
	}
	err = sr.OperatorStoreSnapshot(&sd)
	if err != nil {
		t.Errorf("Fail to store snapshot: %s", err)
	}
	sdl, err = sr.OperatorLoadSnapshot(sd.CarId)
	if err != nil || sdl == nil {
		t.Errorf("Fail to load snapshot: %s", err)
	}
	if sdl.CarId != sd.CarId || sdl.Status != sd.Status {
		t.Errorf("Stored and loaded are mimatch: %v != %v", sd, *sdl)
	}

	sd.Status = "Updated"

	err = sr.OperatorStoreSnapshot(&sd)
	if err != nil {
		t.Errorf("Fail to store snapshot: %s", err)
	}
	sdl, err = sr.OperatorLoadSnapshot(sd.CarId)
	if err != nil || sdl == nil {
		t.Errorf("Fail to load snapshot: %s", err)
	}
	if sdl.CarId != sd.CarId || sdl.Status != sd.Status {
		t.Errorf("Stored and loaded are mimatch: %v != %v", sd, *sdl)
	}

	err = sr.OperatorDeleteSnapshot(sd.CarId)
	if err != nil {
		t.Errorf("Fail to delete snapshot: %s", err)
	}
	post(t)

}
