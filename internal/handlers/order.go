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

type OrderHandler struct {
	storage storage.Storage
}

func NewOrderHandler(s storage.Storage) *OrderHandler {
	return &OrderHandler{storage: s}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if order.CustomerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer name is required"})
		return
	}
	if order.TotalAmount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Total amount must be positive"})
		return
	}

	tx, err := h.storage.BeginTx(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	if order.Status == "" {
		order.Status = "pending"
	}

	if err := tx.CreateOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order: " + err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	orders, err := h.storage.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders: " + err.Error()})
		return
	}

	start := (page - 1) * limit
	end := start + limit
	if start >= len(orders) {
		c.JSON(http.StatusOK, gin.H{
			"orders": []*models.Order{},
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": len(orders),
			},
		})
		return
	}
	if end > len(orders) {
		end = len(orders)
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders[start:end],
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": len(orders),
		},
	})
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.storage.GetOrderByID(id)
	if err != nil {
		if err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, order)
}

// UpdateOrder обновляет заказ с транзакцией
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}
	order.ID = id

	// Начинаем транзакцию
	tx, err := h.storage.BeginTx(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Проверяем существование заказа
	existingOrder, err := tx.GetOrderByID(id)
	if err != nil {
		if err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order: " + err.Error()})
		}
		return
	}

	// Сохраняем createdAt из существующего заказа
	order.CreatedAt = existingOrder.CreatedAt

	if err := tx.UpdateOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order: " + err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// DeleteOrder удаляет заказ с транзакцией
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Начинаем транзакцию
	tx, err := h.storage.BeginTx(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Проверяем существование заказа
	_, err = tx.GetOrderByID(id)
	if err != nil {
		if err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order: " + err.Error()})
		}
		return
	}

	if err := tx.DeleteOrder(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order: " + err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
