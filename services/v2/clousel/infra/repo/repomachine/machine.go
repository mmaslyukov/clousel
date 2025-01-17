package repomachine

import (
	"clousel/core/machine"
	"clousel/infra/repo/driver"
	rec "clousel/infra/repo/errors"
	"clousel/infra/repo/types"
	"clousel/lib/fault"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	TableMachine     = "machine"
	TableGameHistory = "game"
)

type MachineColumn struct {
	MachId    types.Named[types.UUIDString]
	CompanyId types.Named[types.UUIDString]
	GameCost  types.Named[int]
	Status    types.Named[string]
	Fee       types.Named[int]
	Updated   types.Named[types.TimeString]
}

func MachineColumnDefault() MachineColumn {
	prof := MachineColumn{
		MachId:    types.NamedCreateDefault[types.UUIDString]("Id"),
		CompanyId: types.NamedCreateDefault[types.UUIDString]("CompanyId"),
		GameCost:  types.NamedCreateDefault[int]("Cost"),
		Status:    types.NamedCreateDefault[string]("Status"),
		Fee:       types.NamedCreateDefault[int]("Fee"),
		Updated:   types.NamedCreateDefault[types.TimeString]("Updated"),
	}
	return prof
}

func (c *MachineColumn) toEntry() *machine.MachineEntry {
	entry := &machine.MachineEntry{
		MachId:         c.MachId.ValuePtr().Uuid(),
		CompanyId:      c.CompanyId.ValuePtr().Uuid(),
		GameCost:       c.GameCost.Value(),
		Status:         c.Status.Value(),
		Fee:            c.Fee.Value(),
		LastUpdateTime: c.Updated.ValuePtr().Time(),
	}
	return entry
}

type GameColumn struct {
	GameId   types.Named[types.UUIDString]
	MachId   types.Named[types.UUIDString]
	UserId   types.Named[types.UUIDString]
	GameCost types.Named[int]
	Status   types.Named[string]
	Time     types.Named[types.TimeString]
}

func GameColumnDefault() GameColumn {
	prof := GameColumn{
		GameId:   types.NamedCreateDefault[types.UUIDString]("Id"),
		MachId:   types.NamedCreateDefault[types.UUIDString]("MachId"),
		UserId:   types.NamedCreateDefault[types.UUIDString]("UserId"),
		GameCost: types.NamedCreateDefault[int]("Cost"),
		Status:   types.NamedCreateDefault[string]("Status"),
		Time:     types.NamedCreateDefault[types.TimeString]("Started"),
	}
	return prof
}

func (c *GameColumn) toEntry() *machine.GameEntry {
	entry := &machine.GameEntry{
		EventId:   c.GameId.ValuePtr().Uuid(),
		MachId:    c.MachId.ValuePtr().Uuid(),
		UserId:    c.UserId.ValuePtr().Uuid(),
		GameCost:  c.GameCost.Value(),
		Status:    c.Status.Value(),
		StartTime: c.Time.ValuePtr().Time(),
	}
	return entry
}

type RepositoryMachine struct {
	drv driver.IDBDriver
	log *zerolog.Logger
}

func RepositoryMachineCreate(drv driver.IDBDriver, log *zerolog.Logger) *RepositoryMachine {
	return &RepositoryMachine{drv: drv, log: log}
}

/*
IMachineRepoGeneralAdapter
*/
func (r *RepositoryMachine) SaveNewMachineEntry(machId uuid.UUID, companyId uuid.UUID, cost int) fault.IError {
	const fn = "Repository.Machine.SaveNewMachineEntry"
	var err fault.IError
	t := MachineColumnDefault()

	t.MachId.ValuePtr().SetUuid(machId)
	t.CompanyId.ValuePtr().SetUuid(companyId)
	t.GameCost.SetValue(cost)
	t.Status.SetValue(machine.MachineStatusNew)
	t.Fee.SetValue(0)
	t.Updated.ValuePtr().SetTime(time.Now())
	prompt := fmt.Sprintf("insert into '%s' (%s,%s,%s,%s,%s,%s) values ('%s','%s',%d,'%s',%d,'%s')", TableMachine,
		t.MachId.Name(), t.CompanyId.Name(), t.GameCost.Name(), t.Status.Name(), t.Fee.Name(), t.Updated.Name(),
		t.MachId.ValuePtr().Str(), t.CompanyId.ValuePtr().Str(), t.GameCost.Value(), t.Status.Value(), t.Fee.Value(), t.Updated.ValuePtr().Str())
	if e := r.drv.Session(func(db *sql.DB) error {
		var err error
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Send()
		} else {
			r.log.Err(err).Str("SQL", prompt).Msgf("%s: Fail to Add", fn)
		}
		return err
	}); e != nil {
		err = fault.New(rec.ERepoExecPrompt).Msg(e.Error())
	}
	return err
}

