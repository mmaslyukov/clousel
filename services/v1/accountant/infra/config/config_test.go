package config_test

import (
	"accountant/infra/config"
	"fmt"
	"testing"
)

func TestWebhookUrl(t *testing.T) {
	cfg := config.New()
	fmt.Printf("url - %s", cfg.WebhookUrl("123"))
}
