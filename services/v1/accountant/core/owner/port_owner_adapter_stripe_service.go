package owner

type IPortOwnerAdapterStripeService interface {
	WebhookRegister(url string, skey string) (string, string, IError)
	WebhookUpdateUrl(url string, skey string, whkeyId string) IError
}
