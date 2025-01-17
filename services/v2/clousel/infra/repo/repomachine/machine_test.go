package repomachine_test

import (
	"clousel/core/machine"
	"clousel/infra/log"
	"clousel/infra/repo"
	"clousel/infra/repo/driver"
	"clousel/infra/repo/repomachine"
	"clousel/infra/repo/repouser"
	"clousel/lib/fault"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
)

func CreateMachineTable(drv driver.IDBDriver) error {
	prompt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS '%s'(
		'Id' string PRIMARY KEY, 
		'CompanyId' string NOT NULL,
		'Cost' int  NOT NULL,
		'Status' string NOT NULL,
		'Updated' datetime)`, repomachine.TableMachine)
	return drv.Session(func(db *sql.DB) error {
		var err error
		_, err = db.Exec(prompt)
		return err
	})
}

func CreateGameHistoryTable(drv driver.IDBDriver) error {
	prompt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS '%s'(
		'Id' string PRIMARY KEY, 
		'MachId' string NOT NULL UNIQUE,
		'UserId' string NOT NULL UNIQUE,
		'GameCost' int NOT NULL,
		'Tickets' int NOT NULL,
		'Status' string NOT NULL,
		'Time' datetime)`, repomachine.TableGameHistory)
	return drv.Session(func(db *sql.DB) error {
		var err error
		_, err = db.Exec(prompt)
		return err
	})
}

// go test -v -run ^TestIMachineRepoGeneralAdapterInterface$  .\infra\repo\repomachine\
func TestIMachineRepoGeneralAdapterInterface(t *testing.T) {
	const dbPath = "TestIMachineRepoGeneralAdapterInterface.db"
	machId := uuid.New()
	companyId := uuid.New()
	log := log.New()
	var err fault.IError
	drv := repo.DriverSQLite.New(dbPath)
	repo := repo.Machine.New(drv, log)

	for ok := true; ok; ok = false {

		if err := CreateMachineTable(drv); err != nil {
			t.Errorf("Fail to create table '%s', err:%s", repouser.TableUser, err)
			break
		}
		if err = repo.SaveNewMachineEntry(machId, companyId, 2); err != nil {
			t.Errorf("Fail to create new machine entry: %s", err.Error())
			break
		}
		if err = repo.SaveNewMachineEntry(uuid.New(), companyId, 5); err != nil {
			t.Errorf("Fail to create new machine entry: %s", err.Error())
			break
		}
		var m1 *machine.MachineEntry
		if m1, err = repo.ReadMachineById(machId); err != nil {
			t.Errorf("Fail to read machine entry: %s", err.Error())
			break
		}

		if m1.CompanyId != companyId ||
			m1.GameCost != 2 ||
			m1.Status != machine.MachineStatusNew {
			t.Errorf("Client data is mismatch")
			break
		}
		if err = repo.UpdateMachineStatus(machId, machine.MachineStatusOnline); err != nil {
			t.Errorf("Fail to update machine entry: %s", err.Error())
			break
		}
		if err = repo.UpdateMachineCost(machId, 3); err != nil {
			t.Errorf("Fail to update machine entry: %s", err.Error())
			break
		}
		var m2 *machine.MachineEntry
		if m2, err = repo.ReadMachineById(machId); err != nil {
			t.Errorf("Fail to read machine entry: %s", err.Error())
			break
		}

		if m2.CompanyId != companyId ||
			m2.GameCost != 3 ||
			m2.Status != machine.MachineStatusOnline {
			t.Errorf("Client data is mismatch even after update")
			break
		}
		var ml []*machine.MachineEntry
		if ml, err = repo.ReadMachinesByCompanyId(companyId); err != nil {
			t.Errorf("Fail to read machine entry: %s", err.Error())
			break
		}
		if len(ml) != 2 {
			t.Errorf("Expected len is 2, but got %d", len(ml))
			break
		}
		// t.Logf("%+v", *ml[0])
		// t.Logf("%+v", *ml[1])

	}
	os.Remove(dbPath)
}

// go test -v -run ^IMachineRepoGameAdapterInterface$  .\infra\repo\repomachine\
func IMachineRepoGameAdapterInterface(t *testing.T) {
	const dbPath = "IMachineRepoGameAdapterInterface.db"
	gameId := uuid.New()
	machId := uuid.New()
	userId := uuid.New()
	log := log.New()
	var err fault.IError
	drv := repo.DriverSQLite.New(dbPath)
	repo := repo.Machine.New(drv, log)
	for ok := true; ok; ok = false {

		if err := CreateMachineTable(drv); err != nil {
			t.Errorf("Fail to create table '%s', err:%s", repouser.TableUser, err)
			break
		}
		if err = repo.SaveNewGameEntry(gameId, machId, userId, 1); err != nil {
			t.Errorf("Fail to create new machine entry: %s", err.Error())
			break
		}
		var g1 *machine.GameEntry
		if g1, err = repo.ReadGamesByGameId(gameId); err != nil {
			t.Errorf("Fail to read machine entry: %s", err.Error())
			break
		}
		if err = repo.UpdateGameStatus(gameId, machine.GameStatusCompleted); err != nil {
			t.Errorf("Fail to read machine entry: %s", err.Error())
			break
		}
		var g2 *machine.GameEntry
		if g2, err = repo.ReadGamesByGameId(gameId); err != nil {
			t.Errorf("Fail to read machine entry: %s", err.Error())
			break
		}
		var ga1 []*machine.GameEntry
		if ga1, err = repo.ReadGamesByUserId(userId); err != nil {
			t.Errorf("Fail to read machine entry: %s", err.Error())
			break
		}
		var ga2 []*machine.GameEntry
		if ga2, err = repo.ReadGamesByMachineId(machId); err != nil {
			t.Errorf("Fail to read machine entry: %s", err.Error())
			break
		}
		var ga3 []*machine.GameEntry
		if ga3, err = repo.ReadGamesByStatus(machine.GameStatusCompleted); err != nil {
			t.Errorf("Fail to read machine entry: %s", err.Error())
			break
		}
		if g1.UserId != userId || g2.UserId != g1.UserId || g2.Status != machine.GameStatusCompleted {
			t.Errorf("Game status is unexpected: %s", g2.Status)
			break
		}
		if ga1[0].EventId != ga2[0].EventId || ga1[0].EventId != ga3[0].EventId {
			t.Errorf("Game selectors are mismatch")
			break
		}
	}
	os.Remove(dbPath)
}
