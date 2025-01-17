package config

import "time"

const (
	sqlitePath       = "carousel.db"
	mqttUrl          = "tcp://192.168.0.150:1883"
	serverDest       = ":8081"
	rootTopicCloud   = "/clousel/cloud"
	rootTopicCaursel = "/clousel/carousel"
	qos              = 1

	tmMachineMonitor = time.Second * 10
	tmMachineExpired = time.Second * 60
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

func (c *Config) ServerAddress() string {
	return serverDest
}

func (c *Config) ServerKeyPath() string {
	return "cert/server.key"
	// return "cert/dev.local+4-key.pem"
}

func (c *Config) ServerCertPath() string {
	return "cert/server.crt"
	// return "cert/myCA.pem"
	// return "cert/dev.local+4.pem"
}

func (c *Config) DatabseURL() string {
	return sqlitePath
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

func (c *Config) GetNatsUrl() string {
	return "localhost:4222"
}

func (c *Config) GetMachineMonitorTm() time.Duration {
	return tmMachineMonitor
}
func (c *Config) GetMachineExpiredTime() time.Duration {
	return tmMachineExpired
}
