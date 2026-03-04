package models

type CardLabel struct {
	CardInternalID  int64 `json:"card_internal_id" db:"card_internal_id"`
	LabelInternalID int64 `json:"label_internal_id" db:"label_internal_id"`
}
