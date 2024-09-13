package ports

import (
	"context"
	"filesms/internal/core/domain"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uint) error
}

type FileRepository interface {
	Create(ctx context.Context, file *domain.File) error
	GetByID(ctx context.Context, id uint) (*domain.File, error)
	GetByUserID(ctx context.Context, userID uint) ([]*domain.File, error)
	Update(ctx context.Context, file *domain.File) error
	Delete(ctx context.Context, id uint) error
	Search(ctx context.Context, userID uint, query string, fileType string, fromDate, toDate time.Time) ([]*domain.File, error)
	GetExpiredFiles(ctx context.Context, expirationDate time.Time) ([]*domain.File, error)
}
