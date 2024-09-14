package errors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type APIError struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("status %d: %s", e.StatusCode, e.Message)
}

func NewAPIError(statusCode int, message string, details interface{}) APIError {
	return APIError{
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
	}
}

func HandleValidationErrors(err error) APIError {
	var details []string
	for _, err := range err.(validator.ValidationErrors) {
		details = append(details, fmt.Sprintf("%s failed on the '%s' tag", err.Field(), err.Tag()))
	}
	return NewAPIError(http.StatusBadRequest, "Validation failed", details)
}

func (e APIError) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.StatusCode)
	json.NewEncoder(w).Encode(e)
}
