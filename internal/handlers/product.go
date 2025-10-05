package handlers

import (
	"backend-store/internal/models"
	"backend-store/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := h.productService.CreateProduct(c.Request.Context(), &product); err != nil {
		if contains(err.Error(), "validate") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	products, err := h.productService.GetAllProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products: " + err.Error()})
		return
	}

	start := (page - 1) * limit
	end := start + limit

	var paginatedProducts []*models.Product
	if start < len(products) {
		if end > len(products) {
			end = len(products)
		}
		paginatedProducts = products[start:end]
	} else {
		paginatedProducts = []*models.Product{}
	}

	c.JSON(http.StatusOK, gin.H{
		"products": paginatedProducts,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": len(products),
			"pages": (len(products) + limit - 1) / limit,
		},
	})
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.productService.GetProductByID(c.Request.Context(), id)
	if err != nil {
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
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

	if err := h.productService.UpdateProduct(c.Request.Context(), &product); err != nil {
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else if contains(err.Error(), "validate") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := h.productService.DeleteProduct(c.Request.Context(), id); err != nil {
		if contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else if contains(err.Error(), "cannot delete") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
