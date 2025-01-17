package machine

import (
	"clousel/lib/fault"
	"time"

	"github.com/google/uuid"
)

type IMachineSelector interface {
	CompanyId() *uuid.UUID
	MachId() *uuid.UUID
	Status() *MachineStatus
	TimeFrom() *time.Time
	TimeTill() *time.Time
}

type IMachineRestController interface {
	SaveNewMachineEntry(machId uuid.UUID, companyId uuid.UUID, gameCost int) fault.IError
	ReadMachineEntriesBySelector(selector IMachineSelector) (entries []*MachineEntry, err fault.IError)
	ChangeGameCost(machId uuid.UUID, gameCost int) fault.IError
	ChangeSubscriptionFee(machId uuid.UUID, fee int) fault.IError

	PlayRequest(machId uuid.UUID, userId uuid.UUID) (eventId *uuid.UUID, err fault.IError)
	PollRequestStatus(eventId uuid.UUID) (GameStatus, fault.IError)
}
