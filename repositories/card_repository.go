package repositories

import (
	"fmt"
	"hadeboard-be/config"
	"hadeboard-be/internal/models"
	"path/filepath"

	"gorm.io/gorm"
)

type cardRepository struct{}

type CardRepository interface {
	Create(card *models.Card) error
	Update(card *models.Card) error
	Delete(id uint) error
	FindByID(id uint) (*models.Card, error)
	FindByPublicID(publicID string) (*models.Card, error)
	FindByListID(listID string) ([]models.Card, error)
	FindCardPositionByListID(id int64) (*models.CardPosition, error)
	UpdatePosition(listID string, position []string) error
}

func NewCardRepository() CardRepository {
	return &cardRepository{}
}

func (c *cardRepository) Create(card *models.Card) error {
	return config.DB.Create(card).Error
}

func (c *cardRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Card{}, id).Error
}

func (c *cardRepository) Update(card *models.Card) error {
	return config.DB.Save(card).Error
}

func (c *cardRepository) FindByID(id uint) (*models.Card, error) {
	var card models.Card
	err := config.DB.Preload("Labels").Preload("Assigness").First(&card, id).Error

	return &card, err
}

func (c *cardRepository) FindByPublicID(publicID string) (*models.Card, error) {
	var card models.Card
	if err := config.DB.
		Preload("Assignees.User", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("internal_id", "public_id", "name", "email")
		}).
		Preload("Attachments").
		Where("public_id = ?", publicID).
		First(&card).Error; err != nil {
		return nil, err
	}

	baseUrl := config.AppConfig.AppUrl

	for i := range card.Attachments {
		card.Attachments[i].FileURL = fmt.Sprintf("%s/files/%s", baseUrl, filepath.Base(card.Attachments[i].File))
	}

	return &card, nil
}

func (c *cardRepository) FindByListID(listID string) ([]models.Card, error) {
	var cards []models.Card
	err := config.DB.Joins("JOIN lists ON lists.internal_id = cards.list_internal_id").
		Where("lists.public_id = ?", listID).
		Order("position ASC").
		Find(&cards).Error

	return cards, err
}

func (c *cardRepository) FindCardPositionByListID(id int64) (*models.CardPosition, error) {
	var position models.CardPosition

	err := config.DB.Where("list_internal_id = ?", id).First(&position).Error
	if err != nil {
		return nil, err
	}

	return &position, err
}

func (c *cardRepository) UpdatePosition(listID string, position []string) error {
	return config.DB.Model(&models.CardPosition{}).
		Where("list_internal_id = (SELECT internal_id FROM lists WHERE public_id = ?)", listID).
		Update("card_order", position).Error
}
