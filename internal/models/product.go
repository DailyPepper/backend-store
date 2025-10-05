package models

import (
	"errors"
	"time"
)

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return errors.New("product name is required")
	}
	if len(p.Name) > 100 {
		return errors.New("product name is too long")
	}
	if p.Price <= 0 {
		return errors.New("product price must be positive")
	}
	if p.Quantity < 0 {
		return errors.New("product quantity cannot be negative")
	}
	return nil
}
