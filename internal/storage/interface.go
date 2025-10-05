package storage

import (
	"backend-store/internal/models"
	"context"
)

type Storage interface {
	// Products
	CreateProduct(ctx context.Context, product *models.Product) error
	GetProductByID(ctx context.Context, id int) (*models.Product, error)
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id int) error

	// Orders
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrderByID(ctx context.Context, id int) (*models.Order, error)
	GetAllOrders(ctx context.Context) ([]*models.Order, error)
	UpdateOrder(ctx context.Context, order *models.Order) error
	DeleteOrder(ctx context.Context, id int) error

	// Transactions
	BeginTx(ctx context.Context) (StorageTx, error)

	// Закрытие соединения (опционально, только для некоторых реализаций)
	Close() error
}

// StorageTx интерфейс для транзакций
type StorageTx interface {
	Storage
	Commit() error
	Rollback() error
}
