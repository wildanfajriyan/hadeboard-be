package models

import "time"

type BoardMember struct {
	BoardInternalID int64     `json:"board_internal_id" db:"board_internal_id" gorm:"primaryKey"`
	UserInternalID  int64     `json:"user_internal_id" db:"user_internal_id" gorm:"primaryKey"`
	JoinedAt        time.Time `json:"joined_at" db:"joined_at"`
}
