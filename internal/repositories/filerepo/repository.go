package filerepo

import (
	"context"
	"database/sql"
	"errors"
	"filesms/internal/core/domain"

	"github.com/google/uuid"
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
func (r *postgresFileRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.File, error) {
	query := `SELECT id, user_id, name, size, type, url, created_at, updated_at 
              FROM files 
              WHERE id = $1`
	var file domain.File
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&file.ID, &file.UserID, &file.Name, &file.Size, &file.Type, &file.URL, &file.CreatedAt, &file.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("file not found")
		}
		return nil, err
	}
	return &file, nil
}

func (r *postgresFileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.File, error) {
	query := `SELECT id, user_id, name, size, type, url, created_at, updated_at 
              FROM files 
              WHERE user_id = $1 
              ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*domain.File
	for rows.Next() {
		var file domain.File
		if err := rows.Scan(&file.ID, &file.UserID, &file.Name, &file.Size, &file.Type, &file.URL, &file.CreatedAt, &file.UpdatedAt); err != nil {
			return nil, err
		}
		files = append(files, &file)
	}
	return files, nil
}
