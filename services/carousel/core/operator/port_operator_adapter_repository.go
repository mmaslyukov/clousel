package operator

import (
	"time"

	"github.com/google/uuid"
)

type IPortOperatorAdapterRepository interface {
	IsExists(c Carousel) (bool, error)
	Refill(rd *RoundsData) error
	Play(rd *PlayData) error

	Mark(s *StatusData) error
	Confirm(s *StatusData) error
	UpdateTime(rd *Carousel) error
	ClearPendingFlag(rd *EventData) error

	Read(c *Carousel) ([]CompositeData, error)
	// ReadOwned(ownerId string) ([]CompositeData, error)
	ReadAsSnapshot(c *Carousel) (SnapshotData, error)
	SaveSnapshot(c *SnapshotData) error
	// ReadPending(c *Carousel) ([]CompositeData, error)
	ReadPendingTimeout(dur time.Duration) ([]CompositeData, error)
	ReadPending() ([]CompositeData, error)
	ReadWStatus(status string) ([]SnapshotData, error)
	ReadExpired(dur time.Duration) ([]CompositeData, error)

	Remove(c *Carousel) error
	RemoveByEvent(evtId uuid.UUID) error
}
