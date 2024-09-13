package ports

import "filesms/internal/core/domain"

type AuthService interface {
	Login(email, password string) (string, error)
	Register(email, password string) (string, error)
}

type FileService interface {
	UploadFile(ownerID string, file *domain.File) error
	GetFilesByOwner(ownerID string) ([]*domain.File, error)
	ShareFile(fileID string) (string, error)
}
