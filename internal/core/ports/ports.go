package ports

import (
	"context"
	"filesms/internal/core/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type FileRepository interface {
	Create(ctx context.Context, file *domain.File) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.File, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.File, error)
	// Update(ctx context.Context, file *domain.File) error
	// Delete(ctx context.Context, id uuid.UUID) error
	// Search(ctx context.Context, userID uuid.UUID, query string, fileType string, fromDate, toDate time.Time) ([]*domain.File, error)
	// GetExpiredFiles(ctx context.Context, expirationDate time.Time) ([]*domain.File, error)
}
