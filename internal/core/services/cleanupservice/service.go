package cleanupservice

import (
	"context"
	"filesms/internal/core/ports"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type CleanupService struct {
	fileRepo      ports.FileRepository
	storageDir    string
	checkInterval time.Duration
}

func NewCleanupService(fileRepo ports.FileRepository, storageDir string, checkInterval time.Duration) *CleanupService {
	return &CleanupService{
		fileRepo:      fileRepo,
		storageDir:    storageDir,
		checkInterval: checkInterval,
	}
}

func (s *CleanupService) Start(ctx context.Context) {
	fmt.Println("Starting cleanup service...")
	ticker := time.NewTicker(s.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.cleanupExpiredFiles(ctx)
		}
	}
}

func (s *CleanupService) cleanupExpiredFiles(ctx context.Context) {
	// fmt.Println("Checking for expired files...(every 10 seconds)")
	expiredFiles, err := s.fileRepo.GetExpiredFiles(ctx)
	if err != nil {
		log.Printf("Error getting expired files: %v", err)
		return
	}

	var filesToDelete []uuid.UUID
	for _, file := range expiredFiles {
		fmt.Println("Deleting expired file:", file.URL)
		// Delete file from local storage
		err := os.Remove(filepath.Join(s.storageDir, file.URL))
		if err != nil {
			log.Printf("Error deleting file %s: %v", file.URL, err)
			continue
		}
		filesToDelete = append(filesToDelete, file.ID)
	}

	// Delete files from database
	if len(filesToDelete) > 0 {
		err := s.fileRepo.DeleteFiles(ctx, filesToDelete)
		if err != nil {
			log.Printf("Error deleting files from database: %v", err)
		} else {
			log.Printf("Successfully deleted %d expired files", len(filesToDelete))
		}
	}
}
