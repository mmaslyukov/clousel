package store

type IPortStoreAdapterCarouselService interface {
	Refill(carId Carousel, tickets int) IError
}
