package domain

import (
	"time"

	"github.com/google/uuid"
)

type SharedFileURL struct {
	FileID    uuid.UUID `json:"file_id"`
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
