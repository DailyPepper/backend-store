package handlers

import (
	"backend-store/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderService) GetAllOrders(ctx context.Context) ([]*models.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Order), args.Error(1)
}

func (m *MockOrderService) GetOrderByID(ctx context.Context, id int) (*models.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderService) UpdateOrder(ctx context.Context, order *models.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderService) DeleteOrder(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestOrderHandler_CreateOrder_Success(t *testing.T) {
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)
	router := setupRouter()
	router.POST("/orders", handler.CreateOrder)

	orderRequest := map[string]interface{}{
		"user_id": 1,
		"total":   9999,
		"status":  "pending",
		"products": []map[string]interface{}{
			{
				"product_id": 1,
				"quantity":   2,
				"price":      5000,
			},
		},
	}

	mockService.On("CreateOrder", mock.Anything, mock.AnythingOfType("*models.Order")).
		Return(nil).
		Run(func(args mock.Arguments) {
			order := args.Get(1).(*models.Order)
			order.ID = 1
			order.CreatedAt = time.Now()
			order.UpdatedAt = time.Now()
		})

	body, _ := json.Marshal(orderRequest)
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Order
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 1, response.ID)
	assert.Equal(t, 9999, response.Total)
	assert.Equal(t, "pending", response.Status)
	mockService.AssertExpectations(t)
}

func TestOrderHandler_CreateOrder_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)
	router := setupRouter()
	router.POST("/orders", handler.CreateOrder)

	// Act
	body := bytes.NewBufferString(`{"invalid_json"`)
	req, _ := http.NewRequest("POST", "/orders", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

func TestOrderHandler_CreateOrder_ProductNotFound(t *testing.T) {
	// Arrange
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)
	router := setupRouter()
	router.POST("/orders", handler.CreateOrder)

	orderRequest := map[string]interface{}{
		"user_id": 1,
		"total":   9999,
		"status":  "pending",
		"products": []map[string]interface{}{
			{
				"product_id": 1,
				"quantity":   2,
				"price":      5000,
			},
		},
	}

	mockService.On("CreateOrder", mock.Anything, mock.AnythingOfType("*models.Order")).
		Return(errors.New("product 1 not found"))

	// Act
	body, _ := json.Marshal(orderRequest)
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "product 1 not found")
	mockService.AssertExpectations(t)
}

func TestOrderHandler_GetAllOrders_Success(t *testing.T) {
	// Arrange
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)
	router := setupRouter()
	router.GET("/orders", handler.GetAllOrders)

	expectedOrders := []*models.Order{
		{
			ID:     1,
			Total:  9999,
			Status: "completed",
			Products: []models.OrderItem{
				{ProductID: 1, Quantity: 2, Price: 5000},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:     2,
			Total:  14999,
			Status: "pending",
			Products: []models.OrderItem{
				{ProductID: 2, Quantity: 1, Price: 14999},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockService.On("GetAllOrders", mock.Anything).Return(expectedOrders, nil)

	// Act
	req, _ := http.NewRequest("GET", "/orders?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	orders := response["orders"].([]interface{})
	pagination := response["pagination"].(map[string]interface{})

	assert.Len(t, orders, 2)
	assert.Equal(t, float64(1), pagination["page"])
	assert.Equal(t, float64(10), pagination["limit"])
	assert.Equal(t, float64(2), pagination["total"])
	mockService.AssertExpectations(t)
}

func TestOrderHandler_GetAllOrders_Pagination(t *testing.T) {
	// Arrange
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)
	router := setupRouter()
	router.GET("/orders", handler.GetAllOrders)

	orders := make([]*models.Order, 15)
	for i := 0; i < 15; i++ {
		orders[i] = &models.Order{
			ID:     i + 1,
			Total:  1000,
			Status: "pending",
		}
	}

	mockService.On("GetAllOrders", mock.Anything).Return(orders, nil)

	// Act
	req, _ := http.NewRequest("GET", "/orders?page=2&limit=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	returnedOrders := response["orders"].([]interface{})
	pagination := response["pagination"].(map[string]interface{})

	assert.Len(t, returnedOrders, 5)
	assert.Equal(t, float64(2), pagination["page"])
	assert.Equal(t, float64(5), pagination["limit"])
	assert.Equal(t, float64(15), pagination["total"])
	assert.Equal(t, float64(3), pagination["pages"])
	mockService.AssertExpectations(t)
}

func TestOrderHandler_GetOrderByID_Success(t *testing.T) {
	// Arrange
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)
	router := setupRouter()
	router.GET("/orders/:id", handler.GetOrderByID)

	expectedOrder := &models.Order{
		ID:     1,
		Total:  9999,
		Status: "completed",
		Products: []models.OrderItem{
			{ProductID: 1, Quantity: 2, Price: 5000},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("GetOrderByID", mock.Anything, 1).Return(expectedOrder, nil)

	// Act
	req, _ := http.NewRequest("GET", "/orders/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Order
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, expectedOrder.ID, response.ID)
	assert.Equal(t, expectedOrder.Total, response.Total)
	mockService.AssertExpectations(t)
}

func TestOrderHandler_GetOrderByID_InvalidID(t *testing.T) {
	// Arrange
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)
	router := setupRouter()
	router.GET("/orders/:id", handler.GetOrderByID)

	// Act
	req, _ := http.NewRequest("GET", "/orders/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid order ID")
}

func TestOrderHandler_GetOrderByID_NotFound(t *testing.T) {
	// Arrange
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)
	router := setupRouter()
	router.GET("/orders/:id", handler.GetOrderByID)

	mockService.On("GetOrderByID", mock.Anything, 999).Return((*models.Order)(nil), errors.New("order not found"))

	// Act
	req, _ := http.NewRequest("GET", "/orders/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Order not found")
	mockService.AssertExpectations(t)
}

func TestOrderHandler_UpdateOrder_Success(t *testing.T) {
	// Arrange
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)
	router := setupRouter()
	router.PUT("/orders/:id", handler.UpdateOrder)

	orderRequest := map[string]interface{}{
		"user_id": 1,
		"total":   14999,
		"status":  "completed",
		"products": []map[string]interface{}{
			{
				"product_id": 1,
				"quantity":   3,
				"price":      5000,
			},
		},
	}

	mockService.On("UpdateOrder", mock.Anything, mock.AnythingOfType("*models.Order")).
		Return(nil)

	// Act
	body, _ := json.Marshal(orderRequest)
	req, _ := http.NewRequest("PUT", "/orders/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Order
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 1, response.ID)
	assert.Equal(t, 14999, response.Total)
	assert.Equal(t, "completed", response.Status)
	mockService.AssertExpectations(t)
}

func TestOrderHandler_DeleteOrder_Success(t *testing.T) {
	// Arrange
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)
	router := setupRouter()
	router.DELETE("/orders/:id", handler.DeleteOrder)

	mockService.On("DeleteOrder", mock.Anything, 1).Return(nil)

	// Act
	req, _ := http.NewRequest("DELETE", "/orders/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Order deleted successfully")
	mockService.AssertExpectations(t)
}

func TestContains(t *testing.T) {
	assert.True(t, contains("product not found", "product"))
	assert.True(t, contains("insufficient quantity", "insufficient"))
	assert.True(t, contains("validate error", "validate"))
	assert.False(t, contains("some error", "product"))
	assert.False(t, contains("", "product"))
}
