package core

import "github.com/google/uuid"

type IMessage interface {
	Name() string
	Id() uuid.UUID
}
