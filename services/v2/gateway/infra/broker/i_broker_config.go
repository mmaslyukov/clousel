package broker

type IBrokerConfig interface {
	BrokerURL() string
	BrokerUsername() string
	BrokerPassword() string
}
