package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/HelplessPlacebo/backend/auth-service/internal/service/auth"
	"github.com/HelplessPlacebo/backend/pkg/shared"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type RegistrationRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
}

func RegisterRegistration(r chi.Router, svc *auth.AuthService, v *validator.Validate, logger *shared.Logger, endpoint string) {
	r.Post(endpoint, func(w http.ResponseWriter, req *http.Request) {
		var body RegistrationRequest
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

		if appErr := svc.Register(body.Email, body.Password, body.Name); appErr != nil {
			logger.Errorf("failed to register user: %s; underlying: %v", appErr.Message, appErr.Err)
			shared.WriteJSON(w, appErr.Code, map[string]string{"error": appErr.Message})
			return
		}

		shared.WriteJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
	})
}
