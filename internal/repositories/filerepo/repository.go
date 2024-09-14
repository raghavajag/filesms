package filerepo

import (
	"context"
	"database/sql"
	"filesms/internal/core/domain"
)

type postgresFileRepository struct {
	db *sql.DB
}

func NewPostgresFileRepository(db *sql.DB) *postgresFileRepository {
	return &postgresFileRepository{db: db}
}

func (r *postgresFileRepository) Create(ctx context.Context, file *domain.File) error {
	query := `INSERT INTO files (id, user_id, name, size, type, url, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, file.ID, file.UserID, file.Name, file.Size, file.Type, file.URL, file.CreatedAt, file.UpdatedAt)
	return err
}
