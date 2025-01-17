package business

import (
	"clousel/lib/fault"

	"github.com/google/uuid"
)

type IBusinessRestController interface {
	Register(companyName string, email string, password string) fault.IError
	Login(companyName string, password string) (*BusinessEntry, fault.IError)
	AssignKeys(companyId uuid.UUID, skey string, prodId string) fault.IError
	ReadWhkey(companyId uuid.UUID) (string, fault.IError)
}
