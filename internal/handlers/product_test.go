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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductService реализует интерфейс service.ProductService для тестов
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductService) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *MockProductService) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) UpdateProduct(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductService) DeleteProduct(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestProductHandler_CreateProduct_Success(t *testing.T) {
	// Arrange
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.POST("/products", handler.CreateProduct)

	productRequest := map[string]interface{}{
		"name":        "Test Product",
		"description": "Test Description",
		"price":       29.99,
		"quantity":    100,
	}

	mockService.On("CreateProduct", mock.Anything, mock.AnythingOfType("*models.Product")).
		Return(nil).
		Run(func(args mock.Arguments) {
			product := args.Get(1).(*models.Product)
			product.ID = 1
			product.CreatedAt = time.Now()
			product.UpdatedAt = time.Now()
		})

	// Act
	body, _ := json.Marshal(productRequest)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Product
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "Test Product", response.Name)
	assert.Equal(t, "Test Description", response.Description)
	assert.Equal(t, 29.99, response.Price)
	assert.Equal(t, 100, response.Quantity)
	mockService.AssertExpectations(t)
}

func TestProductHandler_CreateProduct_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.POST("/products", handler.CreateProduct)

	// Act
	body := bytes.NewBufferString(`{"invalid_json"`)
	req, _ := http.NewRequest("POST", "/products", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

func TestProductHandler_CreateProduct_ValidationError(t *testing.T) {
	// Arrange
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.POST("/products", handler.CreateProduct)

	productRequest := map[string]interface{}{
		"name":        "", // Пустое имя - должно вызвать ошибку валидации
		"description": "Test Description",
		"price":       29.99,
		"quantity":    100,
	}

	mockService.On("CreateProduct", mock.Anything, mock.AnythingOfType("*models.Product")).
		Return(errors.New("validate: name is required"))

	// Act
	body, _ := json.Marshal(productRequest)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "validate")
	mockService.AssertExpectations(t)
}

func TestProductHandler_CreateProduct_AlreadyExists(t *testing.T) {
	// Arrange
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.POST("/products", handler.CreateProduct)

	productRequest := map[string]interface{}{
		"name":        "Test Product",
		"description": "Test Description",
		"price":       29.99,
		"quantity":    100,
	}

	mockService.On("CreateProduct", mock.Anything, mock.AnythingOfType("*models.Product")).
		Return(errors.New("product already exists"))

	// Act
	body, _ := json.Marshal(productRequest)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "already exists")
	mockService.AssertExpectations(t)
}

func TestProductHandler_CreateProduct_InternalError(t *testing.T) {
	// Arrange
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.POST("/products", handler.CreateProduct)

	productRequest := map[string]interface{}{
		"name":        "Test Product",
		"description": "Test Description",
		"price":       29.99,
		"quantity":    100,
	}

	mockService.On("CreateProduct", mock.Anything, mock.AnythingOfType("*models.Product")).
		Return(errors.New("database error"))

	// Act
	body, _ := json.Marshal(productRequest)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to create product")
	mockService.AssertExpectations(t)
}

