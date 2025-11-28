package api

import (
	"net/http"

	"github.com/HelplessPlacebo/backend/auth-service/internal/service/token"
	"github.com/HelplessPlacebo/backend/pkg/shared"
	"github.com/go-chi/chi/v5"
)

func RegisterLogout(r chi.Router, tokenSvc *token.TokenService, endpoint string) {

	r.Post(endpoint, func(w http.ResponseWriter, req *http.Request) {

		cookie, err := req.Cookie("refresh_token")
		if err != nil {
			shared.WriteJSON(w, http.StatusOK, map[string]string{"status": "already_logged_out"})
			return
		}

		_ = tokenSvc.RemoveRefresh(cookie.Value)

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})

		shared.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

}
