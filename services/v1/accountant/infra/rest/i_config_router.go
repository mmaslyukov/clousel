package rest

type IConfigRouter interface {
	WebhookKey() string
	WebhookUrl(id string) string
	ServerAddress() string
	ServerKeyPath() string
	ServerCertPath() string
}
