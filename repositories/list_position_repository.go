package repositories

import (
	"hadeboard-be/config"
	"hadeboard-be/internal/models"

	"github.com/google/uuid"
)

type ListPositionRepository interface {
	GetByBoard(boardPublicID string) (*models.ListPosition, error)
	CreateOrUpdate(boardPublicID string, listOrder []uuid.UUID) error
	GetListOrder(boardPublicID string) ([]uuid.UUID, error)
	UpdateListOrder(position *models.ListPosition) error
}

type listPositionRepository struct{}

func NewListPositionRepository() ListPositionRepository {
	return &listPositionRepository{}
}

func (l *listPositionRepository) GetByBoard(boardPublicID string) (*models.ListPosition, error) {
	var position models.ListPosition

	err := config.DB.Joins("JOIN boards ON boards.internal_id = list_positions.board_internal_id").
		Where("boards.public_id = ?", boardPublicID).Error

	return &position, err
}

func (l *listPositionRepository) CreateOrUpdate(boardPublicID string, listOrder []uuid.UUID) error {
	return config.DB.Exec(`
		INSERT INTO list_positions (board_internal_id, list_order)
		SELECT internal_id, ? FROM boards WHERE public_id = ?
		ON CONFLICT (board_internal_id)
		DO UPDATE SET list_order = EXCLUDE.list_order
	`, listOrder, boardPublicID).Error
}

func (l *listPositionRepository) GetListOrder(boardPublicID string) ([]uuid.UUID, error) {
	position, err := l.GetByBoard(boardPublicID)
	if err != nil {
		return nil, err
	}

	return position.ListOrder, err
}

func (l *listPositionRepository) UpdateListOrder(position *models.ListPosition) error {
	return config.DB.Model(position).
		Where("internal_id = ?", position.InternalID).
		Update("list_order", position.ListOrder).Error
}
