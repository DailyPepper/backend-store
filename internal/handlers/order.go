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

type OrderHander struct {
	storage storage.Storage
}

func NewOrderHandler(storage storage.Storage) *OrderHander {
	return &OrderHander{storage: storage}
}

func (h *OrderHander) CreateOrder(c *gin.Context) {
	var order models.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	if order.TotalAmount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "TotalAmount <= 0",
		})
	}

	if err := h.storage.CreateOrder(&order); err != nil {
		log.Printf("Error creating order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create order",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"order":   order,
	})

}

func (h *OrderHander) GetAllOrders(c *gin.Context) {
	orders, err := h.storage.GetAllOrders()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения заказов"})
		return
	}

	response, err := json.Marshal(orders)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка форматиорвания ответа"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *OrderHander) GetOrderID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.storage.GetOrderByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHander) DeleteOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	err = h.storage.DeleteOrder(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
