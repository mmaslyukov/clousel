package store

import (
	errs "accountant/core/store/error"
	"fmt"
	"net/url"

	"github.com/rs/zerolog"
)

type StoreDomain struct {
	owner    IPortStoreAdapterOwner
	carousel IPortStoreAdapterCarouselService
	stripe   IPortStoreAdapterStripeService
	bookRepo IPortStoreAdapterBookRepo
	log      *zerolog.Logger
}

func StoreDomainCreate(
	owner IPortStoreAdapterOwner,
	carousel IPortStoreAdapterCarouselService,
	bookRepo IPortStoreAdapterBookRepo,
	stripe IPortStoreAdapterStripeService,
	log *zerolog.Logger) *StoreDomain {
	return &StoreDomain{
		owner:    owner,
		carousel: carousel,
		bookRepo: bookRepo,
		stripe:   stripe,
		log:      log,
	}
}

func (sd *StoreDomain) ReadPriceOptions(carId Carousel) ([]PriceTag, IError) {
	var ierr IError
	var err error
	var prices []PriceTag
	var skey, prodId string
	for ok := true; ok; ok = false {
		if _, skey, prodId, err = sd.owner.ReadKeys(carId); err != nil {
			ierr = errs.New(errs.ECReadKeys).Msgf("Fail to read keys for carousel:%s, error:%s", carId.String(), err.Error())
			break
		}
		if len(skey) == 0 {
			ierr = errs.New(errs.ECReadKeys).Msgf("Fail to read, SecretKey is empty")
			break
		}
		prices, ierr = sd.stripe.ReadPriceListByProdId(skey, prodId, 3)
	}
	return prices, ierr
}

func (sd *StoreDomain) Checkout(carId Carousel, priceId string, dUrl string) (ISession, IError) {
	var ierr IError
	var err error
	var cs ISession
	var skey string
	var pt PriceTag
	for ok := true; ok; ok = false {
		if _, skey, _, err = sd.owner.ReadKeys(carId); err != nil {
			ierr = errs.New(errs.ECReadKeys).Msgf("Fail to read keys for carousel:%s, error:%s", carId.String(), err.Error())
			break
		}
		if pt, ierr = sd.stripe.ReadPriceDetails(skey, priceId); ierr != nil {
			ierr = errs.New(errs.ECStripeCheckoutSession).Err(ierr).Msgf("Fail to read price details")
			break
		}
		query_success := url.Values{}
		query_success.Add("type", "popup_success")
		query_success.Add("msg", "Payment has been confirmed")

		query_error := url.Values{}
		query_error.Add("type", "popup_error")
		query_error.Add("msg", "Something went wrong")
		purls := PaymentResltUrls{
			Success: fmt.Sprintf("%s?%s", dUrl, query_success.Encode()), //TODO set url
			Cancel:  fmt.Sprintf("%s?%s", dUrl, query_error.Encode()),   //TODO set url
		}
		sd.log.Debug().Str("Success url", purls.Success).Str("Error url", purls.Cancel).Send()
		if cs, ierr = sd.stripe.GenCheckoutSessionUrl(skey, priceId, purls); ierr != nil {
			cs = nil
			ierr = errs.New(errs.ECStripeCheckoutSession).Err(ierr).Msgf("Fail to make stripe checkout session")
			break
		}
		entry := BookEntryCreate(cs.Id(), carId, pt.Amount, pt.Tickets)
		if ierr = sd.bookRepo.StoreAddBookEntry(&entry); ierr != nil {
			ierr = errs.New(errs.ECStripeCheckoutSession).Err(ierr).Msgf("Fail to store book entry")
		}
	}
	return cs, ierr
}

// TODO make background thread for monitoring paid but not refilled book entries
func (sd *StoreDomain) ApplyPaymentResults(sessionId string, status string) IError {
	var ierr IError

	for ok := true; ok; ok = false {

		var be BookEntry
		if be, ierr = sd.bookRepo.StoreReadBookEntryBySessionId(sessionId); ierr != nil {
			ierr = errs.New(errs.ECBookRepoRead).Msgf("Fail to read entry").Err(ierr)
			sd.log.Err(ierr).Str("SessionId", sessionId).Str("Status", status).Send()
			break
		}
		if ierr = sd.bookRepo.StoreMarkBookEntryBySessionIdWithData(sessionId, status, nil); ierr != nil {
			ierr = errs.New(errs.ECBookRepoMark).Msgf("Fail to mark entry as paid").Err(ierr)
			sd.log.Err(ierr).Str("SessionId", sessionId).Str("Status", status).Send()
			break
		}
		if status != PaymentStatusPaid {
			ierr = errs.New(errs.ECPayment).Msgf("Payment status is invalid: '%s'", status)
			sd.log.Err(ierr).Str("SessionId", sessionId).Str("Status", status).Send()
			break
		}
		status := BookOrderStatusRefilled
		var bookErrorPtr *string = nil
		if ierr = sd.carousel.Refill(be.CarId, be.Tickets); ierr != nil {
			ierr = errs.New(errs.ECRemoteServiceCarouselRefill).Msgf("Fail to Refill carousel").Err(ierr)
			sd.log.Err(ierr).Int("Tickets", be.Tickets).Str("CarouselId", be.CarId.String()).Send()
			status = BookOrderStatusPendingRefill
			bookError := ierr.Error()
			bookErrorPtr = &bookError
		}
		if ierr = sd.bookRepo.StoreMarkBookEntryBySessionIdWithData(sessionId, status, bookErrorPtr); ierr != nil {
			ierr = errs.New(errs.ECBookRepoMark).Msgf("Fail to mark entry as refilled or canceled").Err(ierr)
			sd.log.Err(ierr).Int("Tickets", be.Tickets).Str("CarouselId", be.CarId.String()).Send()
			break
		}
	}

	return ierr
}
