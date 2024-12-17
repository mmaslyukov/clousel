package operator

import "github.com/google/uuid"

const (
	CarouselStatusNameOnline  = "online"
	CarouselStatusNameOffline = "offline"
	CarouselStatusNameFailure = "failure"
	CarouselStatusNameNew     = "new"
)

const (
	CarouselEventMax = 50
)

type Carousel struct {
	CarId string
}

type EventData struct {
	CarId string `json:"CarouselId"`
	EvtId uuid.UUID
}
type PlayData = EventData

type TicketsData struct {
	CarId  string `json:"CarouselId"`
	Tickets int
	EvtId  uuid.UUID
}

type StatusData struct {
	CarId  string `json:"CarouselId"`
	EvtId  uuid.UUID
	Status *string
	Error  *string
}

type CompositeData struct {
	CarId   string `json:"CarouselId"`
	EvtId   uuid.UUID
	Time    string
	Status  *string
	Tickets  int
	Pending *int
	Error   *string
	Extra   *string
}

type SnapshotData struct {
	CarId  string `json:"CarouselId"`
	Status string
	Tickets int
	Extra  *string
}

const (
	MsgTypeEventHeartbeat = "Evt.Heartbeat"
	MsgTypeEventCompleted = "Evt.Completed"
	MsgTypeResponseAck    = "Res.Ack"
	MsgTypeRequestPlay    = "Req.Play"
)

type Itopic interface {
	Get() string
	PartOf(topic string) bool
	Contains(topic string) bool
	Parent() string
	Subscribable() string
	Appned(node string)
}

type IMessageGeneric interface {
	Target() string
	Name() string
	SetSequenceId(id int)
}

type MessageGeneric struct {
	MsgType string `json:"Type"`
	CarId   string `json:"CarId"`
	SeqId   int    `json:"SeqNum"`
}

func (m *MessageGeneric) Name() string {
	return m.MsgType
}
func (m *MessageGeneric) Target() string {
	return m.CarId
}
func (m *MessageGeneric) SetSequenceId(id int) {
	m.SeqId = id
}

func CreateRequestPlay(CarId string, EvtId string) *RequestPlay {
	return &RequestPlay{MessageGeneric: MessageGeneric{MsgType: MsgTypeRequestPlay, CarId: CarId}, EvtId: EvtId}
}

type RequestPlay struct {
	MessageGeneric
	EvtId string `json:"EvtId"`
}

func (m *RequestPlay) Name() string {
	return m.MsgType
}
func (m *RequestPlay) Target() string {
	return m.CarId
}
func (m *RequestPlay) SetSequenceId(id int) {
	m.SeqId = id
}

type EventHeartbeat struct {
	MessageGeneric
	EvtId string `json:"EvtId"`
	Error string `json:"Error"`
}

func (m *EventHeartbeat) Name() string {
	return m.MsgType
}
func (m *EventHeartbeat) Target() string {
	return m.CarId
}
func (m *EventHeartbeat) SetSequenceId(id int) {
	m.SeqId = id
}

type EventCompleted struct {
	MessageGeneric
	EvtId string `json:"EvtId"`
	Error string `json:"Error"`
}

func (m *EventCompleted) Name() string {
	return m.MsgType
}
func (m *EventCompleted) Target() string {
	return m.CarId
}
func (m *EventCompleted) SetSequenceId(id int) {
	m.SeqId = id
}

type ResponseAck struct {
	MessageGeneric
	CorId string `json:"CorId"`
	Error string `json:"Error"`
}

func (m *ResponseAck) Name() string {
	return m.MsgType
}
func (m *ResponseAck) Target() string {
	return m.CarId
}
func (m *ResponseAck) SetSequenceId(id int) {
	m.SeqId = id
}