func (r *RepositoryMachine) UpdateMachineStatus(machId uuid.UUID, status machine.MachineStatus) fault.IError {
	const fn = "Repository.Machine.UpdateMachineStatus"
	t := MachineColumnDefault()
	t.Status.SetValue(status)
	t.Updated.ValuePtr().SetTime(time.Now())
	return r.updateMachineWith(machId, func() string {
		return fmt.Sprintf("%s='%s', %s='%s'",
			t.Updated.Name(), t.Updated.ValuePtr().Str(),
			t.Status.Name(), t.Status.Value())
	})
}

func (r *RepositoryMachine) UpdateMachineCost(machId uuid.UUID, cost int) fault.IError {
	const fn = "Repository.Machine.UpdateMachineCost"
	t := MachineColumnDefault()
	t.GameCost.SetValue(cost)
	return r.updateMachineWith(machId, func() string {
		return fmt.Sprintf("%s=%d",
			t.GameCost.Name(), t.GameCost.Value())
	})
}

func (r *RepositoryMachine) UpdateFee(machId uuid.UUID, fee int) fault.IError {
	const fn = "Repository.Machine.UpdateFee"
	t := MachineColumnDefault()
	t.Fee.SetValue(fee)
	return r.updateMachineWith(machId, func() string {
		return fmt.Sprintf("%s=%d",
			t.Fee.Name(), t.Fee.Value())
	})
}

func (r *RepositoryMachine) ReadMachineById(machId uuid.UUID) (*machine.MachineEntry, fault.IError) {
	const fn = "Repository.Machine.ReadMachineById"
	t := MachineColumnDefault()
	t.MachId.ValuePtr().SetUuid(machId)
	filter := fmt.Sprintf("%s='%s'", t.MachId.Name(), t.MachId.ValuePtr().Str())
	// r.log.Info().Msg(filter)
	if rows, err := r.readMachineEntriesBy(func() string {
		return filter
	}); err != nil {
		return nil, err
	} else if len(rows) != 1 {
		return nil, fault.New(rec.ERepoUnexpectedEntriesCount).Msgf("Got %d rows, by filter '%s'", len(rows), filter)
	} else {
		return rows[0].toEntry(), nil
	}

	// var entries []*machine.MachineEntry
	// if err == nil {
	// 	for _, r := range rows {
	// 		entries = append(entries, r.toEntry())
	// 	}
	// }
}
func (r *RepositoryMachine) ReadMachinesByCompanyId(companyId uuid.UUID) ([]*machine.MachineEntry, fault.IError) {
	const fn = "Repository.Machine.ReadMachinesByCompanyId"
	t := MachineColumnDefault()
	t.CompanyId.ValuePtr().SetUuid(companyId)
	filter := fmt.Sprintf("%s='%s'", t.CompanyId.Name(), t.CompanyId.ValuePtr().Str())
	rows, err := r.readMachineEntriesBy(func() string {
		return filter
	})
	var entries []*machine.MachineEntry
	if err == nil {
		for _, r := range rows {
			// fmt.Printf("time %v\n", r.Updated.ValuePtr().Str())
			entries = append(entries, r.toEntry())
		}
	}
	return entries, err
}
func (r *RepositoryMachine) ReadMachinesByStatus(companyId uuid.UUID, status machine.MachineStatus) ([]*machine.MachineEntry, fault.IError) {
	const fn = "Repository.Machine.ReadMachinesByStatus"
	t := MachineColumnDefault()
	t.Status.SetValue(status)
	t.CompanyId.ValuePtr().SetUuid(companyId)
	filter := fmt.Sprintf("%s='%s' and %s='%s'",
		t.Status.Name(), t.Status.Value(),
		t.CompanyId.Name(), t.CompanyId.ValuePtr().Str())
	rows, err := r.readMachineEntriesBy(func() string {
		return filter
	})
	var entries []*machine.MachineEntry
	if err == nil {
		for _, r := range rows {
			entries = append(entries, r.toEntry())
		}
	}
	return entries, err
}

