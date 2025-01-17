package business

import "clousel/lib/fault"

type IBusinessStripeAdapter interface {
	WebhookRegister(url string, skey string) (string, string, fault.IError)
	WebhookUpdateUrl(url string, skey string, whkeyId string) fault.IError
}
