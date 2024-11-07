package config

const (
	sqlitePath       = "carousel.db"
	mqttUrl          = "tcp://192.168.0.150:1883"
	serverDest       = "localhost:8081"
	rootTopicCloud   = "/clousel/cloud"
	rootTopicCaursel = "/clousel/carousel"
	qos              = 1
)

// func GetSqlitePath() string {
// 	return sqlitePath
// }

// func GetMQTTUrl() string {
// 	return mqttUrl
// }

// func GetServerDest() string {
// 	return serverDest
// }
// func GetTopicCloud() string {
// 	return rootTopicCloud
// }
// func GetTopicCarousel() string {
// 	return rootTopicCaursel
// }

type Config struct {
}

func New() *Config {
	return &Config{}
}
func (c *Config) DatabseURL() string {
	return sqlitePath
}
func (c *Config) Server() string {
	return serverDest
}
func (c *Config) BrokerURL() string {
	return mqttUrl
}
func (c *Config) RootTopicPub() string {
	return rootTopicCloud
}
func (c *Config) RootTopicSub() string {
	return rootTopicCaursel
}
func (c *Config) DefaultQOS() byte {
	return qos
}
