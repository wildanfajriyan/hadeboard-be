package services

import (
	"errors"
	"hadeboard-be/internal/models"
	"hadeboard-be/repositories"
	"hadeboard-be/utils"

	"github.com/google/uuid"
)

type UserService interface {
	Register(user *models.User) error
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{userRepository}
}

func (s *userService) Register(user *models.User) error {
	existingUser, _ := s.userRepository.FindByEmail(user.Email)
	if existingUser.InternalID != 0 {
		return errors.New("Email already registered")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	user.Role = "user"
	user.PublicID = uuid.New()

	return s.userRepository.Create(user)
}
