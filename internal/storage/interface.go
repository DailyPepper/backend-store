package storage

import "backend-store/internal/models"

type Storage interface {
	// Order
	CreateOrder() ([]*models.Order, error)
	GetAllOrders() ([]*models.Order, error)
	UpdateOrder() ([]*models.OrderItem, error)
	DeleteOrder() ([]*models.OrderItem, error)

	// Product
	CreateProduct() ([]*models.Product, error)
	GetAllProduct() ([]*models.Product, error)
	UpdateProduct() ([]*models.Product, error)
	DeleteProduct() ([]*models.Product, error)
}
