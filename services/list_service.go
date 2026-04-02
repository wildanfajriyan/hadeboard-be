package services

import (
	"errors"
	"fmt"
	"hadeboard-be/config"
	"hadeboard-be/internal/models"
	"hadeboard-be/internal/models/types"
	"hadeboard-be/repositories"
	"hadeboard-be/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ListWithOrder struct {
	Position []uuid.UUID
	Lists    []models.List
}

type ListService interface {
	GetByBoardID(boardPublicID string) (*ListWithOrder, error)
	GetByID(id uint) (*models.List, error)
	GetByPublicID(publicID string) (*models.List, error)
	Create(list *models.List) error
	Update(list *models.List) error
	Delete(id uint) error
	UpdatePosition(boardPublicID string, position []uuid.UUID) error
}

type listService struct {
	listRepository         repositories.ListRepository
	boardRepository        repositories.BoardRepository
	listPositionRepository repositories.ListPositionRepository
}

func NewListService(
	listRepository repositories.ListRepository,
	boardRepository repositories.BoardRepository,
	listPositionRepository repositories.ListPositionRepository) ListService {
	return &listService{listRepository, boardRepository, listPositionRepository}
}

func (l *listService) Create(list *models.List) error {
	board, err := l.boardRepository.FindByPublicID(list.BoardPublicId.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Board not found")
		}
		return fmt.Errorf("Failed to get board: %w", err)
	}
	list.BoardInternalID = board.InternalID

	if list.PublicID == uuid.Nil {
		list.PublicID = uuid.New()
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(list).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to create list : %w", err)
	}

	var position models.ListPosition
	res := tx.Where("board_internal_id = ?", board.InternalID).First(&position)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		position = models.ListPosition{
			PublicID:        uuid.New(),
			BoardInternalID: board.InternalID,
			ListOrder:       types.UUIDArray{list.PublicID},
		}
		if err := tx.Create(&position).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to create list position: %w", err)
		}
	} else if res.Error != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to create list : %w", res.Error)
	} else {
		position.ListOrder = append(position.ListOrder, list.PublicID)
		if err := tx.Model(&position).Update("list_order", position.ListOrder).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to update list position: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Transaction commit failed: %w", err)
	}

	return nil
}

func (l *listService) Delete(id uint) error {
	return l.listRepository.Delete(id)
}

func (l *listService) GetByBoardID(boardPublicID string) (*ListWithOrder, error) {
	_, err := l.boardRepository.FindByPublicID(boardPublicID)
	if err != nil {
		return nil, errors.New("Board not found")
	}

	position, err := l.listPositionRepository.GetListOrder(boardPublicID)
	if err != nil {
		return nil, errors.New("Failed to get list order : " + err.Error())
	}

	lists, err := l.listRepository.FindByBoardID(boardPublicID)
	if err != nil {
		return nil, errors.New("Failed to get list : " + err.Error())
	}

	orderedList := utils.SortingListByPosition(lists, position)

	return &ListWithOrder{
		Position: position,
		Lists:    orderedList,
	}, nil
}

func (l *listService) GetByID(id uint) (*models.List, error) {
	return l.listRepository.FindByID(id)
}

func (l *listService) GetByPublicID(publicID string) (*models.List, error) {
	return l.listRepository.FindByPublicID(publicID)
}

func (l *listService) Update(list *models.List) error {
	return l.listRepository.Update(list)
}

func (l *listService) UpdatePosition(boardPublicID string, positions []uuid.UUID) error {
	board, err := l.boardRepository.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("Board not found")
	}

	position, err := l.listPositionRepository.GetByBoard(board.PublicID.String())
	if err != nil {
		return errors.New("List position not found")
	}

	position.ListOrder = positions
	return l.listPositionRepository.UpdateListOrder(position)
}
