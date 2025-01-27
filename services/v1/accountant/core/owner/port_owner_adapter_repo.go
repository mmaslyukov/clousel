package owner

type IPortOwnerAdapterProfileRepo interface {
	OwnerRegister(email string, password string, role UserRole) IError
	OwnerReadEntry(email string) (OwnerEntry, IError)
	OwnerReadEntryByOwner(ownerId string) (OwnerEntry, IError)
	OwnerAssignStripeKeys(ownerId Owner, pk *string, sk *string) IError
	OwnerAssignWebhook(ownerId Owner, whid string, whkey string) IError
}

type IPortOwnerAdapterProductRepo interface {
	OwnerDeleteCarousel(carId Carousel) IError
	OwnerAddCarousel(ownerId Owner, carId Carousel) IError
	OwnerReadProdEntry(carousel Carousel) (ProductEntry, IError)
	OwnerReadProdEntries(ownerId Owner) ([]ProductEntry, IError)
	OwnerAssignStripeProductId(carId Carousel, prodId Product) IError
}
