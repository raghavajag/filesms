package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func JSON(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func Success(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusOK, message, data)
}

// func Created(w http.ResponseWriter, message string, data interface{}) {
// 	JSON(w, http.StatusCreated, message, data)
// }

// func NoContent(w http.ResponseWriter) {
// 	w.WriteHeader(http.StatusNoContent)
// }

// func BadRequest(w http.ResponseWriter, message string) {
// 	JSON(w, http.StatusBadRequest, message, nil)
// }

// func Unauthorized(w http.ResponseWriter, message string) {
// 	JSON(w, http.StatusUnauthorized, message, nil)
// }

// func Forbidden(w http.ResponseWriter, message string) {
// 	JSON(w, http.StatusForbidden, message, nil)
// }

// func NotFound(w http.ResponseWriter, message string) {
// 	JSON(w, http.StatusNotFound, message, nil)
// }

func InternalServerError(w http.ResponseWriter, message string) {
	JSON(w, http.StatusInternalServerError, message, nil)
}
