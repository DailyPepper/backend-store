package models

import "time"

type Order struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
	// Product
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrderItem struct {
	OrderId   int     `json:"order_id"`
	ProductId int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
