package storage

import (
	"backend-store/internal/models"
	"sync"
)

type MemoryStorage struct {
	mu               sync.RWMutex
	orders           map[int]*models.Order
	products         map[int]*models.Product
	orderItems       map[int][]*models.OrderItem // orderID -> []OrderItem
	orderCounter     int
	productCounter   int
	orderItemCounter int
}

// GetAllOrders implements Storage.
func (m *MemoryStorage) GetAllOrders() ([]*models.Order, error) {
	panic("unimplemented")
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		orders:           make(map[int]*models.Order),
		products:         make(map[int]*models.Product),
		orderItems:       make(map[int][]*models.OrderItem),
		orderCounter:     1,
		productCounter:   1,
		orderItemCounter: 1,
	}
}
