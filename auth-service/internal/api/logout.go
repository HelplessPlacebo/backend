package api

import (
	"net/http"

	"github.com/HelplessPlacebo/backend/auth-service/internal/service/token"
	"github.com/HelplessPlacebo/backend/pkg/shared"
	"github.com/go-chi/chi/v5"
)

func RegisterLogout(r chi.Router, tokenSvc *token.TokenService, endpoint string) {
	r.Post(endpoint, func(w http.ResponseWriter, req *http.Request) {

		refToken, err := req.Cookie("refresh_token")

		if err != nil {
			shared.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		removeErr := tokenSvc.RemoveRefresh(refToken.Value)

		if removeErr != nil {
			shared.WriteJSON(w, removeErr.Code, map[string]string{"error": removeErr.Message})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name: "access_token", Value: "",
			HttpOnly: true, Secure: true, Path: "/", SameSite: http.SameSiteStrictMode,
			MaxAge: 0,
		})
		http.SetCookie(w, &http.Cookie{
			Name: "refresh_token", Value: "",
			HttpOnly: true, Secure: true, Path: "/", SameSite: http.SameSiteStrictMode,
			MaxAge: 0,
		})

		shared.WriteJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
	})
}
