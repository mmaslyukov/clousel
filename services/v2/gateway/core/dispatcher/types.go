package dispatcher

import "github.com/google/uuid"

type MachineStatus = string

const (
	MachineStatusInvalid MachineStatus = "invalid"
	MachineStatusNew     MachineStatus = "new"
	MachineStatusOnline  MachineStatus = "online"
	MachineStatusOffline MachineStatus = "offline"
	MachineStatusFailed  MachineStatus = "failed"
)

/*
NATS related part
*/
type GameEvent struct {
	MachId  uuid.UUID `json:"MachId"`
	EventId uuid.UUID `json:"EventId"`
}

type GameAckCode = int

const (
	GameEventAckOk          GameAckCode = 0
	GameEventAckPublishFail             = -1
	GameEventAckRemoteError             = -2
)

type GameEventAck struct {
	MachId        uuid.UUID `json:"MachId"`
	CorrelationId uuid.UUID `json:"CorrelationId"`
	Code          int       `json:"Code"`
}

type RemoteMachineStatus struct {
	MachId uuid.UUID     `json:"MachId"`
	Status MachineStatus `json:"Status"`
}

/*
MQTT related part
*/

const (
	CarouselEventMax = 50
)

const (
	MsgTypeEventHeartbeat = "Evt.Heartbeat"
	MsgTypeEventCompleted = "Evt.Completed"
	MsgTypeResponseAck    = "Res.Ack"
	MsgTypeRequestPlay    = "Req.Play"
)

type ITopic interface {
	Get() string
	PartOf(topic string) bool
	Contains(topic string) bool
	Parent() string
	Subscribable() string
	Appned(node string)
}

type IMessageGeneric interface {
	Target() uuid.UUID
	Name() string
	SetSequenceId(id int)
}

type MessageGeneric struct {
	MsgType string    `json:"Type"`
	CarId   uuid.UUID `json:"CarId"`
	SeqId   int       `json:"SeqNum"`
}

func (m *MessageGeneric) Name() string {
	return m.MsgType
}
func (m *MessageGeneric) Target() uuid.UUID {
	return m.CarId
}
func (m *MessageGeneric) SetSequenceId(id int) {
	m.SeqId = id
}

func CreateRequestPlay(CarId uuid.UUID, EvtId uuid.UUID) *RequestPlay {
	return &RequestPlay{MessageGeneric: MessageGeneric{MsgType: MsgTypeRequestPlay, CarId: CarId}, EvtId: EvtId}
}

type RequestPlay struct {
	MessageGeneric
	EvtId uuid.UUID `json:"EvtId"`
}

func (m *RequestPlay) Name() string {
	return m.MsgType
}
func (m *RequestPlay) Target() uuid.UUID {
	return m.CarId
}
func (m *RequestPlay) SetSequenceId(id int) {
	m.SeqId = id
}

type EventHeartbeat struct {
	MessageGeneric
	EvtId uuid.UUID `json:"EvtId"`
	Error string    `json:"Error"`
}

func (m *EventHeartbeat) Name() string {
	return m.MsgType
}
func (m *EventHeartbeat) Target() uuid.UUID {
	return m.CarId
}
func (m *EventHeartbeat) SetSequenceId(id int) {
	m.SeqId = id
}

type EventCompleted struct {
	MessageGeneric
	EvtId uuid.UUID `json:"EvtId"`
	Error string    `json:"Error"`
}

func (m *EventCompleted) Name() string {
	return m.MsgType
}
func (m *EventCompleted) Target() uuid.UUID {
	return m.CarId
}
func (m *EventCompleted) SetSequenceId(id int) {
	m.SeqId = id
}

type ResponseAck struct {
	MessageGeneric
	CorId uuid.UUID `json:"CorId"`
	Error string    `json:"Error"`
}

func (m *ResponseAck) Name() string {
	return m.MsgType
}
func (m *ResponseAck) Target() uuid.UUID {
	return m.CarId
}
func (m *ResponseAck) SetSequenceId(id int) {
	m.SeqId = id
}
