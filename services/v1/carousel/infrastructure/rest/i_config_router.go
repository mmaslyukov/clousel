package rest

type IConfigRouter interface {
	ServerAddress() string
	ServerKeyPath() string
	ServerCertPath() string
}
