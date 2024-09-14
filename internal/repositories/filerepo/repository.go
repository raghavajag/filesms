package filerepo

import (
	"context"
	"database/sql"
	"errors"
	"filesms/internal/core/domain"
	"fmt"
	"strings"
	"time"

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
func (r *postgresFileRepository) SaveSharedFileURL(ctx context.Context, sharedFileURL *domain.SharedFileURL) error {
	query := `INSERT INTO shared_file_urls (file_id, url, expires_at, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, sharedFileURL.FileID, sharedFileURL.URL, sharedFileURL.ExpiresAt, sharedFileURL.CreatedAt)
	return err
}
func (r *postgresFileRepository) GetFileIDBySharedURL(ctx context.Context, url string) (uuid.UUID, error) {
	query := `SELECT file_id FROM shared_file_urls WHERE url = $1 AND expires_at > NOW()`
	var fileID uuid.UUID
	err := r.db.QueryRowContext(ctx, query, url).Scan(&fileID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, errors.New("shared URL not found or expired")
		}
		return uuid.Nil, err
	}
	return fileID, nil
}
func (r *postgresFileRepository) Search(ctx context.Context, userID uuid.UUID, params domain.FileSearchParams) ([]*domain.File, error) {
	query := `
		SELECT id, user_id, name, size, type, url, created_at, updated_at
		FROM files
		WHERE user_id = $1
	`
	args := []interface{}{userID}
	argCount := 1

	if params.Query != "" {
		argCount++
		query += fmt.Sprintf(" AND name ILIKE $%d", argCount)
		args = append(args, "%"+params.Query+"%")
	}

	if params.FileType != "" {
		argCount++
		query += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, params.FileType)
	}

	if !params.FromDate.IsZero() {
		argCount++
		query += fmt.Sprintf(" AND created_at >= $%d", argCount)
		args = append(args, params.FromDate)
	}

	if !params.ToDate.IsZero() {
		argCount++
		query += fmt.Sprintf(" AND created_at <= $%d", argCount)
		args = append(args, params.ToDate)
	}

	// Add sorting
	if params.SortBy != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", params.SortBy, params.SortDir)
	} else {
		query += " ORDER BY created_at DESC"
	}

	// Add pagination
	if params.Limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, params.Limit)
	}
	if params.Offset > 0 {
		argCount++
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, params.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*domain.File
	for rows.Next() {
		var file domain.File
		err := rows.Scan(&file.ID, &file.UserID, &file.Name, &file.Size, &file.Type, &file.URL, &file.CreatedAt, &file.UpdatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}

	return files, nil
}
func (r *postgresFileRepository) DeleteFiles(ctx context.Context, fileIDs []uuid.UUID) error {
	// Convert UUID slice to PostgreSQL array format
	pgArray := convertUUIDsToPGArray(fileIDs)

	query := `DELETE FROM files WHERE id = ANY($1)`
	_, err := r.db.ExecContext(ctx, query, pgArray)
	return err
}

// Helper function to convert UUID slice to PostgreSQL array format
func convertUUIDsToPGArray(uuids []uuid.UUID) string {
	uuidStrings := make([]string, len(uuids))
	for i, uuid := range uuids {
		uuidStrings[i] = uuid.String()
	}
	return fmt.Sprintf("{%s}", strings.Join(uuidStrings, ","))
}
func (r *postgresFileRepository) GetExpiredFiles(ctx context.Context) ([]*domain.File, error) {
	query := `
        SELECT id, user_id, name, size, type, url, expiration_date, created_at, updated_at
        FROM files
        WHERE expiration_date < $1
    `
	rows, err := r.db.QueryContext(ctx, query, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*domain.File
	for rows.Next() {
		var file domain.File
		err := rows.Scan(
			&file.ID, &file.UserID, &file.Name, &file.Size, &file.Type, &file.URL,
			&file.ExpirationDate, &file.CreatedAt, &file.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}
	return files, nil
}
