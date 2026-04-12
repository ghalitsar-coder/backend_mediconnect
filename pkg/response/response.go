package response

import (
	"encoding/json"
	"net/http"
)

// APIResponse is the standard JSON envelope for all API responses.
type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// JSON writes a JSON-encoded payload to the ResponseWriter with the given status code.
func JSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// Success writes a 2xx response with a standard success envelope.
func Success(w http.ResponseWriter, statusCode int, message string, data any) {
	JSON(w, statusCode, APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// Error writes an error response with a standard error envelope.
func Error(w http.ResponseWriter, statusCode int, message string) {
	JSON(w, statusCode, APIResponse{
		Status:  "error",
		Message: message,
	})
}
