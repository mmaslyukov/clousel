package cfg

import (
	"fmt"
	"os"
	"time"
)

const (
	wbhook     = "109.110.19.178:4321"
	server     = "localhost:4321"
	sqlitePath = "clousel.db"

	gameStartTm   = time.Second * 30
	gameMonitorTm = time.Second * 5
)

type Config struct {
}

func New() *Config {
	return &Config{}
}

func (c *Config) DatabseUrl() string {
	return sqlitePath
}

func (c *Config) WebhookKey() string {
	return os.Getenv("WHKEY")
}

func (c *Config) ServerAddress() string {
	return server
}

func (c *Config) WebhookUrl(id string) string {
	return fmt.Sprintf("http://%s/webhook/%s", wbhook, id)
}
func (c *Config) GetNatsUrl() string {
	return "localhost:4222"
}

func (c *Config) GameStartTimeout() time.Duration {
	return gameStartTm
}
func (c *Config) GameStartMonitor() time.Duration {
	return gameMonitorTm
}
