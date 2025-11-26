package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/HelplessPlacebo/backend/gateway/internal/proxy"
	"github.com/HelplessPlacebo/backend/pkg/shared"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func RegisterLogin(r chi.Router, p proxy.Client, v *validator.Validate, logger *shared.Logger, endpoint string) {
	r.Post(endpoint, func(w http.ResponseWriter, req *http.Request) {
		var body LoginRequest

		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			logger.Infof("Failed to decode %v request: %v", endpoint, err)
			shared.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}

		if err := v.Struct(body); err != nil {
			logger.Infof("Failed to validate %v request: %v", endpoint, err.Error())
			shared.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		resp, err := p.ForwardJSON(context.Background(), req, endpoint, body)
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
