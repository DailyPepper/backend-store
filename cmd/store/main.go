package main

import (
	"backend-store/config"
	"backend-store/internal/app"
	"backend-store/pkg/logger"
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var log logger.Log

func main() {
	cfg := config.Load()
	log = logger.New(cfg.LogLevel)

	setupLogging(cfg)

	log.Info("Starting application", "mode", cfg.Environment)
	log.Info("Server will start on", "address", cfg.ServerHost+":"+cfg.ServerPort)

	// Передаем логгер как второй параметр
	application, err := app.New(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}
	defer application.Close()

	router := setupRouter(application.Handlers)

	startServer(cfg, router)
}

func setupLogging(cfg *config.Config) {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

func setupRouter(handlers *app.Handlers) *gin.Engine {
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	setupSwagger(router)

	router.GET("/health", healthCheck)

	api := router.Group("/api")
	{
		order := api.Group("/order")
		{
			order.POST("/", handlers.OrderHandler.CreateOrder)
			order.GET("/", handlers.OrderHandler.GetAllOrders)
			order.GET("/:id", handlers.OrderHandler.GetOrderByID)
			order.PUT("/:id", handlers.OrderHandler.UpdateOrder)
			order.DELETE("/:id", handlers.OrderHandler.DeleteOrder)
		}

		product := api.Group("/product")
		{
			product.POST("/", handlers.ProductHandler.CreateProduct)
			product.GET("/", handlers.ProductHandler.GetAllProducts)
			product.GET("/:id", handlers.ProductHandler.GetProductByID)
			product.PUT("/:id", handlers.ProductHandler.UpdateProduct)
			product.DELETE("/:id", handlers.ProductHandler.DeleteProduct)
		}
	}

	return router
}

func setupSwagger(router *gin.Engine) {
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
		ginSwagger.DefaultModelsExpandDepth(-1),
	))
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "OK",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func startServer(cfg *config.Config, router *gin.Engine) {
	srv := &http.Server{
		Addr:         cfg.ServerHost + ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	go func() {
		log.Info("Server starting on", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Info("Server exited")
}
