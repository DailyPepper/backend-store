package service

import (
	"backend-store/internal/models"
	"backend-store/internal/repository"
)

type UserService interface {
	GetUserProfile(userID string) (*models.User, error)
	UpdateUserProfile(user *models.User) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUserProfile(userID string) (*models.User, error) {
	return s.userRepo.GetUserByID(userID)
}

func (s *userService) UpdateUserProfile(user *models.User) error {
	return s.userRepo.UpdateUser(user)
}
