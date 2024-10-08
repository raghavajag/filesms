package filehdl

import (
	"filesms/internal/core/domain"
	"filesms/internal/core/services/filesrv"
	response "filesms/pkg/api"
	"filesms/pkg/errors"
	"filesms/pkg/middleware"
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
	response.Success(w, "File uploaded successfully", uploadedFile)
	return nil
}
func (h *FileHandler) GetFiles(w http.ResponseWriter, r *http.Request) error {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)

	files, err := h.fileService.GetFiles(r.Context(), userID)
	if err != nil {
		return errors.NewAPIError(http.StatusInternalServerError, "Failed to get files", err)
	}
	if len(files) == 0 {
		response.Success(w, "No files found", []domain.File{})
		return nil
	}
	response.Success(w, "Files retrieved successfully", files)
	return nil
}
func (h *FileHandler) ShareFile(w http.ResponseWriter, r *http.Request) error {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	fileIDStr := r.URL.Query().Get("file_id")
	fileID, err := uuid.Parse(fileIDStr)
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
	response.Success(w, "File shared successfully", shareURL)
	return nil
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
	if len(files) == 0 {
		response.Success(w, "No files found", []domain.File{})
		return nil
	}
	response.Success(w, "Files retrieved successfully", files)
	return nil
}

func (h *FileHandler) GetFile(w http.ResponseWriter, r *http.Request) error {
	userId := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	fileIDStr := r.URL.Query().Get("file_id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		return errors.NewAPIError(http.StatusBadRequest, "Invalid file ID", err)
	}

	file, err := h.fileService.GetFile(r.Context(), fileID)
	if err != nil {
		return errors.NewAPIError(http.StatusInternalServerError, "Failed to get file", err)
	}

	if userId != file.UserID {
		return errors.NewAPIError(http.StatusUnauthorized, "Unauthorized access to file", nil)
	}
	response.Success(w, "File retrieved successfully", file)
	return nil
}
