package client

import "clousel/lib/fault"

type IClientStripeAdapter interface {
	ReadPriceListByProdId(skey string, prodId string, limit int) ([]PriceTag, fault.IError)
	ReadPriceDetails(skey string, priceId string) (PriceTag, fault.IError)
	GenCheckoutSessionUrl(email string, skey string, priceId string, url PaymentResultUrls) (ISession, fault.IError)
}

type ISession interface {
	Url() string
	Id() string
}
