package token

import (
	"time"

	"github.com/HelplessPlacebo/backend/pkg/shared"
)

func GenerateRefreshToken(ttl time.Duration) (raw string, hash string, expires time.Time, err error) {
	raw, err = shared.GenerateRandomString(64)
	if err != nil {
		return "", "", time.Time{}, err
	}

	hash = shared.HashSHA256(raw)
	expires = time.Now().Add(ttl)

	return raw, hash, expires, nil
}
