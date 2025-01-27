package operator

import (
	"time"

	"github.com/google/uuid"
)

type IPortOperatorAdapterEventRepository interface {
	OperatorRefill(rd *TicketsData) error
	OperatorPlay(rd *PlayData) error

	OperatorMark(s *StatusData) error
	OperatorConfirm(s *StatusData) error
	OperatorUpdateTime(rd *Carousel) error
	OperatorClearPendingFlag(rd *EventData) error

	OperatorRead(carId string) ([]CompositeData, error)
	OperatorReadAsSnapshot(carId string) (*SnapshotData, error)
	OperatorReadPendingTimeout(dur time.Duration) ([]CompositeData, error)
	OperatorReadPending() ([]CompositeData, error)
	OperatorReadByStatus(status string) ([]SnapshotData, error)
	OperatorReadExpired(dur time.Duration) ([]CompositeData, error)

	OperatorRemove(carId string) error
	OperatorRemoveByEvent(evtId uuid.UUID) error
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
