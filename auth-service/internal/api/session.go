package api

import (
	"net/http"

	"github.com/HelplessPlacebo/backend/auth-service/internal/service/token"
	"github.com/HelplessPlacebo/backend/pkg/shared"
	"github.com/go-chi/chi/v5"
)

func RegisterSession(r chi.Router, tokenSvc *token.TokenService, endpoint string) {
	r.Get(endpoint, func(w http.ResponseWriter, req *http.Request) {

		aToken, err := req.Cookie("access_token")

		if err == nil {
			shared.WriteJSON(w, http.StatusUnauthorized, &shared.AppError{Message: "Invalid token"})
			return
		}

		claims, vErr := tokenSvc.VerifyAccessToken(aToken.Value)

		if vErr != nil {
			shared.WriteJSON(w, vErr.Code, vErr.Message)
			return
		}

		shared.WriteJSON(w, http.StatusCreated, map[string]int{"userID": claims.UserID})

	})

}
