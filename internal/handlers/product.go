package handlers

import (
	"backend-store/internal/models"
	"backend-store/internal/storage"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	storage storage.Storage
}

func NewProductHandler(storage storage.Storage) *ProductHandler {
	return &ProductHandler{storage: storage}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Валидация
	if product.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product name is required"})
		return
	}
	if product.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price must be positive"})
		return
	}

	tx, err := h.storage.BeginTx(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	if err := tx.CreateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product: " + err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	products, err := h.storage.GetAllProduct()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.storage.GetProductByID(id)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}
	product.ID = id

	tx, err := h.storage.BeginTx(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Проверяем существование продукта
	existingProduct, err := tx.GetProductByID(id)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product: " + err.Error()})
		}
		return
	}

	// Сохраняем createdAt из существующего продукта
	product.CreatedAt = existingProduct.CreatedAt

	if err := tx.UpdateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product: " + err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	tx, err := h.storage.BeginTx(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Проверяем существование продукта
	_, err = tx.GetProductByID(id)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product: " + err.Error()})
		}
		return
	}

	if err := tx.DeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product: " + err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