func TestProductHandler_GetAllProducts_Success(t *testing.T) {
	// Arrange
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.GET("/products", handler.GetAllProducts)

	expectedProducts := []*models.Product{
		{
			ID:          1,
			Name:        "Product 1",
			Description: "Description 1",
			Price:       29.99,
			Quantity:    100,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Name:        "Product 2",
			Description: "Description 2",
			Price:       49.99,
			Quantity:    50,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockService.On("GetAllProducts", mock.Anything).Return(expectedProducts, nil)

	// Act
	req, _ := http.NewRequest("GET", "/products?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	products := response["products"].([]interface{})
	pagination := response["pagination"].(map[string]interface{})

	assert.Len(t, products, 2)
	assert.Equal(t, float64(1), pagination["page"])
	assert.Equal(t, float64(10), pagination["limit"])
	assert.Equal(t, float64(2), pagination["total"])
	mockService.AssertExpectations(t)
}

func TestProductHandler_GetAllProducts_Pagination(t *testing.T) {
	// Arrange
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.GET("/products", handler.GetAllProducts)

	products := make([]*models.Product, 15)
	for i := 0; i < 15; i++ {
		products[i] = &models.Product{
			ID:          i + 1,
			Name:        "Product",
			Description: "Description",
			Price:       10.0,
			Quantity:    10,
		}
	}

	mockService.On("GetAllProducts", mock.Anything).Return(products, nil)

	// Act
	req, _ := http.NewRequest("GET", "/products?page=2&limit=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	returnedProducts := response["products"].([]interface{})
	pagination := response["pagination"].(map[string]interface{})

	assert.Len(t, returnedProducts, 5)
	assert.Equal(t, float64(2), pagination["page"])
	assert.Equal(t, float64(5), pagination["limit"])
	assert.Equal(t, float64(15), pagination["total"])
	assert.Equal(t, float64(3), pagination["pages"])
	mockService.AssertExpectations(t)
}

func TestProductHandler_GetAllProducts_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.GET("/products", handler.GetAllProducts)

	mockService.On("GetAllProducts", mock.Anything).Return(([]*models.Product)(nil), errors.New("database error"))

	// Act
	req, _ := http.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to fetch products")
	mockService.AssertExpectations(t)
}

func TestProductHandler_GetProductByID_Success(t *testing.T) {
	// Arrange
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.GET("/products/:id", handler.GetProductByID)

	expectedProduct := &models.Product{
		ID:          1,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       29.99,
		Quantity:    100,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockService.On("GetProductByID", mock.Anything, 1).Return(expectedProduct, nil)

	// Act
	req, _ := http.NewRequest("GET", "/products/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Product
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, expectedProduct.ID, response.ID)
	assert.Equal(t, expectedProduct.Name, response.Name)
	assert.Equal(t, expectedProduct.Description, response.Description)
	assert.Equal(t, expectedProduct.Price, response.Price)
	mockService.AssertExpectations(t)
}

func TestProductHandler_GetProductByID_InvalidID(t *testing.T) {
	// Arrange
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.GET("/products/:id", handler.GetProductByID)

	// Act
	req, _ := http.NewRequest("GET", "/products/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid product ID")
}

func TestProductHandler_GetProductByID_NotFound(t *testing.T) {
	// Arrange
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.GET("/products/:id", handler.GetProductByID)

	mockService.On("GetProductByID", mock.Anything, 999).Return((*models.Product)(nil), errors.New("product not found"))

	// Act
	req, _ := http.NewRequest("GET", "/products/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Product not found")
	mockService.AssertExpectations(t)
}

func TestProductHandler_UpdateProduct_Success(t *testing.T) {
	// Arrange
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.PUT("/products/:id", handler.UpdateProduct)

	productRequest := map[string]interface{}{
		"name":        "Updated Product",
		"description": "Updated Description",
		"price":       39.99,
		"quantity":    150,
	}

	mockService.On("UpdateProduct", mock.Anything, mock.AnythingOfType("*models.Product")).
		Return(nil)

	// Act
	body, _ := json.Marshal(productRequest)
	req, _ := http.NewRequest("PUT", "/products/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Product
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "Updated Product", response.Name)
	assert.Equal(t, "Updated Description", response.Description)
	assert.Equal(t, 39.99, response.Price)
	assert.Equal(t, 150, response.Quantity)
	mockService.AssertExpectations(t)
}

func TestProductHandler_UpdateProduct_NotFound(t *testing.T) {
	// Arrange
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.PUT("/products/:id", handler.UpdateProduct)

	productRequest := map[string]interface{}{
		"name":        "Updated Product",
		"description": "Updated Description",
		"price":       39.99,
		"quantity":    150,
	}

	mockService.On("UpdateProduct", mock.Anything, mock.AnythingOfType("*models.Product")).
		Return(errors.New("product not found"))

	// Act
	body, _ := json.Marshal(productRequest)
	req, _ := http.NewRequest("PUT", "/products/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Product not found")
	mockService.AssertExpectations(t)
}

func TestProductHandler_UpdateProduct_ValidationError(t *testing.T) {
	// Arrange
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.PUT("/products/:id", handler.UpdateProduct)

	productRequest := map[string]interface{}{
		"name":        "", // Пустое имя
		"description": "Updated Description",
		"price":       39.99,
		"quantity":    150,
	}

	mockService.On("UpdateProduct", mock.Anything, mock.AnythingOfType("*models.Product")).
		Return(errors.New("validate: name is required"))

	// Act
	body, _ := json.Marshal(productRequest)
	req, _ := http.NewRequest("PUT", "/products/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "validate")
	mockService.AssertExpectations(t)
}

func TestProductHandler_DeleteProduct_Success(t *testing.T) {
	// Arrange
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.DELETE("/products/:id", handler.DeleteProduct)

	mockService.On("DeleteProduct", mock.Anything, 1).Return(nil)

	// Act
	req, _ := http.NewRequest("DELETE", "/products/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Product deleted successfully")
	mockService.AssertExpectations(t)
}

func TestProductHandler_DeleteProduct_NotFound(t *testing.T) {
	// Arrange
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.DELETE("/products/:id", handler.DeleteProduct)

	mockService.On("DeleteProduct", mock.Anything, 999).Return(errors.New("product not found"))

	// Act
	req, _ := http.NewRequest("DELETE", "/products/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Product not found")
	mockService.AssertExpectations(t)
}

func TestProductHandler_DeleteProduct_CannotDelete(t *testing.T) {
	// Arrange
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter()
	router.DELETE("/products/:id", handler.DeleteProduct)

	mockService.On("DeleteProduct", mock.Anything, 1).Return(errors.New("cannot delete product with existing orders"))

	// Act
	req, _ := http.NewRequest("DELETE", "/products/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "cannot delete")
	mockService.AssertExpectations(t)
}
