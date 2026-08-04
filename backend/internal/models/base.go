package models

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID uuid.UUID `db:"id" json:"id"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`

	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}