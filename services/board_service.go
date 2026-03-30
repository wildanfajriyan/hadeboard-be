package services

import (
	"errors"
	"hadeboard-be/internal/models"
	"hadeboard-be/repositories"

	"github.com/google/uuid"
)

type BoardService interface {
	Create(board *models.Board) error
	Update(board *models.Board) error
	FindByPublicID(publicID string) (*models.Board, error)
	AddMember(boardID string, userIDs []string) error
	RemoveMembers(boardID string, userIDs []string) error
}

type boardService struct {
	boardRepository       repositories.BoardRepository
	userRepository        repositories.UserRepository
	boardMemberRepository repositories.BoardMemberRepository
}

func NewBoardService(
	boardRepository repositories.BoardRepository,
	userRepository repositories.UserRepository,
	boardMemberRepository repositories.BoardMemberRepository) BoardService {
	return &boardService{boardRepository, userRepository, boardMemberRepository}
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

func (s *boardService) Update(board *models.Board) error {
	return s.boardRepository.Update(board)
}

func (s *boardService) FindByPublicID(publicID string) (*models.Board, error) {
	return s.boardRepository.FindByPublicID(publicID)
}

func (s *boardService) AddMember(boardPublicID string, userPublicIDs []string) error {
	board, err := s.boardRepository.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("Board not found")
	}

	var userInternalIDs []uint
	for _, userPublicId := range userPublicIDs {
		user, err := s.userRepository.FindByPublicID(userPublicId)
		if err != nil {
			return errors.New("User not found: " + userPublicId)
		}
		userInternalIDs = append(userInternalIDs, uint(user.InternalID))
	}

	exisitingMembers, err := s.boardMemberRepository.GetMembers(string(board.PublicID.String()))
	if err != nil {
		return err
	}

	memberMap := make(map[uint]bool)
	for _, member := range exisitingMembers {
		memberMap[uint(member.InternalID)] = true
	}

	var newMemberIDs []uint
	for _, userID := range userInternalIDs {
		if !memberMap[userID] {
			newMemberIDs = append(newMemberIDs, userID)
		}
	}

	if len(newMemberIDs) == 0 {
		return nil
	}

	return s.boardRepository.AddMember(uint(board.InternalID), newMemberIDs)
}

func (s *boardService) RemoveMembers(boardPublicID string, userPublicIDs []string) error {
	board, err := s.boardRepository.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("Board not found")
	}

	var userInternalIDs []uint
	for _, userPublicId := range userPublicIDs {
		user, err := s.userRepository.FindByPublicID(userPublicId)
		if err != nil {
			return errors.New("User not found: " + userPublicId)
		}
		userInternalIDs = append(userInternalIDs, uint(user.InternalID))
	}

	exisitingMembers, err := s.boardMemberRepository.GetMembers(string(board.PublicID.String()))
	if err != nil {
		return err
	}

	memberMap := make(map[uint]bool)
	for _, member := range exisitingMembers {
		memberMap[uint(member.InternalID)] = true
	}

	var membersToRemove []uint
	for _, userID := range userInternalIDs {
		if memberMap[userID] {
			membersToRemove = append(membersToRemove, userID)
		}
	}

	return s.boardRepository.RemoveMembers(uint(board.InternalID), membersToRemove)
}
