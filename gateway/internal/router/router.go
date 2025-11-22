package router

import (
	"net/http"

	"github.com/HelplessPlacebo/backend/gateway/internal/api"
	"github.com/HelplessPlacebo/backend/gateway/internal/proxy"
	"github.com/HelplessPlacebo/backend/pkg/shared"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func NewRouter(p proxy.Client, logger *shared.Logger) http.Handler {
	r := chi.NewRouter()
	v := validator.New()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Infof("%s %s", r.Method, r.RequestURI)
			next.ServeHTTP(w, r)
		})
	})

	r.Route("/api/v1", func(r chi.Router) {
		api.RegisterRegistration(r, p, v, logger)
	})

	return r
}
