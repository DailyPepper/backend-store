package models

import (
	"errors"
	"time"
)

type Order struct {
	ID        int         `json:"id"`
	UserID    int         `json:"user_id"`
	Products  []OrderItem `json:"products"`
	Status    string      `json:"status"`
	Total     int         `json:"total"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID        int `json:"id"`
	OrderID   int `json:"order_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
	Price     int `json:"price"`
}

func (o *Order) Validate() error {
	if o.UserID <= 0 {
		return errors.New("user ID is required")
	}
	if len(o.Products) == 0 {
		return errors.New("order must contain at least one product")
	}

	for _, item := range o.Products {
		if err := item.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (oi *OrderItem) Validate() error {
	if oi.ProductID <= 0 {
		return errors.New("product ID is required")
	}
	if oi.Quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if oi.Price < 0 {
		return errors.New("price cannot be negative")
	}
	return nil
}

func (o *Order) CalculateTotal() int {
	total := 0
	for _, item := range o.Products {
		total += item.Price * item.Quantity
	}
	return total
}
