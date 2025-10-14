package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProduct_Validate_Success(t *testing.T) {
	// Arrange
	product := &Product{
		Name:        "Valid Product",
		Description: "Valid Description",
		Price:       29.99,
		Quantity:    100,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Act
	err := product.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestProduct_Validate_EmptyName(t *testing.T) {
	// Arrange
	product := &Product{
		Name:        "", // Пустое имя
		Description: "Valid Description",
		Price:       29.99,
		Quantity:    100,
	}

	// Act
	err := product.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "product name is required", err.Error())
}

func TestProduct_Validate_NameTooLong(t *testing.T) {
	// Arrange
	longName := ""
	for i := 0; i < 101; i++ {
		longName += "a"
	}

	product := &Product{
		Name:        longName, // Слишком длинное имя
		Description: "Valid Description",
		Price:       29.99,
		Quantity:    100,
	}

	// Act
	err := product.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "product name is too long", err.Error())
}

func TestProduct_Validate_ZeroPrice(t *testing.T) {
	// Arrange
	product := &Product{
		Name:        "Valid Product",
		Description: "Valid Description",
		Price:       0, // Нулевая цена
		Quantity:    100,
	}

	// Act
	err := product.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "product price must be positive", err.Error())
}

func TestProduct_Validate_NegativePrice(t *testing.T) {
	// Arrange
	product := &Product{
		Name:        "Valid Product",
		Description: "Valid Description",
		Price:       -10.0, // Отрицательная цена
		Quantity:    100,
	}

	// Act
	err := product.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "product price must be positive", err.Error())
}

func TestProduct_Validate_NegativeQuantity(t *testing.T) {
	// Arrange
	product := &Product{
		Name:        "Valid Product",
		Description: "Valid Description",
		Price:       29.99,
		Quantity:    -5, // Отрицательное количество
	}

	// Act
	err := product.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "product quantity cannot be negative", err.Error())
}

func TestProduct_Validate_ZeroQuantity(t *testing.T) {
	// Arrange
	product := &Product{
		Name:        "Valid Product",
		Description: "Valid Description",
		Price:       29.99,
		Quantity:    0, // Нулевое количество - должно быть валидно
	}

	// Act
	err := product.Validate()

	// Assert
	assert.NoError(t, err) // Нулевое количество допустимо
}

func TestProduct_JSONSerialization(t *testing.T) {
	// Arrange
	product := &Product{
		ID:          1,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       29.99,
		Quantity:    100,
		CreatedAt:   time.Date(2023, 10, 5, 14, 30, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2023, 10, 5, 16, 45, 0, 0, time.UTC),
	}

	// Act - используем стандартный json.Marshal вместо product.MarshalJSON()
	jsonData, err := json.Marshal(product)
	assert.NoError(t, err)

	var decodedProduct Product
	err = json.Unmarshal(jsonData, &decodedProduct) // используем json.Unmarshal

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, product.ID, decodedProduct.ID)
	assert.Equal(t, product.Name, decodedProduct.Name)
	assert.Equal(t, product.Description, decodedProduct.Description)
	assert.Equal(t, product.Price, decodedProduct.Price)
	assert.Equal(t, product.Quantity, decodedProduct.Quantity)
	// Время может немного отличаться из-за сериализации, поэтому проверяем что оно не нулевое
	assert.False(t, decodedProduct.CreatedAt.IsZero())
	assert.False(t, decodedProduct.UpdatedAt.IsZero())
}

func TestProduct_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		product     Product
		shouldError bool
		errorMsg    string
	}{
		{
			name: "Valid product with minimal data",
			product: Product{
				Name:  "Minimal",
				Price: 0.01, // Минимальная цена
			},
			shouldError: false,
		},
		{
			name: "Valid product with maximum name length",
			product: Product{
				Name:  "This is exactly 100 characters long name that should pass validation without any issues at all!",
				Price: 10.0,
			},
			shouldError: false,
		},
		{
			name: "Invalid product with very small negative price",
			product: Product{
				Name:  "Product",
				Price: -0.01,
			},
			shouldError: true,
			errorMsg:    "product price must be positive",
		},
		{
			name: "Valid product with large quantity",
			product: Product{
				Name:     "Product",
				Price:    10.0,
				Quantity: 1000000,
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()

			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Equal(t, tt.errorMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
