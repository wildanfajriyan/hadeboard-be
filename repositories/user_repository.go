package repositories

import (
	"context"
	"hadeboard-be/config"
	"hadeboard-be/internal/models"
	"strings"
	"time"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	FindByPublicID(publicId string) (*models.User, error)
	FindAllPagination(filter, sort string, limit, offset int) ([]*models.User, int64, error)
	Update(user *models.User) error
	Delete(id uint) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return config.DB.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := config.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := config.DB.WithContext(ctx).First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByPublicID(publicId string) (*models.User, error) {
	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := config.DB.WithContext(ctx).Where("public_id = ?", publicId).First(&user).Error
	return &user, err
}

func (r *userRepository) FindAllPagination(filter, sort string, limit, offset int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.DB.WithContext(ctx).Model(&models.User{})

	if filter != "" {
		filterPattern := "%" + filter + "%"
		db.Where("name Ilike ? OR email Ilike ?", filterPattern, filterPattern)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sort != "" {
		switch sort {
		case "-id":
			sort = "-internal_id"
		case "id":
			sort = "internal_id"
		default:
			sort = ""
		}

		if after, ok := strings.CutPrefix(sort, "-"); ok {
			sort = after + " DESC"
		} else {
			sort += " ASC"
		}

		db = db.Order(sort)
	}

	err := db.Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *userRepository) Update(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return config.DB.WithContext(ctx).Model(&models.User{}).Where("public_id = ?", user.PublicID).Updates(map[string]any{
		"name": user.Name,
	}).Error
}

func (r *userRepository) Delete(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return config.DB.WithContext(ctx).Delete(&models.User{}, id).Error
}
