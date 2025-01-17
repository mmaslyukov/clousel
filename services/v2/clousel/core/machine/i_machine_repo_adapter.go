package machine

import (
	"clousel/lib/fault"
	"time"

	"github.com/google/uuid"
)

type IMachineRepoGeneralAdapter interface {
	SaveNewMachineEntry(machId uuid.UUID, companyId uuid.UUID, cost int) fault.IError
	UpdateMachineStatus(machId uuid.UUID, status MachineStatus) fault.IError
	UpdateMachineCost(machId uuid.UUID, cost int) fault.IError
	UpdateFee(machId uuid.UUID, fee int) fault.IError
	ReadMachineById(machId uuid.UUID) (*MachineEntry, fault.IError)
	ReadMachinesByCompanyId(companyId uuid.UUID) ([]*MachineEntry, fault.IError)
	ReadMachinesByStatus(companyId uuid.UUID, status MachineStatus) ([]*MachineEntry, fault.IError)
}

type IMachineRepoGameAdapter interface {
	SaveNewGameEntry(gameId uuid.UUID, machId uuid.UUID, userId uuid.UUID, cost int) fault.IError
	UpdateGameStatus(gameId uuid.UUID, status GameStatus) fault.IError
	ReadGamesByGameId(gameId uuid.UUID) (*GameEntry, fault.IError)
	ReadGamesByMachineId(machId uuid.UUID) ([]*GameEntry, fault.IError)
	ReadGamesByMachineIdWithStatus(machId uuid.UUID, status GameStatus) ([]*GameEntry, fault.IError)
	ReadGamesByUserId(userId uuid.UUID) ([]*GameEntry, fault.IError)
	ReadGamesByStatus(status GameStatus) ([]*GameEntry, fault.IError)
	ReadGamesByStatusAndTime(status GameStatus, ts time.Time) ([]*GameEntry, fault.IError)
}
