package config

import (
	"time"

	"github.com/HelplessPlacebo/backend/pkg/shared"
)

type Config struct {
	DBURL                string
	Port                 string
	APIBase              string
	JWTSecret            string
	AccessTokenTTL       time.Duration
	RefreshTokenTTL      time.Duration
	LoginEndpoint        string
	RegistrationEndpoint string
	LogoutEndpoint       string
}

func Load() *Config {

	secret := shared.String("JWT_SECRET", "")
	if secret == "" {
		panic("JWT_SECRET is required and cannot be empty")
	}

	return &Config{
		DBURL:                shared.String("DATABASE_URL", "postgres://user:pass@localhost:5432/auth?sslmode=disable"),
		Port:                 shared.String("AUTH_SERVICE_PORT", "8081"),
		APIBase:              shared.String("AUTH_SERVICE_BASE_PATH", "/api/v1"),
		JWTSecret:            secret,
		AccessTokenTTL:       15 * time.Minute,
		RefreshTokenTTL:      7 * 24 * time.Hour,
		LoginEndpoint:        shared.String("AUTH_SERVICE_LOGIN_PATH", "/login"),
		RegistrationEndpoint: shared.String("AUTH_SERVICE_REG_PATH", "/registration"),
		LogoutEndpoint:       shared.String("AUTH_SERVICE_LOGOUT_PATH", "/logout"),
	}
}
