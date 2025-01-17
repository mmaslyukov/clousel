package business

import "github.com/google/uuid"

type BusinessEntry struct {
	CompanyId   uuid.UUID //PK
	CompanyName string
	Email       string
	Password    string
	ProdId      *string
	Skey        *string
	Whkey       *string
	Whid        *string
	Enabled     bool
}

func (b *BusinessEntry) Id() uuid.UUID {
	return b.CompanyId
}
