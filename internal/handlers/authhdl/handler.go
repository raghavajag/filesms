package authhdl

import (
	"encoding/json"
	"filesms/internal/core/services/authsrv"
	"filesms/pkg/errors"
	"filesms/pkg/middleware"
	"filesms/pkg/validation"
	"net/http"

	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *authsrv.AuthService
}

type AuthInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func NewAuthHandler(authService *authsrv.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	var input AuthInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return errors.NewAPIError(http.StatusBadRequest, "Invalid JSON", nil)
	}

	if err := validation.ValidateStruct(input); err != nil {
		return err
	}

	user, err := h.authService.Register(r.Context(), input.Email, input.Password)
	if err != nil {
		return errors.NewAPIError(http.StatusInternalServerError, "Failed to register user", nil)
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	var input AuthInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return errors.NewAPIError(http.StatusBadRequest, "Invalid JSON", nil)
	}

	if err := validation.ValidateStruct(input); err != nil {
		return err
	}

	token, err := h.authService.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		return errors.NewAPIError(http.StatusUnauthorized, "Invalid credentials", nil)
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(map[string]string{"token": token})
}
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) error {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)

	user, err := h.authService.Me(r.Context(), userID)
	if err != nil {
		return errors.NewAPIError(http.StatusInternalServerError, "Failed to fetch user", nil)
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(user)
}
