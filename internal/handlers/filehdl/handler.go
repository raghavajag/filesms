package filehdl

import (
	"encoding/json"
	"filesms/internal/core/services/filesrv"
	"filesms/pkg/errors"
	"filesms/pkg/middleware"
	"net/http"

	"github.com/google/uuid"
)

type FileHandler struct {
	fileService *filesrv.FileService
}

func NewFileHandler(fileService *filesrv.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request) error {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)

	file, header, err := r.FormFile("file")
	if err != nil {
		return errors.NewAPIError(http.StatusBadRequest, "Failed to read file", nil)
	}
	defer file.Close()

	uploadedFile, err := h.fileService.Upload(r.Context(), userID, header.Filename, file, header.Size)
	if err != nil {
		return errors.NewAPIError(http.StatusInternalServerError, "Failed to upload file", nil)
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(uploadedFile)
}
func (h *FileHandler) GetFiles(w http.ResponseWriter, r *http.Request) error {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)

	files, err := h.fileService.GetFiles(r.Context(), userID)
	if err != nil {
		return errors.NewAPIError(http.StatusInternalServerError, "Failed to get files", nil)
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(files)
}
