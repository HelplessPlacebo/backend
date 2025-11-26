package shared

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// String returns environment variable or fallback.
func String(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Int returns integer environment variable or fallback.
func Int(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

// Bool returns boolean environment variable or fallback.
func Bool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}

func Duration(key string, def int) time.Duration {
	return time.Duration(Int(key, def)) * time.Second
}

func InitEnv(logger *Logger) {

	paths := []string{
		".env",
		".env.local",
		".env.development",
		".env.example",
		"../.env",
		"../.env.example",
		"../../.env",
		"../../.env.example",
	}

	loaded := false
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			if err := godotenv.Load(p); err == nil {
				logger.Infof("loaded env from %s", p)
				loaded = true
				break
			}
		}
	}

	if !loaded {
		logger.Infof("no .env file loaded (looked in likely locations); relying on OS env")
	}
}
