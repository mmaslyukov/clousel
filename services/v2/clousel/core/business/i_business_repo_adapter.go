package business

import (
	"clousel/lib/fault"

	"github.com/google/uuid"
)

type IBusinessRepoAdapter interface {
	SaveNewBusinessEntry(companyId uuid.UUID, companyName string, email string, password string) fault.IError
	ReadBusinessEntryByName(companyName string) (*BusinessEntry, fault.IError)
	ReadBusinessEntryById(companyId uuid.UUID) (*BusinessEntry, fault.IError)
	AssignKeys(companyId uuid.UUID, skey string, prodId string, whid string, whkey string) fault.IError
}
