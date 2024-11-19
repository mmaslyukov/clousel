package operator

import (
	"time"

	"github.com/google/uuid"
)

type IPortOperatorAdapterEventRepository interface {
	Refill(rd *RoundsData) error
	Play(rd *PlayData) error

	Mark(s *StatusData) error
	Confirm(s *StatusData) error
	UpdateTime(rd *Carousel) error
	ClearPendingFlag(rd *EventData) error

	Read(carId string) ([]CompositeData, error)
	ReadAsSnapshot(carId string) (*SnapshotData, error)
	ReadPendingTimeout(dur time.Duration) ([]CompositeData, error)
	ReadPending() ([]CompositeData, error)
	ReadByStatus(status string) ([]SnapshotData, error)
	ReadExpired(dur time.Duration) ([]CompositeData, error)

	Remove(carId string) error
	RemoveByEvent(evtId uuid.UUID) error
}

type IPortOperatorAdapterCarouselRepository interface {
	OperatorIsExistsCarousel(carId string) (bool, error)
	OperarotReadAllCarouselIds() ([]string, error)
}

type IPortOperatorAdapterSnapshotRepository interface {
	OperatorLoadSnapshot(carId string) (*SnapshotData, error)
	OperatorStoreSnapshot(sd *SnapshotData) error
	OperatorDeleteSnapshot(carId string) error
}
