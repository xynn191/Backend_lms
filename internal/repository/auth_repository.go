package repository

import (
	"github.com/xynn191/Backend_lms/internal/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FindUserByEmail(email string) (*models.User, error)
	FindUserByID(id string) (*models.User, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("email = ? AND is_active = ?", email, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindUserByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("id = ? AND is_active = ?", id, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
