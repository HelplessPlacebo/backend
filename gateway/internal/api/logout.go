package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/HelplessPlacebo/backend/gateway/internal/proxy"
	"github.com/HelplessPlacebo/backend/pkg/shared"
	"github.com/go-chi/chi/v5"
)

func RegisterLogout(r chi.Router, p proxy.Client, logger *shared.Logger, endpoint string) {
	r.Post(endpoint, func(w http.ResponseWriter, req *http.Request) {

		resp, err := p.ForwardJSON(context.Background(), endpoint, nil)
		if err != nil {
			logger.Errorf("proxy error: %v", err)
			shared.WriteJSON(w, http.StatusBadGateway, map[string]string{"error": "upstream error"})
			return
		}
		defer resp.Body.Close()

		for k, vals := range resp.Header {
			for _, v := range vals {
				w.Header().Add(k, v)
			}
		}

		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)

		json.NewEncoder(w).Encode(json.RawMessage{})
	})
}
