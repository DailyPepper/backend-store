package storage

import (
	"backend-store/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sqlx.DB
}

type PostgresTx struct {
	tx *sqlx.Tx
}

// NewPostgresStorage создает новое подключение к PostgreSQL
func NewPostgresStorage(databaseURL string) (*PostgresStorage, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

// Init инициализирует таблицы в базе данных
func (p *PostgresStorage) Init() error {
	// Создание таблицы продуктов
	_, err := p.db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			price DECIMAL(10,2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create products table: %w", err)
	}

	// Создание таблицы заказов
	_, err = p.db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			customer_name VARCHAR(255) NOT NULL,
			total_amount DECIMAL(10,2) NOT NULL,
			status VARCHAR(50) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create orders table: %w", err)
	}

	return nil
}

// Close закрывает подключение к базе данных
func (p *PostgresStorage) Close() error {
	return p.db.Close()
}

func (p *PostgresStorage) BeginTx(ctx context.Context) (StorageTx, error) {
	tx, err := p.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return &PostgresTx{tx: tx}, nil
}

// Добавляем недостающий метод BeginTx для PostgresTx
func (pt *PostgresTx) BeginTx(ctx context.Context) (StorageTx, error) {
	return pt, nil
}

// Реализация методов Storage для PostgresTx
func (pt *PostgresTx) CreateProduct(product *models.Product) error {
	query := `
	INSERT INTO products (name, description, price) 
	VALUES ($1, $2, $3) 
	RETURNING id, created_at`

	return pt.tx.QueryRow(
		query,
		product.Name,
		product.Description,
		product.Price,
	).Scan(&product.ID, &product.CreatedAt)
}

func (pt *PostgresTx) GetAllProduct() ([]*models.Product, error) {
	query := `SELECT id, name, description, price, created_at FROM products`
	var products []*models.Product
	err := pt.tx.Select(&products, query)
	return products, err
}

func (pt *PostgresTx) GetProductByID(id int) (*models.Product, error) {
	query := `SELECT id, name, description, price, created_at FROM products WHERE id = $1`
	var product models.Product
	err := pt.tx.Get(&product, query, id)
	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	return &product, err
}

func (pt *PostgresTx) UpdateProduct(product *models.Product) error {
	query := `UPDATE products SET name = $1, description = $2, price = $3 WHERE id = $4`
	result, err := pt.tx.Exec(query, product.Name, product.Description, product.Price, product.ID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (pt *PostgresTx) DeleteProduct(id int) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := pt.tx.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (pt *PostgresTx) CreateOrder(order *models.Order) error {
	query := `
	INSERT INTO orders (customer_name, total_amount, status) 
	VALUES ($1, $2, $3) 
	RETURNING id, created_at`

	return pt.tx.QueryRow(
		query,
		order.CustomerName,
		order.TotalAmount,
		order.Status,
	).Scan(&order.ID, &order.CreatedAt)
}

func (pt *PostgresTx) GetAllOrders() ([]*models.Order, error) {
	query := `SELECT id, customer_name, total_amount, status, created_at FROM orders`
	var orders []*models.Order
	err := pt.tx.Select(&orders, query)
	return orders, err
}

func (pt *PostgresTx) GetOrderByID(id int) (*models.Order, error) {
	query := `SELECT id, customer_name, total_amount, status, created_at FROM orders WHERE id = $1`
	var order models.Order
	err := pt.tx.Get(&order, query, id)
	if err == sql.ErrNoRows {
		return nil, errors.New("order not found")
	}
	return &order, err
}

func (pt *PostgresTx) UpdateOrder(order *models.Order) error {
	query := `UPDATE orders SET customer_name = $1, total_amount = $2, status = $3 WHERE id = $4`
	result, err := pt.tx.Exec(query, order.CustomerName, order.TotalAmount, order.Status, order.ID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (pt *PostgresTx) DeleteOrder(id int) error {
	query := `DELETE FROM orders WHERE id = $1`
	result, err := pt.tx.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (pt *PostgresTx) Commit() error {
	return pt.tx.Commit()
}

func (pt *PostgresTx) Rollback() error {
	return pt.tx.Rollback()
}

func (p *PostgresStorage) CreateProduct(product *models.Product) error {
	query := `
	INSERT INTO products (name, description, price) 
	VALUES ($1, $2, $3) 
	RETURNING id, created_at`

	return p.db.QueryRow(
		query,
		product.Name,
		product.Description,
		product.Price,
	).Scan(&product.ID, &product.CreatedAt)
}

func (p *PostgresStorage) GetAllProduct() ([]*models.Product, error) {
	query := `SELECT id, name, description, price, created_at FROM products`
	var products []*models.Product
	err := p.db.Select(&products, query)
	return products, err
}

func (p *PostgresStorage) GetProductByID(id int) (*models.Product, error) {
	query := `SELECT id, name, description, price, created_at FROM products WHERE id = $1`
	var product models.Product
	err := p.db.Get(&product, query, id)
	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	return &product, err
}

func (p *PostgresStorage) UpdateProduct(product *models.Product) error {
	query := `UPDATE products SET name = $1, description = $2, price = $3 WHERE id = $4`
	result, err := p.db.Exec(query, product.Name, product.Description, product.Price, product.ID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (p *PostgresStorage) DeleteProduct(id int) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := p.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (p *PostgresStorage) CreateOrder(order *models.Order) error {
	query := `
	INSERT INTO orders (customer_name, total_amount, status) 
	VALUES ($1, $2, $3) 
	RETURNING id, created_at`

	return p.db.QueryRow(
		query,
		order.CustomerName,
		order.TotalAmount,
		order.Status,
	).Scan(&order.ID, &order.CreatedAt)
}

func (p *PostgresStorage) GetAllOrders() ([]*models.Order, error) {
	query := `SELECT id, customer_name, total_amount, status, created_at FROM orders`
	var orders []*models.Order
	err := p.db.Select(&orders, query)
	return orders, err
}

func (p *PostgresStorage) GetOrderByID(id int) (*models.Order, error) {
	query := `SELECT id, customer_name, total_amount, status, created_at FROM orders WHERE id = $1`
	var order models.Order
	err := p.db.Get(&order, query, id)
	if err == sql.ErrNoRows {
		return nil, errors.New("order not found")
	}
	return &order, err
}

func (p *PostgresStorage) UpdateOrder(order *models.Order) error {
	query := `UPDATE orders SET customer_name = $1, total_amount = $2, status = $3 WHERE id = $4`
	result, err := p.db.Exec(query, order.CustomerName, order.TotalAmount, order.Status, order.ID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (p *PostgresStorage) DeleteOrder(id int) error {
	query := `DELETE FROM orders WHERE id = $1`
	result, err := p.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}
