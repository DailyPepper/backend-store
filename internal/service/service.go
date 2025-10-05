package service

import (
	"backend-store/internal/models"
	"backend-store/internal/storage"
	"context"
	"errors"
	"fmt"
	"time"
)

type orderService struct {
	storage storage.Storage
}

func NewOrderService(storage storage.Storage) OrderService {
	return &orderService{storage: storage}
}

func (s *orderService) CreateOrder(ctx context.Context, order *models.Order) error {
	if err := order.Validate(); err != nil {
		return err
	}

	if order.Status == "" {
		order.Status = "pending"
	}

	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now

	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	for _, item := range order.Products {
		product, err := tx.GetProductByID(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("product %d not found: %w", item.ProductID, err)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("insufficient quantity for product %d: available %d, requested %d",
				item.ProductID, product.Quantity, item.Quantity)
		}
	}

	if err := tx.CreateOrder(ctx, order); err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return tx.Commit()
}

func (s *orderService) GetAllOrders(ctx context.Context) ([]*models.Order, error) {
	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	orders, err := tx.GetAllOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *orderService) GetOrderByID(ctx context.Context, id int) (*models.Order, error) {
	if id <= 0 {
		return nil, errors.New("invalid order ID")
	}

	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	order, err := tx.GetOrderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) UpdateOrder(ctx context.Context, order *models.Order) error {
	if err := order.Validate(); err != nil {
		return err
	}

	order.UpdatedAt = time.Now()

	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	existingOrder, err := tx.GetOrderByID(ctx, order.ID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	order.CreatedAt = existingOrder.CreatedAt

	if err := tx.UpdateOrder(ctx, order); err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return tx.Commit()
}

func (s *orderService) DeleteOrder(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid order ID")
	}

	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.GetOrderByID(ctx, id)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	if err := tx.DeleteOrder(ctx, id); err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	return tx.Commit()
}

type productService struct {
	storage storage.Storage
}

func NewProductService(storage storage.Storage) ProductService {
	return &productService{storage: storage}
}

func (s *productService) CreateProduct(ctx context.Context, product *models.Product) error {
	if err := product.Validate(); err != nil {
		return err
	}

	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := tx.CreateProduct(ctx, product); err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return tx.Commit()
}

func (s *productService) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	products, err := tx.GetAllProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return products, nil
}

func (s *productService) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	if id <= 0 {
		return nil, errors.New("invalid product ID")
	}

	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	product, err := tx.GetProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) UpdateProduct(ctx context.Context, product *models.Product) error {
	if err := product.Validate(); err != nil {
		return err
	}

	product.UpdatedAt = time.Now()

	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	existingProduct, err := tx.GetProductByID(ctx, product.ID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	product.CreatedAt = existingProduct.CreatedAt

	if err := tx.UpdateProduct(ctx, product); err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return tx.Commit()
}

func (s *productService) DeleteProduct(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid product ID")
	}

	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.GetProductByID(ctx, id)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	if err := tx.DeleteProduct(ctx, id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return tx.Commit()
}
