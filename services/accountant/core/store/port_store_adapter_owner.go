package store

import "github.com/google/uuid"

type IPortStoreAdapterOwner interface {
	ReadKeys(carId uuid.UUID) (pkey string, skey string, prodId string, err error)
}
