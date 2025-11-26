package router

import (
	"net/http"

	"github.com/HelplessPlacebo/backend/auth-service/config"
	"github.com/HelplessPlacebo/backend/auth-service/internal/api"
	"github.com/HelplessPlacebo/backend/auth-service/internal/service/auth"
	"github.com/HelplessPlacebo/backend/auth-service/internal/service/token"
	"github.com/HelplessPlacebo/backend/pkg/shared"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func NewRouter(authSvc *auth.AuthService, tokensSvc *token.TokenService, cfg *config.Config, logger *shared.Logger) http.Handler {
	r := chi.NewRouter()
	v := validator.New()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Infof("%s %s", r.Method, r.RequestURI)
			next.ServeHTTP(w, r)
		})
	})

	r.Route(cfg.APIBase, func(r chi.Router) {
		api.RegisterRegistration(r, authSvc, v, logger, cfg.RegistrationEndpoint)
		api.RegisterLogin(r, authSvc, tokensSvc, v, logger, cfg.AccessTokenTTL, cfg.RefreshTokenTTL, cfg.LoginEndpoint)
		api.RegisterLogout(r, tokensSvc, cfg.LogoutEndpoint)
	})
	return r
}
