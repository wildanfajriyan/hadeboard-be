package models

import (
	"time"

	"github.com/google/uuid"
)

type Label struct {
	InternalID int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID   uuid.UUID `json:"public_id" db:"public_id"`
	Name       string    `json:"name" db:"name"`
	Color      string    `json:"color" db:"color"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
