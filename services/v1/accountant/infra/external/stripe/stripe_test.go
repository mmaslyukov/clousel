package stripe_test

import (
	"testing"

	"accountant/infra/logger"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhookendpoint"
	// "github.com/stripe/stripe-go"
	// "github.com/stripe/stripe-go/webhookendpoint"
)

func TestWebhookRegistration(t *testing.T) {
	stripe.Key = "sk_test_51PajXoRubpSlGSkxRr6WpEzbhLnZH7fV8ly3yhPNWKsHG7ArdsKQAjVXj6iftvOIiBs5Prp5732t4YbBTQ54v9zI00tAccea11"

	params := &stripe.WebhookEndpointParams{
		EnabledEvents: []*string{
			stripe.String("charge.succeeded"),
			stripe.String("charge.failed"),
			stripe.String("checkout.session.completed"),
		},
		URL: stripe.String("https://example.com/my/webhook/endpoint"),
	}
	log := logger.New()
	result, _ := webhookendpoint.New(params)
	log.Debug().Msgf("Result %v", result)
}
