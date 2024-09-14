package domain

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID             uuid.UUID `json:"id" validate:"required,uuid4"`
	Name           string    `json:"name" validate:"required,min=1,max=255"`
	Size           int64     `json:"size" validate:"required,gt=0"`
	UserID         uuid.UUID `json:"user_id" validate:"required,uuid4"`
	Type           string    `json:"type" validate:"required,min=1,max=255"`
	URL            string    `json:"url" validate:"required,min=1,max=255"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	ExpirationDate time.Time `json:"expiration_date"`
}
