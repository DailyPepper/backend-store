package service

import (
	"backend-store/internal/models"
	"context"
)

type ProductService interface {
	CreateProduct(ctx context.Context, product *models.Product) error
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
	GetProductByID(ctx context.Context, id int) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id int) error
}

type OrderService interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetAllOrders(ctx context.Context) ([]*models.Order, error)
	GetOrderByID(ctx context.Context, id int) (*models.Order, error)
	UpdateOrder(ctx context.Context, order *models.Order) error
	DeleteOrder(ctx context.Context, id int) error
}
