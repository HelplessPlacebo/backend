package router

import (
	"net/http"
	"time"

	"github.com/HelplessPlacebo/backend/gateway/config"
	"github.com/HelplessPlacebo/backend/gateway/internal/api"
	"github.com/HelplessPlacebo/backend/gateway/internal/authclient"
	"github.com/HelplessPlacebo/backend/gateway/internal/middleware"
	"github.com/HelplessPlacebo/backend/gateway/internal/proxy"
	"github.com/HelplessPlacebo/backend/pkg/shared"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func NewRouter(p proxy.Client, logger *shared.Logger, cfg *config.Config) http.Handler {
	r := chi.NewRouter()
	v := validator.New()

	baseUrl := cfg.AuthBaseURL()
	authClient := authclient.New(baseUrl+cfg.SessionEndpoint, baseUrl+cfg.RefreshSessionEndpoint, 2*time.Second)

	authMw := middleware.AuthMiddleware(authClient, 1000*time.Millisecond, 1500*time.Millisecond)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Infof("%s %s", r.Method, r.RequestURI)
			next.ServeHTTP(w, r)
		})
	})

	r.Route("/api/v1", func(r chi.Router) {
		api.RegisterRegistration(r, p, v, logger, cfg.RegistrationEndpoint)
		api.RegisterLogin(r, p, v, logger, cfg.LoginEndpoint)
		api.RegisterLogout(r, p, logger, cfg.LogoutEndpoint)
	})

	r.Group(func(r chi.Router) {
		r.Use(authMw)
	})

	return r
}
