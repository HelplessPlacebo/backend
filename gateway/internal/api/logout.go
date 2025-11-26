package api

import (
	"io"
	"net/http"

	"github.com/HelplessPlacebo/backend/gateway/internal/proxy"
	"github.com/HelplessPlacebo/backend/pkg/shared"
	"github.com/go-chi/chi/v5"
)

func RegisterLogout(r chi.Router, p proxy.Client, logger *shared.Logger, endpoint string) {
	r.Post(endpoint, func(w http.ResponseWriter, req *http.Request) {
		resp, err := p.ForwardJSON(req.Context(), req, endpoint, nil)
		if err != nil {
			logger.Errorf("proxy error: %v", err)
			shared.WriteJSON(w, http.StatusBadGateway, map[string]string{"error": "upstream error"})
			return
		}
		defer resp.Body.Close()

		for k, vals := range resp.Header {
			if k == "Set-Cookie" {
				continue
			}
			for _, v := range vals {
				w.Header().Add(k, v)
			}
		}

		shared.ProxyCookiesFromResponse(resp, w)

		w.WriteHeader(resp.StatusCode)

		_, _ = io.Copy(w, resp.Body)
	})
}
