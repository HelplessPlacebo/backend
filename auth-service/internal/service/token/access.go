package token

import (
	"time"

	"github.com/HelplessPlacebo/backend/pkg/shared"
	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaims struct {
	UserID int `json:"sub"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID int, ttl time.Duration, secret string) (string, *shared.AppError) {
	claims := AccessTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", shared.Internal("failed to sign access token", err)
	}

	return signed, nil
}
