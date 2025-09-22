package main

import (
	"backend-store/internal/handlers"
	"backend-store/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	store := storage.NewMemoryStorage()
	orderHandler := handlers.NewOrderHandler(store)
	productHandler := handlers.NewProductHandler(store)

	api := router.Group("/api")
	{
		order := api.Group("/order")
		{
			order.POST("/", orderHandler.CreateOrder)
			order.GET("/", orderHandler.GetAllOrders)
			order.GET("/:id", orderHandler.GetOrderID)
			order.DELETE("/:id", orderHandler.DeleteOrder)
		}

		product := api.Group("/product")
		{
			product.POST("/", productHandler.CreateProduct)
			product.GET("/", productHandler.GetAllProduct)
			product.GET("/:id", productHandler.GetProductByID)
			product.DELETE("/:id", productHandler.DeleteProduct)
		}
	}

	router.Run(":8080")
}
