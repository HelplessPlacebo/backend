package config

import "github.com/HelplessPlacebo/backend/pkg/shared"

type Config struct {
	DBURL   string
	Port    string
	APIBase string
}

func Load() *Config {
	return &Config{
		DBURL:   shared.String("DATABASE_URL", "postgres://user:pass@localhost:5432/auth?sslmode=disable"),
		Port:    shared.String("AUTH_SERVICE_PORT", "8081"),
		APIBase: shared.String("AUTH_SERVICE_BASE_PATH", "/api/v1"),
	}
}
