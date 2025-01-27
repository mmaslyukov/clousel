package store

type IPortBookControllerApi interface {
	ReadPriceOptions(carId Carousel) ([]PriceTag, IError)
	Checkout(carId Carousel, priceId string, dUrl string) (ISession, IError)
	ApplyPaymentResults(sessionId string, status string) IError
}
