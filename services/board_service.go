package services

import (
	"errors"
	"hadeboard-be/internal/models"
	"hadeboard-be/repositories"

	"github.com/google/uuid"
)

type BoardService interface {
	Create(board *models.Board) error
}

type boardService struct {
	boardRepository repositories.BoardRepository
	userRepository  repositories.UserRepository
}

func NewBoardService(boardRepository repositories.BoardRepository, userRepository repositories.UserRepository) BoardService {
	return &boardService{boardRepository, userRepository}
}

func (s *boardService) Create(board *models.Board) error {
	user, err := s.userRepository.FindByPublicID(board.OwnerPublicID.String())
	if err != nil {
		return errors.New("owner not found")
	}

	board.PublicID = uuid.New()
	board.OwnerInternalID = user.InternalID
	return s.boardRepository.Create(board)
}