/*
IMachineRepoGameAdapter
*/
func (r *RepositoryMachine) SaveNewGameEntry(gameId uuid.UUID, machId uuid.UUID, userId uuid.UUID, cost int) fault.IError {
	const fn = "Repository.Machine.SaveNewGameEntry"
	var err fault.IError
	t := GameColumnDefault()
	t.GameId.ValuePtr().SetUuid(gameId)
	t.MachId.ValuePtr().SetUuid(machId)
	t.UserId.ValuePtr().SetUuid(userId)
	t.GameCost.SetValue(cost)
	t.Status.SetValue(machine.GameStatusPending)
	t.Time.ValuePtr().SetTime(time.Now())
	prompt := fmt.Sprintf("insert into '%s' (%s,%s,%s,%s,%s,%s) values ('%s','%s','%s',%d,'%s','%s')", TableGameHistory,
		t.GameId.Name(), t.MachId.Name(), t.UserId.Name(), t.GameCost.Name(), t.Status.Name(), t.Time.Name(),
		t.GameId.ValuePtr().Str(), t.MachId.ValuePtr().Str(), t.UserId.ValuePtr().Str(), t.GameCost.Value(), t.Status.Value(), t.Time.ValuePtr().Str(),
	)
	if e := r.drv.Session(func(db *sql.DB) error {
		var err error
		if _, err = db.Exec(prompt); err == nil {
			r.log.Debug().Str("SQL", prompt).Send()
		} else {
			r.log.Err(err).Str("SQL", prompt).Msgf("%s: Fail to Add", fn)
		}
		return err
	}); e != nil {
		err = fault.New(rec.ERepoExecPrompt).Msg(e.Error())
	}
	return err
}
func (r *RepositoryMachine) UpdateGameStatus(gameId uuid.UUID, status machine.GameStatus) fault.IError {
	const fn = "Repository.Machine.UpdateGameStatus"
	t := GameColumnDefault()
	t.Status.SetValue(status)
	return r.updateGameWith(gameId, func() string {
		return fmt.Sprintf("%s='%s'",
			t.Status.Name(), t.Status.Value())
	})
}
func (r *RepositoryMachine) ReadGamesByGameId(gameId uuid.UUID) (*machine.GameEntry, fault.IError) {
	const fn = "Repository.Machine.ReadGamesByGameId"
	t := GameColumnDefault()
	t.GameId.ValuePtr().SetUuid(gameId)
	filter := fmt.Sprintf("%s='%s'", t.GameId.Name(), t.GameId.ValuePtr().Str())
	if rows, err := r.readGameHistoryEntriesBy(func() string {
		return filter
	}); err != nil {
		return nil, err
	} else if len(rows) != 1 {
		return nil, fault.New(rec.ERepoUnexpectedEntriesCount).Msgf("Got %d rows, by filter '%s'", len(rows), filter)
	} else {
		return rows[0].toEntry(), nil
	}
}
func (r *RepositoryMachine) ReadGamesByMachineId(machId uuid.UUID) ([]*machine.GameEntry, fault.IError) {
	const fn = "Repository.Machine.ReadGamesByMachineId"
	t := GameColumnDefault()
	t.MachId.ValuePtr().SetUuid(machId)
	filter := fmt.Sprintf("%s='%s'", t.MachId.Name(), t.MachId.ValuePtr().Str())
	rows, err := r.readGameHistoryEntriesBy(func() string {
		return filter
	})
	var entries []*machine.GameEntry
	if err == nil {
		for _, r := range rows {
			entries = append(entries, r.toEntry())
		}
	}
	return entries, err
	// panic(fmt.Sprintf("%s not implemented", fn))
}

