package config

const (
	sqlitePath       = "sqlite.db"
	mqttUrl          = "tcp://192.168.0.150:1883"
	serverDest       = "localhost:8080"
	rootTopicCloud   = "/clousel/cloud"
	rootTopicCaursel = "/clousel/carousel"
)

// type ConfigBrokerInterface interface {
// 	GetTopicRoot() string
// }

// type ConfigBrokerCloud struct {
// }

// type ConfigBrokerCarousel struct {
// }

// func (c *ConfigBrokerCarousel) GetTopicRoot() string {
// 	return rootTopicCaursel
// }
// func (c *ConfigBrokerCloud) GetTopicRoot() string {
// 	return rootTopicCloud
// }
// func NewConfigBroker(role string) ConfigBrokerInterface {
// 	switch role {
// 	case "cloud":
// 		return &ConfigBrokerCloud{}
// 	case "carousel":
// 		return &ConfigBrokerCarousel{}
// 	}
// 	return &ConfigBrokerCloud{}
// }
func GetSqlitePath() string {
	return sqlitePath
}

func GetMQTTUrl() string {
	return mqttUrl
}

func GetServerDest() string {
	return serverDest
}
func GetTopicCloud() string {
	return rootTopicCloud
}
func GetTopicCarousel() string {
	return rootTopicCaursel
}
