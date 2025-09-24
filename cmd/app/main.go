package main

import (
	"backend-store/internal/config"
	"backend-store/internal/handlers"
	"backend-store/internal/storage"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	var store storage.Storage

	if cfg.DatabaseURL != "" {
		log.Println("Database URL:", cfg.DatabaseURL)
		postgresStore, err := storage.NewPostgresStorage(cfg.DatabaseURL)
		if err != nil {
			log.Fatal("Failed to connect to PostgreSQL:", err)
		}
		defer postgresStore.Close()

		if err := postgresStore.Init(); err != nil {
			log.Fatal("Failed to init database:", err)
		}

		store = postgresStore
		log.Println("Using PostgreSQL storage")
	} else {
		store = storage.NewMemoryStorage()
		log.Println("Using in-memory storage")
	}

	orderHandler := handlers.NewOrderHandler(store)
	productHandler := handlers.NewProductHandler(store)

	router := gin.Default()

	api := router.Group("/api")
	{
		order := api.Group("/order")
		{
			order.POST("/", orderHandler.CreateOrder)
			order.GET("/", orderHandler.GetAllOrders)
			order.GET("/:id", orderHandler.GetOrderByID)
			order.PUT("/:id", orderHandler.UpdateOrder)
			order.DELETE("/:id", orderHandler.DeleteOrder)
		}

		product := api.Group("/product")
		{
			product.POST("/", productHandler.CreateProduct)
			product.GET("/", productHandler.GetAllProducts)
			product.GET("/:id", productHandler.GetProductByID)
			product.PUT("/:id", productHandler.UpdateProduct)
			product.DELETE("/:id", productHandler.DeleteProduct)
		}
	}

	log.Println("Server starting on :8080")
	router.Run(":8080")
}
