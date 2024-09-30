package services

import (
	"errors"
	"payment-app/models"
	"payment-app/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
    RegisterUser(user *models.User) error
    LoginUser(email, password string) (*models.User, error)
}

type userService struct {
    userRepository repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
    return &userService{
        userRepository: userRepo,
    }
}

func (s *userService) RegisterUser(user *models.User) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.Password = string(hashedPassword)

    return s.userRepository.Create(user)
}

func (s *userService) LoginUser(email, password string) (*models.User, error) {
    user, err := s.userRepository.GetByEmail(email)
    if err != nil {
        return nil, errors.New("user not found")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, errors.New("invalid credentials")
    }

    return user, nil
}
