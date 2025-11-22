package service

import (
	"github.com/HelplessPlacebo/backend/auth-service/internal/storage"
	"github.com/HelplessPlacebo/backend/pkg/shared"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users  *storage.UserRepo
	logger *shared.Logger
}

func NewAuthService(users *storage.UserRepo, logger *shared.Logger) *AuthService {
	return &AuthService{users: users, logger: logger}
}

func (a *AuthService) Register(email, password, name string) *shared.AppError {
	_, err := a.users.GetByEmail(email)
	if err == nil {
		return shared.Conflict("user already exists", nil)
	}
	if ae, ok := shared.IsAppError(err); ok && ae.Code != 404 {
		a.logger.Errorf("failed to check user: %s; underlying: %v", email, err)
		return shared.Internal("db error", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return shared.Internal("failed to hash password", err)
	}

	if err := a.users.CreateUser(email, string(hash), name); err != nil {
		if ae, ok := shared.IsAppError(err); ok {
			return ae
		}
		return shared.Internal("failed to create user", err)
	}

	a.logger.Infof("user registered: %s", email)
	return nil
}
