package storage

import (
	"backend-store/internal/models"
	"context"
	"errors"
	"sync"
)

type MemoryStorage struct {
	products     map[int]*models.Product
	orders       map[int]*models.Order
	productIDSeq int
	orderIDSeq   int
	mu           sync.RWMutex
}

// MemoryTx - транзакция для in-memory хранилища
type MemoryTx struct {
	storage *MemoryStorage
	mu      *sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		products: make(map[int]*models.Product),
		orders:   make(map[int]*models.Order),
	}
}

func (m *MemoryStorage) BeginTx(ctx context.Context) (StorageTx, error) {
	m.mu.Lock()

	return &MemoryTx{
		storage: m,
		mu:      &m.mu,
	}, nil
}

func (mt *MemoryTx) BeginTx(ctx context.Context) (StorageTx, error) {
	return mt, nil
}

func (mt *MemoryTx) Commit() error {
	mt.mu.Unlock()
	return nil
}

func (mt *MemoryTx) Rollback() error {
	mt.mu.Unlock()
	return nil
}

// Реализация методов Storage для MemoryTx
func (mt *MemoryTx) CreateProduct(product *models.Product) error {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	mt.storage.productIDSeq++
	product.ID = mt.storage.productIDSeq
	mt.storage.products[product.ID] = product
	return nil
}

func (mt *MemoryTx) GetAllProduct() ([]*models.Product, error) {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	products := make([]*models.Product, 0, len(mt.storage.products))
	for _, p := range mt.storage.products {
		products = append(products, p)
	}
	return products, nil
}

func (mt *MemoryTx) GetProductByID(id int) (*models.Product, error) {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	product, exists := mt.storage.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (mt *MemoryTx) UpdateProduct(product *models.Product) error {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	if _, exists := mt.storage.products[product.ID]; !exists {
		return errors.New("product not found")
	}
	mt.storage.products[product.ID] = product
	return nil
}

func (mt *MemoryTx) DeleteProduct(id int) error {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	if _, exists := mt.storage.products[id]; !exists {
		return errors.New("product not found")
	}
	delete(mt.storage.products, id)
	return nil
}

func (mt *MemoryTx) CreateOrder(order *models.Order) error {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	mt.storage.orderIDSeq++
	order.ID = mt.storage.orderIDSeq
	mt.storage.orders[order.ID] = order
	return nil
}

func (mt *MemoryTx) GetAllOrders() ([]*models.Order, error) {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	orders := make([]*models.Order, 0, len(mt.storage.orders))
	for _, o := range mt.storage.orders {
		orders = append(orders, o)
	}
	return orders, nil
}

func (mt *MemoryTx) GetOrderByID(id int) (*models.Order, error) {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	order, exists := mt.storage.orders[id]
	if !exists {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (mt *MemoryTx) UpdateOrder(order *models.Order) error {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	if _, exists := mt.storage.orders[order.ID]; !exists {
		return errors.New("order not found")
	}
	mt.storage.orders[order.ID] = order
	return nil
}

func (mt *MemoryTx) DeleteOrder(id int) error {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	if _, exists := mt.storage.orders[id]; !exists {
		return errors.New("order not found")
	}
	delete(mt.storage.orders, id)
	return nil
}

// Реализация методов Storage для MemoryStorage (без транзакций)
func (m *MemoryStorage) CreateProduct(product *models.Product) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.productIDSeq++
	product.ID = m.productIDSeq
	m.products[product.ID] = product
	return nil
}

func (m *MemoryStorage) GetAllProduct() ([]*models.Product, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	products := make([]*models.Product, 0, len(m.products))
	for _, p := range m.products {
		products = append(products, p)
	}
	return products, nil
}

func (m *MemoryStorage) GetProductByID(id int) (*models.Product, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	product, exists := m.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (m *MemoryStorage) UpdateProduct(product *models.Product) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.products[product.ID]; !exists {
		return errors.New("product not found")
	}
	m.products[product.ID] = product
	return nil
}

func (m *MemoryStorage) DeleteProduct(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.products[id]; !exists {
		return errors.New("product not found")
	}
	delete(m.products, id)
	return nil
}

func (m *MemoryStorage) CreateOrder(order *models.Order) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.orderIDSeq++
	order.ID = m.orderIDSeq
	m.orders[order.ID] = order
	return nil
}

func (m *MemoryStorage) GetAllOrders() ([]*models.Order, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	orders := make([]*models.Order, 0, len(m.orders))
	for _, o := range m.orders {
		orders = append(orders, o)
	}
	return orders, nil
}

func (m *MemoryStorage) GetOrderByID(id int) (*models.Order, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	order, exists := m.orders[id]
	if !exists {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (m *MemoryStorage) UpdateOrder(order *models.Order) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.orders[order.ID]; !exists {
		return errors.New("order not found")
	}
	m.orders[order.ID] = order
	return nil
}

func (m *MemoryStorage) DeleteOrder(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.orders[id]; !exists {
		return errors.New("order not found")
	}
	delete(m.orders, id)
	return nil
}
