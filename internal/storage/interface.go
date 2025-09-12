package storage

import "backend-store/internal/models"

type Storage interface {

	// Order
	GetAllOrders() ([]*models.Order, error)
}
