package models

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	InternalID     int64      `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID       uuid.UUID  `json:"public_id" db:"public_id"`
	ListInternalID int64      `json:"list_internal_id" db:"list_internal_id"`
	Title          string     `json:"title" db:"title"`
	Description    string     `json:"description" db:"description"`
	Position       int        `json:"position" db:"position"`
	DueDate        *time.Time `json:"due_date,omitempty" db:"due_date"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	// UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`

	Assignees   []CardAssignee   `json:"assignees,omitempty" gorm:"foreignKey:CardInternalID;references:InternalID"`
	Attachments []CardAttachment `json:"attachments,omitempty" gorm:"foreignKey:CardInternalID;references:InternalID"`
	Labels      []CardLabel      `json:"labels,omitempty" gorm:"foreignKey:CardInternalID;references:InternalID"`
}
