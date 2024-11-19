package manager

type IPortManagerAdapterCarouselRepository interface {
	ManagerIsExistsCarousel(carId string) (bool, error)
	ManagerRemoveCarousel(carId string) error
	ManagerRemoveOwner(ownerId string) error
	ManagerReadOwnedCarousel(ownerId string) ([]Carousel, error)
	ManagerAddCarousel(c Carousel) error
}

type IPortManagerAdapterSnapshotRepository interface {
	ManagerStoreNewSnapshot(carId string) error
}
