package middleware

import (
	"filesms/pkg/errors"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ErrorHandlerFunc func(http.ResponseWriter, *http.Request) error

func ErrorHandler(next ErrorHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)
		if err != nil {
			log.Printf("Error: %v", err)
			var apiErr errors.APIError
			switch e := err.(type) {
			case validator.ValidationErrors:
				apiErr = errors.HandleValidationErrors(e)
			case errors.APIError:
				apiErr = e
			default:
				apiErr = errors.NewAPIError(http.StatusInternalServerError, "Internal server error", nil)
			}
			apiErr.Write(w)
		}
	}
}
