package token

import (
	"errors"
	"time"

	"github.com/HelplessPlacebo/backend/auth-service/config"
	"github.com/HelplessPlacebo/backend/auth-service/internal/storage"
	"github.com/HelplessPlacebo/backend/pkg/shared"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secret     string
	accessTTL  time.Duration
	refreshTTL time.Duration
	repo       *storage.RefreshTokenRepo
	cfg        *config.Config
}

func NewTokenService(secret string, accessTTL, refreshTTL time.Duration, repo *storage.RefreshTokenRepo) *TokenService {
	return &TokenService{
		secret:     secret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		repo:       repo,
	}
}

func (s *TokenService) GenerateTokenPair(userID int) (access string, refresh string, appErr *shared.AppError) {
	access, accessErr := GenerateAccessToken(userID, s.accessTTL, s.secret)

	if accessErr != nil {
		return "", "", accessErr
	}

	rawRefresh, hash, expires, refreshErr := GenerateRefreshToken(s.refreshTTL)

	if refreshErr != nil {
		return "", "", shared.Internal("failed to generate refresh token", refreshErr)
	}

	err := s.repo.SaveHash(hash, userID, expires)

	if err != nil {
		return "", "", shared.Internal("failed to save refresh token", refreshErr)
	}

	return access, rawRefresh, nil
}

func (s *TokenService) RemoveHashedRefresh(hash string) (appErr *shared.AppError) {
	err := s.repo.DeleteHashed(hash)

	if err != nil {
		return shared.Internal("failed delete hashed refresh token", err)
	}

	return nil
}

func (s *TokenService) RemoveRefresh(token string) (appErr *shared.AppError) {
	err := s.repo.Delete(token)

	if err != nil {
		return shared.Internal("failed delete refresh token", err)
	}

	return nil
}

func (s *TokenService) GetUserIDByRefresh(refresh string) (userID int, exp time.Time, err *shared.AppError) {
	return s.repo.Find(refresh)
}

func (s *TokenService) FindRefreshTokenByUserID(userID int) (hash string, exp time.Time, err *shared.AppError) {
	return s.repo.FindByUserID(userID)
}

func (s *TokenService) VerifyAccessToken(raw string) (*AccessClaims, *shared.AppError) {
	token, err := jwt.ParseWithClaims(
		raw,
		&AccessClaims{},
		func(t *jwt.Token) (interface{}, error) {

			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, shared.BadRequest("invalid method", errors.New("unexpected signing method"))
			}

			return []byte(s.secret), nil
		},
	)
	if err != nil {
		return nil, shared.Internal("failed to parse token", err)
	}

	claims, ok := token.Claims.(*AccessClaims)

	if !ok || !token.Valid {
		return nil, shared.BadRequest("token expired", errors.New("invalid token"))
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, shared.BadRequest("token expired", errors.New("token expired"))
	}

	return claims, nil
}
