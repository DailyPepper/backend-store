package storage

import (
	"backend-store/internal/models"
	"context"
)

type Storage interface {
	CreateProduct(product *models.Product) error
	GetAllProduct() ([]*models.Product, error)
	GetProductByID(id int) (*models.Product, error)
	UpdateProduct(product *models.Product) error
	DeleteProduct(id int) error

	CreateOrder(order *models.Order) error
	GetAllOrders() ([]*models.Order, error)
	GetOrderByID(id int) (*models.Order, error)
	UpdateOrder(order *models.Order) error
	DeleteOrder(id int) error

	BeginTx(ctx context.Context) (StorageTx, error)
}

type StorageTx interface {
	Storage
	Commit() error
	Rollback() error
	BeginTx(ctx context.Context) (StorageTx, error)
}
