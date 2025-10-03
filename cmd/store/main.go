package main

import (
	"backend-store/internal/config"
	"backend-store/internal/handlers"
	"backend-store/internal/storage"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	router.GET("/openapi.yaml", func(c *gin.Context) {
		openAPIPath := filepath.Join("api", "openapi.yaml")

		content, err := ioutil.ReadFile(openAPIPath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "OpenAPI spec not found"})
			return
		}

		c.Data(http.StatusOK, "application/yaml; charset=utf-8", content)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/openapi.yaml"),
		ginSwagger.DeepLinking(true),
	))

	router.GET("/docs/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/openapi.yaml"),
		ginSwagger.DeepLinking(true),
	))

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

	log.Println("Docs UI: http://localhost:8080/docs/index.html")
	router.Run(":8080")
}
