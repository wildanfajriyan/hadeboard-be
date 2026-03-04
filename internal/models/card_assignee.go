package models

type CardAssignee struct {
	CardInternalID int64 `json:"card_internal_id" db:"card_internal_id"`
	UserInternalID int64 `json:"user_internal_id" db:"user_internal_id"`
}
