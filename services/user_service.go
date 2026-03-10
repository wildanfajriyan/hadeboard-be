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
	Login(email, password string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByPublicID(publicId string) (*models.User, error)
	GetAllPagination(filter, sort string, limit, offset int) ([]*models.User, int64, error)
	Update(user *models.User) error
	Delete(id uint) error
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

func (s *userService) Login(email, password string) (*models.User, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, errors.New("Invalid Credentials")
	}

	t, err := utils.CheckPasswordHash(user.Password, password)
	if err != nil {
		return nil, err
	}

	if !t {
		return nil, errors.New("Invalid Credentials")
	}

	return user, nil
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	return s.userRepository.FindByID(id)
}

func (s *userService) GetByPublicID(publicId string) (*models.User, error) {
	return s.userRepository.FindByPublicID(publicId)
}

func (s *userService) GetAllPagination(filter string, sort string, limit int, offset int) ([]*models.User, int64, error) {
	return s.userRepository.FindAllPagination(filter, sort, limit, offset)
}

func (s *userService) Update(user *models.User) error {
	return s.userRepository.Update(user)
}

func (s *userService) Delete(id uint) error {
	return s.userRepository.Delete(id)
}
