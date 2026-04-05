package models

import "github.com/google/uuid"

type UserAssignee struct {
	InternalID int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID   uuid.UUID `json:"public_id" db:"public_id"`
	Name       string    `json:"name" db:"name"`
	Email      string    `json:"email" db:"email" gorm:"unique"`
}

type CardAssignee struct {
	CardInternalID int64 `json:"card_internal_id" db:"card_internal_id"`
	UserInternalID int64 `json:"user_internal_id" db:"user_internal_id"`

	User UserAssignee `json:"user" gorm:"foreignKey:UserInternalID;references:InternalID"`
}

func (UserAssignee) TableName() string {
	return "users"
}
