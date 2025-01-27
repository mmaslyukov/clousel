package stripe

import (
	"clousel/core/client"
	"clousel/lib/fault"
	"slices"

	"github.com/rs/zerolog"
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
	log *zerolog.Logger
}

func StripeGatewayCreate(log *zerolog.Logger) *StripeGateway {
	return &StripeGateway{
		log: log,
	}
}

func (s *StripeGateway) ReadPriceListByProdId(skey string, prodId string, limit int) ([]*client.PriceTag, fault.IError) {
	var priceArray []*client.PriceTag
	var ierr fault.IError

	params := &stripe.PriceListParams{}
	params.Limit = stripe.Int64(int64(limit))
	pc := price.Client{B: stripe.GetBackend(stripe.APIBackend), Key: skey}
	result := pc.List(params)
	for result.Err() == nil && result.Next() {
		if result.Price().Product.ID != prodId {
			continue
		}

		var tickets int
		if result.Price().TransformQuantity != nil {
			tickets = int(result.Price().TransformQuantity.DivideBy)
		} else {
			tickets = 1
		}
		priceArray = append(priceArray, &client.PriceTag{
			PriceId: result.Price().ID,
			Amount:  int(result.Price().UnitAmount),
			Tickets: tickets,
		})

		slices.SortFunc(priceArray, func(a, b *client.PriceTag) int {
			return a.Amount - b.Amount
		})

	}
	if result.Err() != nil {
		ierr = fault.New(EStripeReadPrice).Msgf("%s", result.Err())
	}

	return priceArray, ierr
}

func (s *StripeGateway) GenCheckoutSessionUrl(email string, skey string, priceId string, url client.PaymentResultUrls) (client.ISession, fault.IError) {
	var ierr fault.IError
	sc := session.Client{B: stripe.GetBackend(stripe.APIBackend), Key: skey}
	params := &stripe.CheckoutSessionParams{
		CustomerEmail: &email,
		SuccessURL:    &url.Success, //stripe.String(domainURL + "/success.html?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:     &url.Cancel,  //stripe.String(domainURL + "/canceled.html"),
		Mode:          stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Quantity: stripe.Int64(1),
				Price:    stripe.String(priceId),
			},
		},
	}

	ss, err := sc.New(params)
	if err != nil {
		ierr = fault.New(EStripeCheckoutSession).Msgf("Error while creating session %v", err.Error())
	}
	return &CheckoutSessionWrapper{s: ss}, ierr
}

func (s *StripeGateway) ReadPriceDetails(skey string, priceId string) (client.PriceTag, fault.IError) {
	var ierr fault.IError
	var pt client.PriceTag
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
		ierr = fault.New(EStripeReadPrice).Msgf("Error while reading price: %v", err.Error())
	}

	return pt, ierr
}

func (s *StripeGateway) WebhookRegister(url string, skey string) (string, string, fault.IError) {
	var ierr fault.IError
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
		ierr = fault.New(EStripeRegWebhook).Msgf("Error while registering a webhook url: %v", err.Error())
	}

	return whid, whkey, ierr
}

func (s *StripeGateway) WebhookUpdateUrl(url string, skey string, whkeyId string) fault.IError {
	var ierr fault.IError
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
		ierr = fault.New(EStripeRegWebhook).Msgf("Error while updating a webhook url: %v", err.Error())
	}

	return ierr
}
