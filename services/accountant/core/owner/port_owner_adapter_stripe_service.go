package owner

type IPortOwnerAdapterStripeService interface {
	RegisterWebhook(url string, skey string, whkeyId *string) (string, string, IError)
}
