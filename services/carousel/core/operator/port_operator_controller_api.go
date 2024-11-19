package operator

type IPortOperatorControllerApi interface {
	Refill(c Carousel, rounds int) error
	Play(c Carousel) error
	Read(c Carousel) (*SnapshotData, error)
	ReadByStatus(status string) ([]SnapshotData, error)
	ReadPending() ([]CompositeData, error)
}
