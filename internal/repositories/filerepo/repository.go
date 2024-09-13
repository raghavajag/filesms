package filerepo

import (
	"database/sql"
)

type postgresFileRepository struct {
	db *sql.DB
}

func NewPostgresFileRepository(db *sql.DB) *postgresFileRepository {
	return &postgresFileRepository{db: db}
}
