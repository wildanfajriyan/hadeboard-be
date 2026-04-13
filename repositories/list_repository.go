package repositories

import (
	"context"
	"hadeboard-be/config"
	"hadeboard-be/internal/models"
	"time"

	"github.com/google/uuid"
)

type ListRepository interface {
	Create(list *models.List) error
	Update(list *models.List) error
	Delete(id uint) error
	UpdatePosition(boardPublicID string, position []string) error
	GetCardPosition(listPublicID string) ([]uuid.UUID, error)
	FindByBoardID(boardID string) ([]models.List, error)
	FindByPublicID(publicID string) (*models.List, error)
	FindByID(id uint) (*models.List, error)
}

type listRepository struct {
}

func NewListRepository() ListRepository {
	return &listRepository{}
}

func (l *listRepository) Create(list *models.List) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return config.DB.WithContext(ctx).Create(list).Error
}

func (l *listRepository) Update(list *models.List) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return config.DB.WithContext(ctx).Model(&models.List{}).
		Where("public_id = ?", list.PublicID).
		Updates(map[string]interface{}{
			"title": list.Title,
		}).Error
}

func (l *listRepository) Delete(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return config.DB.WithContext(ctx).Delete(&models.List{}, id).Error
}

func (l *listRepository) UpdatePosition(boardPublicID string, position []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return config.DB.WithContext(ctx).Model(&models.ListPosition{}).
		Where("board_internal_id = (SELECT internal_id FROM boards WHERE public_id = ?)", boardPublicID).
		Update("list_order", position).Error
}

func (l *listRepository) GetCardPosition(listPublicID string) ([]uuid.UUID, error) {
	var position models.CardPosition

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := config.DB.WithContext(ctx).Joins("JOIN lists ON list.internal_id = card_positions.list_internal_id").
		Where("lists.public_id = ?", listPublicID).Error

	return position.CardOrder, err
}

func (l *listRepository) FindByBoardID(boardID string) ([]models.List, error) {
	var list []models.List

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := config.DB.WithContext(ctx).Where("board_public_id = ?", boardID).Order("internal_id ASC").Find(&list).Error

	return list, err
}

func (l *listRepository) FindByPublicID(publicID string) (*models.List, error) {
	var list models.List

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := config.DB.WithContext(ctx).Where("public_id = ?", publicID).First(&list).Error

	return &list, err
}

func (l *listRepository) FindByID(id uint) (*models.List, error) {
	var list models.List

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := config.DB.WithContext(ctx).First(&list, id).Error

	return &list, err
}
