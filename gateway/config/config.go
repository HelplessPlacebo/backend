package config

import (
	"fmt"

	"github.com/HelplessPlacebo/backend/pkg/shared"
)

type Config struct {
	AuthHost string
	AuthPort string
	AuthBase string
	Port     string
}

func Load() *Config {
	return &Config{
		AuthHost: shared.String("AUTH_SERVICE_HOST", "localhost"),
		AuthPort: shared.String("AUTH_SERVICE_PORT", "8081"),
		AuthBase: shared.String("AUTH_SERVICE_BASE_PATH", "/api/v1"),
		Port:     shared.String("GATEWAY_PORT", "8080"),
	}
}

func (c *Config) AuthBaseURL() string {
	return fmt.Sprintf("http://%s:%s%s", c.AuthHost, c.AuthPort, c.AuthBase)
}
