package storage

import (
	"backend-store/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sqlx.DB
}

type PostgresTx struct {
	tx *sqlx.Tx
}

func NewPostgresStorage(databaseURL string) (*PostgresStorage, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

func (p *PostgresStorage) Init() error {
	_, err := p.db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			price INTEGER NOT NULL,
			quantity INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create products table: %w", err)
	}

	_, err = p.db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			status VARCHAR(50) NOT NULL,
			total INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create orders table: %w", err)
	}

	_, err = p.db.Exec(`
		CREATE TABLE IF NOT EXISTS order_items (
			id SERIAL PRIMARY KEY,
			order_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			price INTEGER NOT NULL,
			FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
			FOREIGN KEY (product_id) REFERENCES products(id)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create order_items table: %w", err)
	}

	return nil
}

func (p *PostgresStorage) SetMaxOpenConns(conns int) {
	p.db.SetMaxOpenConns(conns)
}

func (p *PostgresStorage) SetMaxIdleConns(conns int) {
	p.db.SetMaxIdleConns(conns)
}

func (p *PostgresStorage) SetConnMaxLifetime(lifetime time.Duration) {
	p.db.SetConnMaxLifetime(lifetime)
}

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

func (p *PostgresStorage) CreateProduct(ctx context.Context, product *models.Product) error {
	query := `
	INSERT INTO products (name, description, price, quantity) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id, created_at, updated_at`

	return p.db.QueryRowContext(ctx,
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Quantity,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
}

func (p *PostgresStorage) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	query := `SELECT id, name, description, price, quantity, created_at, updated_at FROM products`
	var products []*models.Product
	err := p.db.SelectContext(ctx, &products, query)
	return products, err
}

func (p *PostgresStorage) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	query := `SELECT id, name, description, price, quantity, created_at, updated_at FROM products WHERE id = $1`
	var product models.Product
	err := p.db.GetContext(ctx, &product, query, id)
	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	return &product, err
}

func (p *PostgresStorage) UpdateProduct(ctx context.Context, product *models.Product) error {
	query := `
		UPDATE products 
		SET name = $1, description = $2, price = $3, quantity = $4, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $5
	`
	result, err := p.db.ExecContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.Quantity,
		product.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (p *PostgresStorage) DeleteProduct(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (p *PostgresStorage) CreateOrder(ctx context.Context, order *models.Order) error {
	total := 0
	for _, item := range order.Products {
		var productPrice int
		err := p.db.GetContext(ctx, &productPrice, "SELECT price FROM products WHERE id = $1", item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product price: %w", err)
		}
		total += productPrice * item.Quantity
		item.Price = productPrice
	}
	order.Total = total

	orderQuery := `
		INSERT INTO orders (user_id, status, total) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at`

	err := p.db.QueryRowContext(ctx,
		orderQuery,
		order.ID,
		order.Status,
		order.Total,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return err
	}

	itemQuery := `
		INSERT INTO order_items (order_id, product_id, quantity, price) 
		VALUES ($1, $2, $3, $4)`

	for _, item := range order.Products {
		_, err = p.db.ExecContext(ctx, itemQuery, order.ID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return fmt.Errorf("failed to create order item: %w", err)
		}
	}

	return nil
}

func (p *PostgresStorage) GetAllOrders(ctx context.Context) ([]*models.Order, error) {
	query := `SELECT id, user_id, status, total, created_at, updated_at FROM orders`
	var orders []*models.Order
	err := p.db.SelectContext(ctx, &orders, query)
	if err != nil {
		return nil, err
	}

	for i := range orders {
		itemsQuery := `
			SELECT id, order_id, product_id, quantity, price 
			FROM order_items 
			WHERE order_id = $1`

		var items []models.OrderItem
		err := p.db.SelectContext(ctx, &items, itemsQuery, orders[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get order items: %w", err)
		}
		orders[i].Products = items
	}

	return orders, nil
}

func (p *PostgresStorage) GetOrderByID(ctx context.Context, id int) (*models.Order, error) {
	query := `SELECT id, user_id, status, total, created_at, updated_at FROM orders WHERE id = $1`
	var order models.Order
	err := p.db.GetContext(ctx, &order, query, id)
	if err == sql.ErrNoRows {
		return nil, errors.New("order not found")
	}
	if err != nil {
		return nil, err
	}

	itemsQuery := `
		SELECT id, order_id, product_id, quantity, price 
		FROM order_items 
		WHERE order_id = $1`

	var items []models.OrderItem
	err = p.db.SelectContext(ctx, &items, itemsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	order.Products = items

	return &order, nil
}

func (p *PostgresStorage) UpdateOrder(ctx context.Context, order *models.Order) error {
	total := 0
	for _, item := range order.Products {
		var productPrice int
		err := p.db.GetContext(ctx, &productPrice, "SELECT price FROM products WHERE id = $1", item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product price: %w", err)
		}
		total += productPrice * item.Quantity
		item.Price = productPrice
	}
	order.Total = total

	query := `
		UPDATE orders 
		SET user_id = $1, status = $2, total = $3, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $4`

	result, err := p.db.ExecContext(ctx, query, order.ID, order.Status, order.Total)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("order not found")
	}

	_, err = p.db.ExecContext(ctx, "DELETE FROM order_items WHERE order_id = $1", order.ID)
	if err != nil {
		return fmt.Errorf("failed to delete old order items: %w", err)
	}

	itemQuery := `
		INSERT INTO order_items (order_id, product_id, quantity, price) 
		VALUES ($1, $2, $3, $4)`

	for _, item := range order.Products {
		_, err = p.db.ExecContext(ctx, itemQuery, order.ID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return fmt.Errorf("failed to create order item: %w", err)
		}
	}

	return nil
}

func (p *PostgresStorage) DeleteOrder(ctx context.Context, id int) error {
	query := `DELETE FROM orders WHERE id = $1`
	result, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (pt *PostgresTx) BeginTx(ctx context.Context) (StorageTx, error) {
	return pt, nil
}

func (pt *PostgresTx) CreateProduct(ctx context.Context, product *models.Product) error {
	query := `
	INSERT INTO products (name, description, price, quantity) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id, created_at, updated_at`

	return pt.tx.QueryRowContext(ctx,
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Quantity,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
}

func (pt *PostgresTx) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	query := `SELECT id, name, description, price, quantity, created_at, updated_at FROM products`
	var products []*models.Product
	err := pt.tx.SelectContext(ctx, &products, query)
	return products, err
}

func (pt *PostgresTx) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	query := `SELECT id, name, description, price, quantity, created_at, updated_at FROM products WHERE id = $1`
	var product models.Product
	err := pt.tx.GetContext(ctx, &product, query, id)
	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	return &product, err
}

func (pt *PostgresTx) UpdateProduct(ctx context.Context, product *models.Product) error {
	query := `
		UPDATE products 
		SET name = $1, description = $2, price = $3, quantity = $4, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $5
	`
	result, err := pt.tx.ExecContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.Quantity,
		product.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (pt *PostgresTx) DeleteProduct(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := pt.tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (pt *PostgresTx) CreateOrder(ctx context.Context, order *models.Order) error {
	total := 0
	for _, item := range order.Products {
		var productPrice int
		err := pt.tx.GetContext(ctx, &productPrice, "SELECT price FROM products WHERE id = $1", item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product price: %w", err)
		}
		total += productPrice * item.Quantity
		item.Price = productPrice
	}
	order.Total = total

	orderQuery := `
		INSERT INTO orders (user_id, status, total) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at`

	err := pt.tx.QueryRowContext(ctx,
		orderQuery,
		order.ID,
		order.Status,
		order.Total,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return err
	}

	itemQuery := `
		INSERT INTO order_items (order_id, product_id, quantity, price) 
		VALUES ($1, $2, $3, $4)`

	for _, item := range order.Products {
		_, err = pt.tx.ExecContext(ctx, itemQuery, order.ID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return fmt.Errorf("failed to create order item: %w", err)
		}
	}

	return nil
}

func (pt *PostgresTx) GetAllOrders(ctx context.Context) ([]*models.Order, error) {
	query := `SELECT id, user_id, status, total, created_at, updated_at FROM orders`
	var orders []*models.Order
	err := pt.tx.SelectContext(ctx, &orders, query)
	if err != nil {
		return nil, err
	}

	for i := range orders {
		itemsQuery := `
			SELECT id, order_id, product_id, quantity, price 
			FROM order_items 
			WHERE order_id = $1`

		var items []models.OrderItem
		err := pt.tx.SelectContext(ctx, &items, itemsQuery, orders[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get order items: %w", err)
		}
		orders[i].Products = items
	}

	return orders, nil
}

func (pt *PostgresTx) GetOrderByID(ctx context.Context, id int) (*models.Order, error) {
	query := `SELECT id, user_id, status, total, created_at, updated_at FROM orders WHERE id = $1`
	var order models.Order
	err := pt.tx.GetContext(ctx, &order, query, id)
	if err == sql.ErrNoRows {
		return nil, errors.New("order not found")
	}
	if err != nil {
		return nil, err
	}

	itemsQuery := `
		SELECT id, order_id, product_id, quantity, price 
		FROM order_items 
		WHERE order_id = $1`

	var items []models.OrderItem
	err = pt.tx.SelectContext(ctx, &items, itemsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	order.Products = items

	return &order, nil
}

func (pt *PostgresTx) UpdateOrder(ctx context.Context, order *models.Order) error {
	total := 0
	for _, item := range order.Products {
		var productPrice int
		err := pt.tx.GetContext(ctx, &productPrice, "SELECT price FROM products WHERE id = $1", item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product price: %w", err)
		}
		total += productPrice * item.Quantity
		item.Price = productPrice
	}
	order.Total = total

	query := `
		UPDATE orders 
		SET user_id = $1, status = $2, total = $3, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $4`

	result, err := pt.tx.ExecContext(ctx, query, order.ID, order.Status, order.Total)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("order not found")
	}

	_, err = pt.tx.ExecContext(ctx, "DELETE FROM order_items WHERE order_id = $1", order.ID)
	if err != nil {
		return fmt.Errorf("failed to delete old order items: %w", err)
	}

	itemQuery := `
		INSERT INTO order_items (order_id, product_id, quantity, price) 
		VALUES ($1, $2, $3, $4)`

	for _, item := range order.Products {
		_, err = pt.tx.ExecContext(ctx, itemQuery, order.ID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return fmt.Errorf("failed to create order item: %w", err)
		}
	}

	return nil
}

func (pt *PostgresTx) DeleteOrder(ctx context.Context, id int) error {
	query := `DELETE FROM orders WHERE id = $1`
	result, err := pt.tx.ExecContext(ctx, query, id)
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

func (pt *PostgresTx) Close() error {
	return nil
}

func (p *PostgresStorage) UpdateProductQuantity(ctx context.Context, id int, quantity int) error {
	query := `UPDATE products SET quantity = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := p.db.ExecContext(ctx, query, quantity, id)
	return err
}