func (r *RepositoryMachine) ReadGamesByMachineIdWithStatus(machId uuid.UUID, status machine.GameStatus) ([]*machine.GameEntry, fault.IError) {
	const fn = "Repository.Machine.ReadGamesByMachineIdWithStatus"
	t := GameColumnDefault()
	t.MachId.ValuePtr().SetUuid(machId)
	t.Status.SetValue(status)
	filter := fmt.Sprintf("%s='%s' and  %s='%s'",
		t.MachId.Name(), t.MachId.ValuePtr().Str(),
		t.Status.Name(), t.Status.Value())
	rows, err := r.readGameHistoryEntriesBy(func() string {
		return filter
	})
	var entries []*machine.GameEntry
	if err == nil {
		for _, r := range rows {
			entries = append(entries, r.toEntry())
		}
	}
	return entries, err
}

func (r *RepositoryMachine) ReadGamesByUserId(userId uuid.UUID) ([]*machine.GameEntry, fault.IError) {
	const fn = "Repository.Machine.ReadGamesByUserId"
	t := GameColumnDefault()
	t.UserId.ValuePtr().SetUuid(userId)
	filter := fmt.Sprintf("%s='%s'", t.UserId.Name(), t.UserId.ValuePtr().Str())
	rows, err := r.readGameHistoryEntriesBy(func() string {
		return filter
	})
	var entries []*machine.GameEntry
	if err == nil {
		for _, r := range rows {
			entries = append(entries, r.toEntry())
		}
	}
	return entries, err
	// panic(fmt.Sprintf("%s not implemented", fn))
}
func (r *RepositoryMachine) ReadGamesByStatus(status machine.GameStatus) ([]*machine.GameEntry, fault.IError) {
	const fn = "Repository.Machine.ReadGamesByStatus"
	t := GameColumnDefault()
	t.Status.SetValue(status)
	filter := fmt.Sprintf("%s='%s'", t.Status.Name(), t.Status.Value())
	rows, err := r.readGameHistoryEntriesBy(func() string {
		return filter
	})
	var entries []*machine.GameEntry
	if err == nil {
		for _, r := range rows {
			entries = append(entries, r.toEntry())
		}
	}
	return entries, err
	// panic(fmt.Sprintf("%s not implemented", fn))
}

func (r *RepositoryMachine) ReadGamesByStatusAndTime(status machine.GameStatus, ts time.Time) ([]*machine.GameEntry, fault.IError) {
	const fn = "Repository.Machine.ReadGamesByStatusAndTime"
	t := GameColumnDefault()
	t.Status.SetValue(status)
	t.Time.ValuePtr().SetTime(ts)
	filter := fmt.Sprintf("%s='%s' and %s < '%s'",
		t.Status.Name(), t.Status.Value(),
		t.Time.Name(), t.Time.ValuePtr().Str())
	rows, err := r.readGameHistoryEntriesBy(func() string {
		return filter
	})
	var entries []*machine.GameEntry
	if err == nil {
		for _, r := range rows {
			entries = append(entries, r.toEntry())
		}
	}
	return entries, err
	// panic(fmt.Sprintf("%s not implemented", fn))
}

/*
Private methods
*/

func (r *RepositoryMachine) updateMachineWith(machId uuid.UUID, set func() string) fault.IError {
	const fn = "Repository.Machine.updateMachineWith"
	var err fault.IError
	t := MachineColumnDefault()
	t.MachId.ValuePtr().SetUuid(machId)
	prompt := fmt.Sprintf("update '%s' set %s where %s='%s'", TableMachine,
		set(),
		t.MachId.Name(), t.MachId.ValuePtr().Str())

	if e := r.drv.Session(func(db *sql.DB) error {
		var e error
		if _, e = db.Exec(prompt); e == nil {
			r.log.Debug().Str("SQL", prompt).Msgf("%s: Success", fn)
		} else {
			r.log.Err(e).Str("SQL", prompt).Msgf("%s: Failure", fn)
		}
		return e
	}); e != nil {
		err = fault.New(rec.ERepoExecPrompt).Msgf("Fail to update machine %s", e.Error())
	}
	return err
}

