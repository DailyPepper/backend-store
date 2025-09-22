package storage

import "backend-store/internal/models"

type Storage interface {
	// Order
	CreateOrder(order *models.Order) error
	GetAllOrders() ([]*models.Order, error)
	GetOrderByID(id int) (*models.Order, error)
	UpdateOrder(order *models.Order) error
	DeleteOrder(id int) error

	// Product
	CreateProduct(product *models.Product) error
	GetAllProduct() ([]*models.Product, error)
	GetProductByID(id int) (*models.Product, error)
	UpdateProduct(product *models.Product) error
	DeleteProduct(id int) error
}
