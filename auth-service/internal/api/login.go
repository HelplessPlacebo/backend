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

func SetAuthCookie(w http.ResponseWriter, access string, refresh string, accessMaxAge int, refreshMaxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    access,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   accessMaxAge,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   refreshMaxAge,
	})
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
			return
		}

		user, loginErr := authSvc.Login(body.Email, body.Password)

		if loginErr != nil {
			shared.WriteJSON(w, loginErr.Code, map[string]string{"error": loginErr.Message})
			return
		}

		refreshCookie, err := req.Cookie("refresh_token")

		if err == nil {

			existingUserID, expiresAt, dbErr := tokenSvc.CheckExistingLogin(refreshCookie.Value)

			if dbErr == nil && existingUserID == user.ID && expiresAt.After(time.Now()) {
				shared.WriteJSON(w, http.StatusOK, map[string]string{"status": "already_logged_in"})
				return
			}
		}

		activeHash, activeExp, findErr := tokenSvc.FindActiveByUserID(user.ID)
		if findErr == nil && activeExp.After(time.Now()) {
			_ = tokenSvc.RemoveHashedRefresh(activeHash)

			access, refreshRaw, genErr := tokenSvc.GenerateTokenPair(user.ID)
			if genErr == nil {

				SetAuthCookie(w, access, refreshRaw, int(accessTokenTTL.Seconds()), int(refreshTokenTTL.Seconds()))

				return
			}
		}

		access, refreshRaw, genErr := tokenSvc.GenerateTokenPair(user.ID)
		if genErr != nil {
			shared.WriteJSON(w, genErr.Code, map[string]string{"error": genErr.Message})
			return
		}

		SetAuthCookie(w, access, refreshRaw, int(accessTokenTTL.Seconds()), int(refreshTokenTTL.Seconds()))

		shared.WriteJSON(w, http.StatusCreated, map[string]string{"status": "ok"})

	})

}
