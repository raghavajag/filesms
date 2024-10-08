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
	SaveSharedFileURL(ctx context.Context, sharedFileURL *domain.SharedFileURL) error
	Search(ctx context.Context, userID uuid.UUID, params domain.FileSearchParams) ([]*domain.File, error)
	GetExpiredFiles(ctx context.Context) ([]*domain.File, error)
	DeleteFiles(ctx context.Context, fileIDs []uuid.UUID) error
	// Update(ctx context.Context, file *domain.File) error
	// Delete(ctx context.Context, id uuid.UUID) error
}
