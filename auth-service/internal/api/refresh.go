package api

import (
	"net/http"
	"time"

	"github.com/HelplessPlacebo/backend/auth-service/internal/service/token"
	"github.com/HelplessPlacebo/backend/pkg/shared"
	"github.com/go-chi/chi/v5"
)

func RegisterRefresh(r chi.Router, tokenSvc *token.TokenService, accessTokenTTL time.Duration, refreshTokenTTL time.Duration, endpoint string) {
	r.Post(endpoint, func(w http.ResponseWriter, req *http.Request) {

		refCookie, err := req.Cookie("refresh_token")
		refToken := refCookie.Value

		if err != nil {
			shared.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
			return
		}

		uID, expAt, findErr := tokenSvc.GetUserIDByRefresh(refToken)

		if findErr != nil || expAt.Before(time.Now()) {
			shared.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "session expired"})
			return
		}

		removeRefErr := tokenSvc.RemoveRefresh(refToken)

		if removeRefErr != nil {
			shared.WriteJSON(w, removeRefErr.Code, map[string]string{"error": removeRefErr.Message})

		}

		access, refreshRaw, genErr := tokenSvc.GenerateTokenPair(uID)

		if genErr != nil {
			shared.WriteJSON(w, genErr.Code, map[string]string{"error": genErr.Message})
			return
		}

		SetAuthCookie(w, access, refreshRaw, int(accessTokenTTL.Seconds()), int(refreshTokenTTL.Seconds()))

		shared.WriteJSON(w, http.StatusCreated, map[string]string{"status": "ok"})

	})

}
