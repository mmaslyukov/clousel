package store

type IPortStoreAdapterStripeService interface {
	ReadPriceListByProdId(skey string, prodId string, limit int) ([]PriceTag, IError)
	ReadPriceDetails(skey string, priceId string) (PriceTag, IError)
	GenCheckoutSessionUrl(skey string, priceId string, url PaymentResltUrls) (ISession, IError)
}

type ISession interface {
	Url() string
	Id() string
}
