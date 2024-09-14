package filesrv

import (
	"context"
	"filesms/internal/core/domain"
	"filesms/internal/core/ports"
	"filesms/pkg/storage"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type FileService struct {
	fileRepo ports.FileRepository
	storage  *storage.LocalStorage
	baseURL  string
}

func NewFileService(fileRepo ports.FileRepository, storage *storage.LocalStorage, baseURL string) *FileService {
	return &FileService{
		fileRepo: fileRepo,
		storage:  storage,
		baseURL:  baseURL,
	}
}

func (s *FileService) Upload(ctx context.Context, userID uuid.UUID, fileName string, content io.Reader, fileSize int64) (*domain.File, error) {
	// Generate a unique filename
	uniqueFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileName)
	fmt.Println(uniqueFileName)
	// Save the file to local storage
	filePath, err := s.storage.Save(uniqueFileName, content)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Create file metadata
	file := &domain.File{
		UserID:    userID,
		Name:      fileName,
		Size:      fileSize,
		Type:      filepath.Ext(fileName),
		URL:       filePath,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ID:        uuid.New(),
	}

	// Save file metadata to database
	err = s.fileRepo.Create(ctx, file)
	if err != nil {
		// If database insert fails, delete the file from storage
		_ = s.storage.Delete(uniqueFileName)
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}

	return file, nil
}
func (s *FileService) GetFiles(ctx context.Context, userID uuid.UUID) ([]*domain.File, error) {
	return s.fileRepo.GetByUserID(ctx, userID)
}

func (s *FileService) GetFileByID(ctx context.Context, fileID uuid.UUID) (*domain.File, error) {
	return s.fileRepo.GetByID(ctx, fileID)
}
