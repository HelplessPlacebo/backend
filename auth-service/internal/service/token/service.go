package token

import (
	"time"

	"github.com/HelplessPlacebo/backend/auth-service/internal/storage"
	"github.com/HelplessPlacebo/backend/pkg/shared"
)

type TokenService struct {
	secret     string
	accessTTL  time.Duration
	refreshTTL time.Duration
	repo       *storage.RefreshTokenRepo
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

func (s *TokenService) CheckExistingLogin(refresh string) (userID int, exp time.Time, err error) {
	return s.repo.Find(refresh)
}

func (s *TokenService) FindActiveByUserID(userID int) (hash string, exp time.Time, err error) {
	return s.repo.FindByUserID(userID)
}
