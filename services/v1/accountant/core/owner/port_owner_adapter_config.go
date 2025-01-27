package owner

type IPortOwnerAdapterProfileConfig interface {
	WebhookUrl(id string) string
}
