package handlers

import (
	"backend-store/internal/models"
	"backend-store/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := h.orderService.CreateOrder(c.Request.Context(), &order); err != nil {
		switch {
		case contains(err.Error(), "product") && contains(err.Error(), "not found"):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case contains(err.Error(), "insufficient quantity"):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case contains(err.Error(), "validate"):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	orders, err := h.orderService.GetAllOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders: " + err.Error()})
		return
	}

	start := (page - 1) * limit
	end := start + limit

	var paginatedOrders []*models.Order
	if start < len(orders) {
		if end > len(orders) {
			end = len(orders)
		}
		paginatedOrders = orders[start:end]
	} else {
		paginatedOrders = []*models.Order{}
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": paginatedOrders,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": len(orders),
			"pages": (len(orders) + limit - 1) / limit,
		},
	})
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.orderService.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
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

	if err := h.orderService.UpdateOrder(c.Request.Context(), &order); err != nil {
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else if contains(err.Error(), "validate") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	if err := h.orderService.DeleteOrder(c.Request.Context(), id); err != nil {
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
