package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/HelplessPlacebo/backend/auth-service/internal/service/auth"
	"github.com/HelplessPlacebo/backend/auth-service/internal/service/token"
	"github.com/HelplessPlacebo/backend/pkg/shared"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func RegisterLogin(r chi.Router, authSvc *auth.AuthService, tokenSvc *token.TokenService, v *validator.Validate, logger *shared.Logger, accessTokenTTL time.Duration, refreshTokenTTL time.Duration, endpoint string) {
	r.Post(endpoint, func(w http.ResponseWriter, req *http.Request) {
		var body LoginRequest
		b, _ := io.ReadAll(req.Body)
		if err := json.Unmarshal(b, &body); err != nil {
			shared.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			logger.Errorf("invalid json %v", err)
			return
		}
		if err := v.Struct(body); err != nil {
			shared.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			logger.Errorf("invalid json %v", err.Error())
			return
		}

		user, loginErr := authSvc.Login(body.Email, body.Password)

		if loginErr != nil {
			logger.Errorf("failed to login user: %s; underlying: %v", loginErr.Message, loginErr.Err)
			shared.WriteJSON(w, loginErr.Code, map[string]string{"error": loginErr.Message})
			return
		}

		access, refresh, err := tokenSvc.GenerateTokenPair(user.ID)

		if err != nil {
			logger.Errorf("failed to generate tokens: %s; underlying: %v", err.Message, err.Err)
			shared.WriteJSON(w, err.Code, map[string]string{"error": err.Message})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name: "access_token", Value: access,
			HttpOnly: true, Secure: true, Path: "/", SameSite: http.SameSiteStrictMode,
			MaxAge: int(accessTokenTTL),
		})
		http.SetCookie(w, &http.Cookie{
			Name: "refresh_token", Value: refresh,
			HttpOnly: true, Secure: true, Path: "/", SameSite: http.SameSiteStrictMode,
			MaxAge: int(refreshTokenTTL),
		})

		shared.WriteJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
	})
}
