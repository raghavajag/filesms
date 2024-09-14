package cleanupservice

import (
	"context"
	"filesms/internal/core/ports"
	"fmt"
	"time"
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

func (s *CleanupService) cleanupExpiredFiles(_ context.Context) {
	fmt.Println("Cleaning up expired files...")
}
