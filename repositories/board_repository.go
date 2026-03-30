package repositories

import (
	"hadeboard-be/config"
	"hadeboard-be/internal/models"
	"time"
)

type BoardRepository interface {
	Create(board *models.Board) error
	Update(board *models.Board) error
	FindByPublicID(publicID string) (*models.Board, error)
	AddMember(boardID uint, userIDs []uint) error
	RemoveMembers(boardID uint, userIDs []uint) error
}

type boardRepository struct{}

func NewBoardRepository() BoardRepository {
	return &boardRepository{}
}

func (r *boardRepository) Create(board *models.Board) error {
	return config.DB.Create(board).Error
}

func (r *boardRepository) Update(board *models.Board) error {
	return config.DB.Model(&models.Board{}).Where("public_id = ?", board.PublicID).Updates(map[string]interface{}{
		"title":       board.Title,
		"description": board.Description,
		"due_date":    board.DueDate,
	}).Error
}

func (r *boardRepository) FindByPublicID(publicID string) (*models.Board, error) {
	var board models.Board

	err := config.DB.Where("public_id = ?", publicID).First(&board).Error
	return &board, err
}

func (r *boardRepository) AddMember(boardID uint, userIDs []uint) error {
	if len(userIDs) == 0 {
		return nil
	}

	dateNow := time.Now()
	var members []models.BoardMember
	for _, userID := range userIDs {
		members = append(members, models.BoardMember{
			BoardInternalID: int64(boardID),
			UserInternalID:  int64(userID),
			JoinedAt:        dateNow,
		})
	}

	return config.DB.Create(&members).Error
}

func (r *boardRepository) RemoveMembers(boardID uint, userIDs []uint) error {
	if len(userIDs) == 0 {
		return nil
	}

	return config.DB.
		Where("board_internal_id = ? AND user_internal_id IN (?)", boardID, userIDs).
		Delete(&models.BoardMember{}).Error
}