func (r *RepositoryMachine) updateGameWith(gameId uuid.UUID, set func() string) fault.IError {
	const fn = "Repository.Machine.updateMachineWith"
	var err fault.IError
	t := GameColumnDefault()
	t.GameId.ValuePtr().SetUuid(gameId)
	prompt := fmt.Sprintf("update '%s' set %s where %s='%s'", TableGameHistory,
		set(),
		t.GameId.Name(), t.GameId.ValuePtr().Str())

	if e := r.drv.Session(func(db *sql.DB) error {
		var e error
		if _, e = db.Exec(prompt); e == nil {
			r.log.Debug().Str("SQL", prompt).Msgf("%s: Success", fn)
		} else {
			r.log.Err(e).Str("SQL", prompt).Msgf("%s: Failure", fn)
		}
		return e
	}); e != nil {
		err = fault.New(rec.ERepoExecPrompt).Msgf("Fail to update machine %s", e.Error())
	}
	return err
}

func (r *RepositoryMachine) readMachineEntriesBy(where func() string) ([]*MachineColumn, fault.IError) {
	const fn = "Repository.User.readMachineEntriesBy"
	var err fault.IError
	var e error
	var entries []*MachineColumn
	prompt := fmt.Sprintf("select * from '%s' where %s", TableMachine, where())
	if e = r.drv.Session(func(db *sql.DB) error {
		var rows *sql.Rows
		if rows, e = db.Query(prompt); e == nil {
			defer rows.Close()
			for rows.Next() {
				t := MachineColumnDefault()
				if e := rows.Scan(
					t.MachId.ValuePtr().Ptr(),
					t.CompanyId.ValuePtr().Ptr(),
					t.GameCost.ValuePtr(),
					t.Status.ValuePtr(),
					t.Fee.ValuePtr(),
					t.Updated.ValuePtr().Ptr(),
				); e == nil {
					entries = append(entries, &t)
				} else {
					r.log.Err(e).Msgf("%s:", fn)
				}
			}
			r.log.Debug().Str("SQL", prompt).Msgf("%s: Success", fn)
		} else {
			r.log.Err(e).Str("SQL", prompt).Msgf("%s: Fail to Read from '%s' table", fn, TableMachine)
		}
		return e
	}); e != nil {
		err = fault.New(rec.ERepoQueryData).Msg(e.Error())
		r.log.Err(e).Str("SQL", prompt).Msgf("%s: Fail to Read from '%s' table", fn, TableMachine)
	}
	return entries, err
}

func (r *RepositoryMachine) readGameHistoryEntriesBy(where func() string) ([]*GameColumn, fault.IError) {
	const fn = "Repository.User.readGameHistoryEntriesBy"
	var err fault.IError
	var e error
	var entries []*GameColumn
	prompt := fmt.Sprintf("select * from '%s' where %s", TableGameHistory, where())
	if e = r.drv.Session(func(db *sql.DB) error {
		var rows *sql.Rows
		if rows, e = db.Query(prompt); e == nil {
			defer rows.Close()
			for rows.Next() {
				t := GameColumnDefault()
				if e := rows.Scan(
					t.GameId.ValuePtr().Ptr(),
					t.MachId.ValuePtr().Ptr(),
					t.UserId.ValuePtr().Ptr(),
					t.GameCost.ValuePtr(),
					t.Status.ValuePtr(),
					t.Time.ValuePtr().Ptr(),
				); e == nil {
					entries = append(entries, &t)
				} else {
					r.log.Err(e).Msgf("%s:", fn)
				}
			}
			r.log.Debug().Str("SQL", prompt).Msgf("%s: Success", fn)
		} else {
			r.log.Err(e).Str("SQL", prompt).Msgf("%s: Fail to Read from '%s' table", fn, TableMachine)
		}
		return e
	}); e != nil {
		err = fault.New(rec.ERepoQueryData).Msg(e.Error())
	}
	return entries, err
}
