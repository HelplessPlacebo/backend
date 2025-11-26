package token

import (
	"time"

	"github.com/HelplessPlacebo/backend/auth-service/internal/storage"
)

func GenerateRefreshToken(ttl time.Duration) (raw string, hash string, expires time.Time, err error) {
	raw, err = storage.GenerateRandomString(64)
	if err != nil {
		return "", "", time.Time{}, err
	}

	hash = storage.HashSHA256(raw)
	expires = time.Now().Add(ttl)

	return raw, hash, expires, nil
}
