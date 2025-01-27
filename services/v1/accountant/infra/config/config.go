package config

import (
	"fmt"
	"os"
)

const (
	wbhook     = "109.110.19.178:4321"
	server     = "localhost:4321"
	carser     = "localhost:8081"
	sqlitePath = "accountant.db"
)

type Config struct {
}

func New() *Config {
	return &Config{}
}

func (c *Config) DatabseUrl() string {
	return sqlitePath
}

// testing data. Later will be in the database
//
//	func (c *Config) PublishKey() string {
//		return os.Getenv("PKEY")
//	}
//
//	func (c *Config) SecretKey() string {
//		return os.Getenv("SKEY")
//	}
func (c *Config) WebhookKey() string {
	return os.Getenv("WHKEY")
}

func (c *Config) ServerAddress() string {
	return server
}

func (c *Config) ServerKeyPath() string {
	return os.Getenv("SRVKEYFILE")
}

func (c *Config) ServerCertPath() string {
	return os.Getenv("SRVCERTFILE")
}

func (c *Config) ExternalServiceCarouselRegisterUrl() string {
	return fmt.Sprintf("http://%s/carousel", carser)
}

func (c *Config) ExternalServiceCarouselRefillUrl() string {
	return fmt.Sprintf("http://%s/carousel/refill", carser)
}

func (c *Config) WebhookUrl(id string) string {
	return fmt.Sprintf("http://%s/webhook/%s", wbhook, id)
	// return fmt.Sprintf("http://%s/webhook/%s", server, id)
}

// WebhookUrl() string
