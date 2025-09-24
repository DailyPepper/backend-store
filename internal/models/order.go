package models

import "time"

type Order struct {
	ID           int       `json:"id" db:"id"`
	CustomerName string    `json:"customer_name" db:"customer_name"`
	TotalAmount  float64   `json:"total_amount" db:"total_amount"`
	Status       string    `json:"status" db:"status"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type OrderItem struct {
	OrderId   int     `json:"order_id"`
	ProductId int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
