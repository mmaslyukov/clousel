package business

type IBusinessConfigAdapter interface {
	WebhookUrl(id string) string
}
