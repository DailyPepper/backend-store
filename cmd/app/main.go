package main

import (
	"backend-store/internal/handlers"
	"backend-store/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	store := storage.NewMemoryStorage()
	orderHandler := handlers.NewOrderHandler(store)

	api := router.Group("/api")
	{
		order := api.Group("/order")
		{
			// Обертываем http.Handler в gin.HandlerFunc
			order.GET("/", gin.WrapH(http.HandlerFunc(orderHandler.GetAllOrders)))
		}
	}

	router.Run(":8080")
}
