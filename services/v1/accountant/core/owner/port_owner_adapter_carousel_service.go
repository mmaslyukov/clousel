package owner

type IPortStoreAdapterCarouselService interface {
	Register(ownerId Owner, carId Carousel) IError
}
