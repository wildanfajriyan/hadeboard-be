package repositories

import (
	"context"
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
	GetMyBoardPaginate(userPublicID, filter, sort string, limit, offset int) ([]models.Board, int64, error)
}

type boardRepository struct{}

func NewBoardRepository() BoardRepository {
	return &boardRepository{}
}

func (r *boardRepository) Create(board *models.Board) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return config.DB.WithContext(ctx).Create(board).Error
}

func (r *boardRepository) Update(board *models.Board) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return config.DB.WithContext(ctx).Model(&models.Board{}).Where("public_id = ?", board.PublicID).Updates(map[string]interface{}{
		"title":       board.Title,
		"description": board.Description,
		"due_date":    board.DueDate,
	}).Error
}

func (r *boardRepository) FindByPublicID(publicID string) (*models.Board, error) {
	var board models.Board

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := config.DB.WithContext(ctx).Where("public_id = ?", publicID).First(&board).Error
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

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return config.DB.WithContext(ctx).Create(&members).Error
}

func (r *boardRepository) RemoveMembers(boardID uint, userIDs []uint) error {
	if len(userIDs) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return config.DB.
		WithContext(ctx).
		Where("board_internal_id = ? AND user_internal_id IN (?)", boardID, userIDs).
		Delete(&models.BoardMember{}).Error
}

func (r *boardRepository) GetMyBoardPaginate(userPublicID string, filter string, sort string, limit int, offset int) ([]models.Board, int64, error) {
	var board []models.Board
	var total int64

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := config.DB.WithContext(ctx).Model(&models.Board{}).
		Where("owner_public_id = ? OR internal_id IN ("+
			"SELECT board_members.board_internal_id FROM board_members "+
			"JOIN users ON users.internal_id = board_members.user_internal_id "+
			"WHERE users.public_id = ?)", userPublicID, userPublicID)

	if filter != "" {
		query = query.Where("title ILIKE ?", "%"+filter+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sort != "" {
		query = query.Order(sort)
	} else {
		query = query.Order("created_at desc")
	}

	if err := query.Limit(limit).Offset(offset).Find(&board).Error; err != nil {
		return nil, 0, err
	}

	return board, total, nil
}
