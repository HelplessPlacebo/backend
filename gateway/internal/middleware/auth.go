package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/HelplessPlacebo/backend/gateway/internal/authclient"
)

type ctxKeyUserID struct{}

func AuthMiddleware(authClient authclient.AuthClient, sessionTimeout, refreshTimeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx1, cancel1 := context.WithTimeout(r.Context(), sessionTimeout)
			ok, sess, status, err := authClient.RequestSession(ctx1, r)
			cancel1()

			if err != nil {
				http.Error(w, "auth service error", http.StatusBadGateway)
				return
			}

			if ok {
				ctx := context.WithValue(r.Context(), ctxKeyUserID{}, sess.UserID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if status != http.StatusUnauthorized {
				http.Error(w, http.StatusText(status), status)
				return
			}

			ctx2, cancel2 := context.WithTimeout(r.Context(), refreshTimeout)
			okRef, sessRef, setCookies, refreshStatus, err := authClient.RequestRefresh(ctx2, r)
			cancel2()

			if err != nil {
				http.Error(w, "auth service error", http.StatusBadGateway)
				return
			}

			if okRef {
				for _, c := range setCookies {
					http.SetCookie(w, c)
				}
				ctx := context.WithValue(r.Context(), ctxKeyUserID{}, sessRef.UserID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if refreshStatus != 0 {
				http.Error(w, http.StatusText(refreshStatus), refreshStatus)
				return
			}

			http.Error(w, "unauthorized", http.StatusUnauthorized)
		})
	}
}
