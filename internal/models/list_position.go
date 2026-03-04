package models

import (
	"hadeboard-be/internal/models/types"

	"github.com/google/uuid"
)

type ListPosition struct {
	InternalID      int64           `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID        uuid.UUID       `json:"public_id" db:"public_id"`
	BoardInternalID int64           `json:"board_internal_id" db:"board_internal_id"`
	ListOrder       types.UUIDArray `json:"list_order"`
}
