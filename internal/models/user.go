package models

import (
	"time"
)

type User struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // скрываем при сериализации в JSON
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Phone        string    `json:"phone" db:"phone"`
	Address      string    `json:"address" db:"address"`
	Role         string    `json:"role" db:"role"` // user, admin, etc.
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserCreateRequest - запрос на создание пользователя
type UserCreateRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}

// UserUpdateRequest - запрос на обновление пользователя
type UserUpdateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}

// UserResponse - ответ с данными пользователя
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginRequest - запрос на вход
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse - ответ с токеном
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// ToResponse преобразует User в UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		Address:   u.Address,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// Validate проверяет валидность пользователя
func (u *User) Validate() error {
	if u.Email == "" {
		return ErrInvalidEmail
	}
	if u.PasswordHash == "" {
		return ErrInvalidPassword
	}
	if u.FirstName == "" {
		return ErrInvalidFirstName
	}
	if u.LastName == "" {
		return ErrInvalidLastName
	}
	return nil
}

// Ошибки валидации
var (
	ErrInvalidEmail     = NewValidationError("email is required")
	ErrInvalidPassword  = NewValidationError("password is required")
	ErrInvalidFirstName = NewValidationError("first name is required")
	ErrInvalidLastName  = NewValidationError("last name is required")
	ErrUserNotFound     = NewValidationError("user not found")
	ErrUserExists       = NewValidationError("user already exists")
)

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Message string
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

func (e *ValidationError) Error() string {
	return e.Message
}
