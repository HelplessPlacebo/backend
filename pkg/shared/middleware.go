package shared

import (
	"net/http"
)

// ErrorHandler middleware: if handler returned an *AppError via context or panic, convert to JSON.
// For simplicity we'll use pattern: handlers return *AppError or nil; we wrap in adapter below.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) *AppError

// Adapt converts HandlerFunc to http.HandlerFunc with AppError handling.
func Adapt(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// recover panic, wrap and respond
				appErr := Internal("internal server error", nil)
				WriteJSON(w, appErr.Code, map[string]string{"error": appErr.Message})
			}
		}()

		if err := h(w, r); err != nil {
			WriteJSON(w, err.Code, map[string]string{"error": err.Message})
		}
	}
}
