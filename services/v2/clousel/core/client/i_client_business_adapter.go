package client

import (
	"clousel/lib/fault"
)

type IClientBusinessAdapter interface {
	ClientReadKeys(companyName string) (string, string, fault.IError)
	IsCompanyExists(companyName string) (bool, fault.IError)
	// ClientReadId(companyName string) (uuid.UUID, fault.IError)
}
