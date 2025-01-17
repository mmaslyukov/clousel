package router

type IConfigRouter interface {
	WebhookKey() string
	WebhookUrl(id string) string
	ServerAddress() string
}
