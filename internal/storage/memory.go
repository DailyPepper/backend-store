package storage

import (
	"backend-store/internal/models"
	"errors"
	"sync"
)

type MemoryStorage struct {
	orders    []*models.Order
	products  []*models.Product
	orderID   int
	productID int
	mu        sync.RWMutex
}

// CreateOrder implements Storage.
func (m *MemoryStorage) CreateOrder(order *models.Order) error {
	panic("unimplemented")
}

// GetAllOrders implements Storage.
func (m *MemoryStorage) GetAllOrders() ([]*models.Order, error) {
	panic("unimplemented")
}

// GetOrderByID implements Storage.
func (m *MemoryStorage) GetOrderByID(id int) (*models.Order, error) {
	panic("unimplemented")
}

// UpdateOrder implements Storage.
func (m *MemoryStorage) UpdateOrder(order *models.Order) error {
	panic("unimplemented")
}

// UpdateProduct implements Storage.
func (m *MemoryStorage) UpdateProduct(product *models.Product) error {
	panic("unimplemented")
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		orders:    make([]*models.Order, 0),
		products:  make([]*models.Product, 0),
		orderID:   1,
		productID: 1,
	}
}

func (m *MemoryStorage) CreateProduct(product *models.Product) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	product.ID = m.productID
	m.productID++
	m.products = append(m.products, product)

	return nil
}

func (m *MemoryStorage) GetAllProduct() ([]*models.Product, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.products, nil
}

func (m *MemoryStorage) GetProductByID(id int) (*models.Product, error) {
	m.mu.RLock() // Блокировка для чтения
	defer m.mu.RUnlock()

	for _, product := range m.products {
		if product.ID == id {
			return product, nil
		}
	}

	return nil, errors.New("product not found")
}

func (m *MemoryStorage) DeleteProduct(id int) error {
	m.mu.Lock() // Блокировка для записи
	defer m.mu.Unlock()

	for i, product := range m.products {
		if product.ID == id {
			// Удаляем элемент из слайса
			m.products = append(m.products[:i], m.products[i+1:]...)
			return nil
		}
	}

	return errors.New("product not found")
}

func (m *MemoryStorage) DeleteOrder(id int) error {
	m.mu.Lock() // Блокировка для записи
	defer m.mu.Unlock()

	for i, product := range m.products {
		if product.ID == id {
			// Удаляем элемент из слайса
			m.products = append(m.products[:i], m.products[i+1:]...)
			return nil
		}
	}

	return errors.New("product not found")
}
