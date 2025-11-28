package server

import (
	"backend-store/config"
	"backend-store/internal/app"
	"backend-store/internal/auth"
	"backend-store/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, handlers *app.Handlers, authHandlers *auth.AuthHandlers, authClient *auth.AuthClient) *gin.Engine {
	setupGinMode(cfg)

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	SetupSwagger(router)
	setupAuthRoutes(router, authHandlers)
	setupAPIRoutes(router, handlers, authClient)

	return router
}

func setupGinMode(cfg *config.Config) {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

func setupAuthRoutes(router *gin.Engine, authHandlers *auth.AuthHandlers) {
	router.POST("/api/auth/register", authHandlers.RegisterHandler)
	router.POST("/api/auth/login", authHandlers.LoginHandler)
}

func setupAPIRoutes(router *gin.Engine, handlers *app.Handlers, authClient *auth.AuthClient) {
	api := router.Group("/api")
	{
		setupProductRoutes(api, handlers, authClient)
		setupOrderRoutes(api, handlers, authClient)
	}
}

func setupProductRoutes(api *gin.RouterGroup, handlers *app.Handlers, authClient *auth.AuthClient) {
	product := api.Group("/product")
	{
		product.GET("/", handlers.ProductHandler.GetAllProducts)
		product.GET("/:id", handlers.ProductHandler.GetProductByID)

		protected := product.Group("")
		protected.Use(middleware.AuthMiddleware(authClient))
		{
			protected.POST("/", handlers.ProductHandler.CreateProduct)
			protected.PUT("/:id", handlers.ProductHandler.UpdateProduct)
			protected.DELETE("/:id", handlers.ProductHandler.DeleteProduct)
		}
	}
}

func setupOrderRoutes(api *gin.RouterGroup, handlers *app.Handlers, authClient *auth.AuthClient) {
	order := api.Group("/order")
	order.Use(middleware.AuthMiddleware(authClient))
	{
		order.GET("/", handlers.OrderHandler.GetAllOrders)
		order.POST("/", handlers.OrderHandler.CreateOrder)
		order.GET("/:id", handlers.OrderHandler.GetOrderByID)
		order.PUT("/:id", handlers.OrderHandler.UpdateOrder)
		order.DELETE("/:id", handlers.OrderHandler.DeleteOrder)
	}
}
