package repository

import (
	"backend-store/internal/models"
	"backend-store/internal/storage"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id string) error
}

type userRepository struct {
	storage storage.Storage
}

func NewUserRepository(storage storage.Storage) UserRepository {
	return &userRepository{storage: storage}
}

func (r *userRepository) CreateUser(user *models.User) error {
	// TODO: реализовать
	return nil
}

func (r *userRepository) GetUserByID(id string) (*models.User, error) {
	// TODO: реализовать
	return &models.User{}, nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	// TODO: реализовать
	return &models.User{}, nil
}

func (r *userRepository) UpdateUser(user *models.User) error {
	// TODO: реализовать
	return nil
}

func (r *userRepository) DeleteUser(id string) error {
	// TODO: реализовать
	return nil
}
