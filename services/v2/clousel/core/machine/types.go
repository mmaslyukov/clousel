package machine

import (
	"time"

	"github.com/google/uuid"
)

type MachineStatus = string

const (
	MachineStatusInvalid MachineStatus = "invalid"
	MachineStatusNew     MachineStatus = "new"
	MachineStatusOnline  MachineStatus = "online"
	MachineStatusOffline MachineStatus = "offline"
	MachineStatusFailed  MachineStatus = "failed"
)

type GameStatus = string

const (
	GameStatusIvalid    MachineStatus = "invalid"
	GameStatusPending   MachineStatus = "pending"
	GameStatusCompleted MachineStatus = "completed"
	GameStatusFailed    MachineStatus = "failed"
)

type MachineEntry struct {
	MachId         uuid.UUID     `json:"MachId"`    //PK
	CompanyId      uuid.UUID     `json:"CompanyId"` //FK
	GameCost       int           `json:"Cost"`
	Status         MachineStatus `json:"Status"`
	Fee            int           `json:"Fee"`
	LastUpdateTime time.Time     `json:"LastUpdate"`
}

type GameEntry struct {
	EventId   uuid.UUID //PK (GameId)
	UserId    uuid.UUID //FK
	MachId    uuid.UUID
	GameCost  int
	Status    GameStatus
	StartTime time.Time
}

func (g *GameEntry) toGameEvent() *GameEvent {
	return &GameEvent{
		EventId: g.EventId,
		MachId:  g.MachId,
	}
}

type GameEvent struct {
	MachId  uuid.UUID `json:"MachId"`
	EventId uuid.UUID `json:"EventId"`
}
type GameEventAck struct {
	MachId        uuid.UUID `json:"MachId"`
	CorrelationId uuid.UUID `json:"CorrelationId"`
	Code          int       `json:"Code"`
}

type RemoteMachineStatus struct {
	MachId uuid.UUID     `json:"MachId"`
	Status MachineStatus `json:"Status"`
}
