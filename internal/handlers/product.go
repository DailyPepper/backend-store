package handlers

import (
	"backend-store/internal/models"
	"backend-store/internal/storage"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	storage storage.Storage
}

func NewProductHandler(storage storage.Storage) *ProductHandler {
	return &ProductHandler{storage: storage}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Валидация
	if product.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Product name is required",
		})
		return
	}

	if product.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Price must be greater than 0",
		})
		return
	}

	if product.Stock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Stock cannot be negative",
		})
		return
	}

	if err := h.storage.CreateProduct(&product); err != nil {
		log.Printf("Error creating product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create product",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"product": product,
	})
}

func (h *ProductHandler) GetAllProduct(c *gin.Context) {
	products, err := h.storage.GetAllProduct()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Нет такого товара"})
		return
	}

	response, err := json.Marshal(products)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка форматиорвания ответа"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.storage.GetProductByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	err = h.storage.DeleteProduct(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
