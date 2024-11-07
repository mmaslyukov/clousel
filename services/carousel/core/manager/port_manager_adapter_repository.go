package manager

type IPortManagerAdapterRepository interface {
	Remove(c Carousel) error
	ReadOwned(ownerId string) ([]Carousel, error)
	AddCarousel(c Carousel) error
	AddEventWithStatusNew(c Carousel) error
	IsCarouselExistsInEvents(c Carousel) (bool, error)
	// IsExists(c Carousel) (bool, error)
}
