package repositories

import (
	"payment-app/config"
	"payment-app/models"
)

type UserRepository interface {
    Create(user *models.User) error
    GetByEmail(email string) (*models.User, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
    return &userRepository{}
}

func (r *userRepository) Create(user *models.User) error {
    return config.DB.Create(user).Error
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
    var user models.User
    if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
