package services

import (
	"errors"
	"fmt"
	"hadeboard-be/config"
	"hadeboard-be/internal/models"
	"hadeboard-be/internal/models/types"
	"hadeboard-be/repositories"
	"sort"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CardService interface {
	Create(card *models.Card, listPublicID string) error
	Update(card *models.Card, listPublicID string) error
	Delete(id uint) error
	GeyByID(id uint) (*models.Card, error)
	GetByPublicID(publicID string) (*models.Card, error)
	GetByListID(listPublicID string) ([]models.Card, error)
}

type cardService struct {
	cardRepository repositories.CardRepository
	listRepository repositories.ListRepository
	userRepository repositories.UserRepository
}

func NewCardService(
	cardRepository repositories.CardRepository,
	listRepository repositories.ListRepository,
	userRepository repositories.UserRepository,
) CardService {
	return &cardService{cardRepository, listRepository, userRepository}
}

func SortCardByPosition(cards []models.Card, order []uuid.UUID) []models.Card {
	orderMap := make(map[uuid.UUID]int)
	for i, id := range order {
		orderMap[id] = i
	}

	defaultIdx := len(order)
	sort.SliceStable(cards, func(i, j int) bool {
		idxI, okI := orderMap[cards[i].PublicID]
		if !okI {
			idxI = defaultIdx
		}

		idxJ, okJ := orderMap[cards[j].PublicID]
		if !okJ {
			idxJ = defaultIdx
		}

		if idxI == idxJ {
			return cards[i].CreatedAt.Before(cards[j].CreatedAt)
		}

		return idxI < idxJ
	})

	return cards
}

func (c *cardService) Create(card *models.Card, listPublicID string) error {
	list, err := c.listRepository.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("List not found: %w", err)
	}
	card.ListInternalID = list.InternalID

	if card.PublicID == uuid.Nil {
		card.PublicID = uuid.New()
	}
	card.CreatedAt = time.Now()

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := tx.Create(card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to create card: %w", err)
	}

	var position models.CardPosition
	if err := tx.Model(&models.CardPosition{}).
		Where("list_internal_id = ?", list.InternalID).
		First(&position).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			position = models.CardPosition{
				PublicID:       uuid.New(),
				ListInternalID: list.InternalID,
				CardOrder:      types.UUIDArray{card.PublicID},
			}

			if err := tx.Create(&position).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("Failed to create card position: %w", err)
			}
		} else {
			tx.Rollback()
			return fmt.Errorf("Failed to get card position: %w", err)
		}
	} else {
		position.CardOrder = append(position.CardOrder, card.PublicID)
		if err := tx.Model(&models.CardPosition{}).
			Where("internal_id = ?", position.InternalID).
			Update("card_order", position.CardOrder).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to update card position: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Transaction commit failed: %w", err)
	}

	return nil
}

func (c *cardService) Update(card *models.Card, listPublicID string) error {
	existingCard, err := c.cardRepository.FindByPublicID(card.PublicID.String())
	if err != nil {
		return fmt.Errorf("Card not found: %w", err)
	}

	newList, err := c.listRepository.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("List not found: %w", err)
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if existingCard.ListInternalID != newList.InternalID {
		var oldCardPosition models.CardPosition
		if err := tx.Where("list_internal_id = ?", existingCard.ListInternalID).First(&oldCardPosition).Error; err != nil {
			filtered := make(types.UUIDArray, 0, len(oldCardPosition.CardOrder))
			for _, id := range oldCardPosition.CardOrder {
				if id != existingCard.PublicID {
					filtered = append(filtered, id)
				}
			}

			if err := tx.Model(&models.CardPosition{}).
				Where("internal_id = ?", oldCardPosition.InternalID).
				Update("card_order", types.UUIDArray(filtered)).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("Failed to update old card position: %w", err)
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return fmt.Errorf("Failed to get old card position: %w", err)
		}

		var newCardPosition models.CardPosition
		res := tx.Where("list_internal_id = ?", newList.InternalID).First(&newCardPosition)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			newCardPosition = models.CardPosition{
				PublicID:       uuid.New(),
				ListInternalID: newList.InternalID,
				CardOrder:      types.UUIDArray{existingCard.PublicID},
			}

			if err := tx.Create(&newCardPosition).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("Failed to create card position for new list: %w", err)
			}
		} else if res.Error == nil {
			updateOrder := append(newCardPosition.CardOrder, existingCard.PublicID)
			if err := tx.Model(&models.CardPosition{}).
				Where("internal_id = ?", newCardPosition.InternalID).
				Update("card_order", types.UUIDArray(updateOrder)).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("Failed to update new card position: %w", err)
			}
		} else {
			tx.Rollback()
			return fmt.Errorf("Failed to get new card position: %w", res.Error)
		}
	}

	card.InternalID = existingCard.InternalID
	card.PublicID = existingCard.PublicID
	card.ListInternalID = existingCard.ListInternalID

	if err := tx.Save(card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to update card: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Transaction commit failed: %w", err)
	}

	return nil
}

func (c *cardService) Delete(id uint) error {
	return c.cardRepository.Delete(id)
}

func (c *cardService) GetByListID(listPublicID string) ([]models.Card, error) {
	list, err := c.listRepository.FindByPublicID(listPublicID)
	if err != nil {
		return nil, fmt.Errorf("List not found: %w", err)
	}

	position, err := c.cardRepository.FindCardPositionByListID(list.InternalID)
	if err != nil {
		return nil, fmt.Errorf("Card position not found: %w", err)
	}

	cards, err := c.cardRepository.FindByListID(listPublicID)
	if err != nil {
		return nil, fmt.Errorf("Cards not found: %w", err)
	}

	if position != nil && len(position.CardOrder) > 0 {
		cards = SortCardByPosition(cards, position.CardOrder)
	}

	return cards, nil
}

func (c *cardService) GetByPublicID(publicID string) (*models.Card, error) {
	return c.cardRepository.FindByPublicID(publicID)
}

func (c *cardService) GeyByID(id uint) (*models.Card, error) {
	return c.cardRepository.FindByID(id)
}
