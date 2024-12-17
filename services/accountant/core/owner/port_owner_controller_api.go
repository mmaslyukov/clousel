package owner

type IPortOwnerControllerOwnerApi interface {
	Register(email string, password string) IError
	Login(email string, password string) (Token, IError)
	AddCarousel(token Token, carId Carousel, prodId *string) IError
	AssignProdId(token Token, carId Carousel, prodId Product) IError
	AssignSkeys(token Token, skey string) IError
	AssignPkeys(token Token, pkey string) IError
	RefreshWebhook(token Token) IError
	ReadWhkey(ownerId string) (string, IError)
}

type IPortOwnerControllerAdminApi interface {
	ChangerCarouselOwner(carId Carousel, owner Owner) IError
}
