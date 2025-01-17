package machine

import (
	"clousel/lib/fault"

	"github.com/google/uuid"
)

type IMachineClientAdapter interface {
	ApplyGameCost(eventId uuid.UUID, userId uuid.UUID, tickets int) fault.IError
	IsCanPay(userId uuid.UUID, cost int) (bool, fault.IError)
}
