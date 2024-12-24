package stripe

import (
	erro "accountant/core/owner/error"
	"accountant/core/store"
	errs "accountant/core/store/error"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/webhookendpoint"
)

type CheckoutSessionWrapper struct {
	s *stripe.CheckoutSession
}

func (l *CheckoutSessionWrapper) Url() string {
	return l.s.URL
}
func (l *CheckoutSessionWrapper) Id() string {
	return l.s.ID
}

type StripeGateway struct {
}

func StripeGatewayCreate() *StripeGateway {
	return &StripeGateway{}
}

func (s *StripeGateway) ReadPriceListByProdId(skey string, prodId string, limit int) ([]store.PriceTag, errs.IError) {
	var priceArray []store.PriceTag
	var ierr errs.IError

	params := &stripe.PriceListParams{}
	params.Limit = stripe.Int64(int64(limit))
	pc := price.Client{B: stripe.GetBackend(stripe.APIBackend), Key: skey}
	result := pc.List(params)
	for result.Next() {
		if result.Price().Product.ID != prodId {
			continue
		}

		var tickets int
		if result.Price().TransformQuantity != nil {
			tickets = int(result.Price().TransformQuantity.DivideBy)
		} else {
			tickets = 1
		}

		priceArray = append(priceArray, store.PriceTag{
			PriceId: result.Price().ID,
			Amount:  int(result.Price().UnitAmount),
			Tickets: tickets,
		})
	}
	return priceArray, ierr
}

func (s *StripeGateway) GenCheckoutSessionUrl(skey string, priceId string, url store.PaymentResltUrls) (store.ISession, errs.IError) {
	var ierr errs.IError

	sc := session.Client{B: stripe.GetBackend(stripe.APIBackend), Key: skey}
	params := &stripe.CheckoutSessionParams{
		SuccessURL: &url.Success, //stripe.String(domainURL + "/success.html?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  &url.Cancel,  //stripe.String(domainURL + "/canceled.html"),
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Quantity: stripe.Int64(1),
				Price:    stripe.String(priceId),
			},
		},
	}

	ss, err := sc.New(params)
	if err != nil {
		ierr = errs.New(errs.ECStripeCheckout).Msgf("Error while creating session %v", err.Error())
	}
	return &CheckoutSessionWrapper{s: ss}, ierr
}

func (s *StripeGateway) ReadPriceDetails(skey string, priceId string) (store.PriceTag, errs.IError) {
	var ierr errs.IError
	var pt store.PriceTag
	c := price.Client{B: stripe.GetBackend(stripe.APIBackend), Key: skey}
	if price, err := c.Get(priceId, &stripe.PriceParams{}); err == nil {
		var tickets int
		if price.TransformQuantity != nil {
			tickets = int(price.TransformQuantity.DivideBy)
		} else {
			tickets = 1
		}
		pt.PriceId = priceId
		pt.Amount = int(price.UnitAmount)
		pt.Tickets = tickets
	} else {
		ierr = errs.New(errs.ECStripeReadPrice).Msgf("Error while reading price: %v", err.Error())
	}

	return pt, ierr
}

func (s *StripeGateway) WebhookRegister(url string, skey string) (string, string, erro.IError) {
	var ierr erro.IError
	var whkey, whid string
	wc := webhookendpoint.Client{B: stripe.GetBackend(stripe.APIBackend), Key: skey}
	params := &stripe.WebhookEndpointParams{
		EnabledEvents: []*string{
			// stripe.String("charge.succeeded"),
			// stripe.String("charge.failed"),
			stripe.String("checkout.session.completed"),
		},
		URL: stripe.String(url),
	}

	if result, err := wc.New(params); err == nil {
		whkey = result.Secret
		whid = result.ID
	} else {
		ierr = erro.New(erro.ECStripeRegWebhook).Msgf("Error while registering a webhook url: %v", err.Error())
	}

	return whid, whkey, ierr
}

func (s *StripeGateway) WebhookUpdateUrl(url string, skey string, whkeyId string) erro.IError {
	var ierr erro.IError
	wc := webhookendpoint.Client{B: stripe.GetBackend(stripe.APIBackend), Key: skey}
	params := &stripe.WebhookEndpointParams{
		EnabledEvents: []*string{
			// stripe.String("charge.succeeded"),
			// stripe.String("charge.failed"),
			stripe.String("checkout.session.completed"),
		},
		URL: stripe.String(url),
	}

	if _, err := wc.Update(whkeyId, params); err != nil {
		ierr = erro.New(erro.ECStripeRegWebhook).Msgf("Error while updating a webhook url: %v", err.Error())
	}

	return ierr
}

// func (s *StripeGateway) RegisterWebhook(url string, skey string, whkeyId *string) (string, string, erro.IError) {
// 	var ierr erro.IError
// 	var whkey, whid string
// 	wc := webhookendpoint.Client{B: stripe.GetBackend(stripe.APIBackend), Key: skey}
// 	params := &stripe.WebhookEndpointParams{
// 		EnabledEvents: []*string{
// 			// stripe.String("charge.succeeded"),
// 			// stripe.String("charge.failed"),
// 			stripe.String("checkout.session.completed"),
// 		},
// 		URL: stripe.String(url),
// 	}

// 	if whkeyId == nil {
// 		if result, err := wc.New(params); err == nil {
// 			whkey = result.Secret
// 			whid = result.ID
// 		} else {
// 			ierr = erro.New(erro.ECStripeRegWebhook).Msgf("Error while registering a webhook url: %v", err.Error())
// 		}
// 	} else {
// 		if result, err := wc.Update(*whkeyId, params); err == nil {
// 			whkey = result.Secret
// 			whid = *whkeyId
// 		} else {
// 			ierr = erro.New(erro.ECStripeRegWebhook).Msgf("Error while updating a webhook url: %v", err.Error())
// 		}
// 	}

// 	return whid, whkey, ierr
// }
