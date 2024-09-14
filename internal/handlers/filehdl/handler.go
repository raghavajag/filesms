package filehdl

import (
	"encoding/json"
	"filesms/internal/core/domain"
	"filesms/internal/core/services/filesrv"
	"filesms/pkg/errors"
	"filesms/pkg/middleware"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
func (h *FileHandler) ShareFile(w http.ResponseWriter, r *http.Request) error {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	fileIDStr := r.URL.Query().Get("file_id")
	fileID, err := uuid.Parse(fileIDStr)
	fmt.Println(fileID)
	if err != nil {
		return errors.NewAPIError(http.StatusBadRequest, "Invalid file ID", err)
	}

	expirationTime := 24 * time.Hour // Default to 24 hours
	if expStr := r.URL.Query().Get("expiration"); expStr != "" {
		expDuration, err := time.ParseDuration(expStr)
		if err == nil {
			expirationTime = expDuration
		}
	}

	shareURL, err := h.fileService.ShareFile(r.Context(), fileID, userID, expirationTime)
	if err != nil {
		return errors.NewAPIError(http.StatusInternalServerError, "Failed to share file", err)
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(map[string]string{"share_url": shareURL})
}
func (h *FileHandler) SearchFiles(w http.ResponseWriter, r *http.Request) error {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return errors.NewAPIError(http.StatusUnauthorized, "Unauthorized", nil)
	}

	params := domain.FileSearchParams{
		Query:    r.URL.Query().Get("query"),
		FileType: r.URL.Query().Get("type"),
		SortBy:   r.URL.Query().Get("sort_by"),
		SortDir:  r.URL.Query().Get("sort_dir"),
	}

	if fromDate := r.URL.Query().Get("from"); fromDate != "" {
		params.FromDate, _ = time.Parse(time.RFC3339, fromDate)
	}

	if toDate := r.URL.Query().Get("to"); toDate != "" {
		params.ToDate, _ = time.Parse(time.RFC3339, toDate)
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		params.Limit, _ = strconv.Atoi(limit)
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		params.Offset, _ = strconv.Atoi(offset)
	}

	files, err := h.fileService.SearchFiles(r.Context(), userID, params)
	if err != nil {
		return errors.NewAPIError(http.StatusInternalServerError, "Failed to search files", err)
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(files)
}
