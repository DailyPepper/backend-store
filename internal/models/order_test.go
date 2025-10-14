package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrder_JSONSerialization(t *testing.T) {
	// Arrange
	order := &Order{
		ID:     1,
		UserID: 1,
		Total:  9999,
		Status: "completed",
		Products: []OrderItem{
			{
				ID:        1,
				OrderID:   1,
				ProductID: 1,
				Quantity:  2,
				Price:     5000,
			},
		},
		CreatedAt: time.Date(2023, 10, 5, 14, 30, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 10, 5, 16, 45, 0, 0, time.UTC),
	}

	// Act - используем стандартный json.Marshal
	jsonData, err := json.Marshal(order)
	assert.NoError(t, err)

	var decodedOrder Order
	err = json.Unmarshal(jsonData, &decodedOrder) // используем json.Unmarshal

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, order.ID, decodedOrder.ID)
	assert.Equal(t, order.UserID, decodedOrder.UserID)
	assert.Equal(t, order.Total, decodedOrder.Total)
	assert.Equal(t, order.Status, decodedOrder.Status)
	assert.Len(t, decodedOrder.Products, 1)
	assert.Equal(t, order.Products[0].ProductID, decodedOrder.Products[0].ProductID)
	assert.Equal(t, order.Products[0].Quantity, decodedOrder.Products[0].Quantity)
	assert.Equal(t, order.Products[0].Price, decodedOrder.Products[0].Price)
}

func TestOrderItem_JSONSerialization(t *testing.T) {
	// Arrange
	orderItem := &OrderItem{
		ID:        1,
		OrderID:   1,
		ProductID: 1,
		Quantity:  2,
		Price:     5000,
	}

	// Act - используем стандартный json.Marshal
	jsonData, err := json.Marshal(orderItem)
	assert.NoError(t, err)

	var decodedOrderItem OrderItem
	err = json.Unmarshal(jsonData, &decodedOrderItem) // используем json.Unmarshal

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, orderItem.ID, decodedOrderItem.ID)
	assert.Equal(t, orderItem.OrderID, decodedOrderItem.OrderID)
	assert.Equal(t, orderItem.ProductID, decodedOrderItem.ProductID)
	assert.Equal(t, orderItem.Quantity, decodedOrderItem.Quantity)
	assert.Equal(t, orderItem.Price, decodedOrderItem.Price)
}

func TestOrder_EmptyProducts(t *testing.T) {
	// Arrange
	order := &Order{
		ID:        1,
		UserID:    1,
		Total:     0,
		Status:    "pending",
		Products:  []OrderItem{}, // Пустой список продуктов
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	jsonData, err := json.Marshal(order)
	assert.NoError(t, err)

	var decodedOrder Order
	err = json.Unmarshal(jsonData, &decodedOrder)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, order.ID, decodedOrder.ID)
	assert.Equal(t, order.UserID, decodedOrder.UserID)
	assert.Equal(t, order.Total, decodedOrder.Total)
	assert.Equal(t, order.Status, decodedOrder.Status)
	assert.Empty(t, decodedOrder.Products) // Должен быть пустой массив
}

func TestOrder_NilProducts(t *testing.T) {
	// Arrange
	order := &Order{
		ID:        1,
		UserID:    1,
		Total:     0,
		Status:    "pending",
		Products:  nil, // nil вместо массива
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	jsonData, err := json.Marshal(order)
	assert.NoError(t, err)

	var decodedOrder Order
	err = json.Unmarshal(jsonData, &decodedOrder)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, order.ID, decodedOrder.ID)
	assert.Equal(t, order.UserID, decodedOrder.UserID)
	assert.Equal(t, order.Total, decodedOrder.Total)
	assert.Equal(t, order.Status, decodedOrder.Status)
	assert.Nil(t, decodedOrder.Products) // Должен остаться nil
}
