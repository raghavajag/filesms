package middleware

import (
	response "filesms/pkg/api"
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
			switch e := err.(type) {
			case validator.ValidationErrors:
				apiErr := errors.HandleValidationErrors(e)
				response.JSON(w, apiErr.StatusCode, apiErr.Message, apiErr.Details)
			case errors.APIError:
				response.JSON(w, e.StatusCode, e.Message, e.Details)
			default:
				response.InternalServerError(w, "Internal server error")
			}
		}
	}
}
