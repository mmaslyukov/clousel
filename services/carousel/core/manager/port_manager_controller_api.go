package manager

type IPortManagerControllerApi interface {
	Register(c Carousel) error
	Unregister(c Carousel) error
	Read(ownerId string) ([]Carousel, error)
}
