package port_persistency

import (
	. "carousel_service/internal/utils"
)

// type Pair struct {
// 	name  string
// 	value any
// }

type RecordInterface[T any] interface {
	Create(record T) error
	Update(record T) error
	// UpdateBy(carouselId string, sets []Pair) error
	Delete(carouselId string) error
	Read(carouselId string) (Optional[T], error)
	ReadOneBy(where func() string) (Optional[T], error)
	ReadManyBy(where func() string) (Optional[[]T], error)
	ReadAll() (Optional[[]T], error)
}

type CarouselRecord struct {
	CarouselId string
	OwnerId    string
	RoundTime  Optional[int]
}

func NewDefaultCarouselRecord() CarouselRecord {
	return CarouselRecord{
		CarouselId: "",
		RoundTime:  NewOptionalValue[int](0),
	}
}

const (
	StatusNameOnline  = "online"
	StatusNameOffline = "offline"
	StatusNameoNew    = "new"
)

type StatusRecord struct {
	CarouselId  string
	Status      Optional[string]
	RoundsReady Optional[int]
	Time        string
}

func NewDefaultStatusRecord() StatusRecord {
	return StatusRecord{
		CarouselId:  "",
		Status:      NewOptionalValue[string](""),
		RoundsReady: NewOptionalValue[int](0),
		Time:        "",
	}
}

type EventRecord struct {
	EventId string
	Time    string
	Type    string
	Data    string
}

type PersistencyInterface interface {
	Close() error
	Status() RecordInterface[StatusRecord]
	Carousel() RecordInterface[CarouselRecord]
}
type PersistencyGatewayInterface interface {
	Open() (PersistencyInterface, error)
}
