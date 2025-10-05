package main

import (
	"backend-store/config"
	"backend-store/internal/app"
	"context"
	"io/ioutil"
	"log"
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

func main() {
	cfg := config.Load()

	setupLogging(cfg)

	log.Printf("Starting application in %s mode", cfg.Environment)
	log.Printf("Server will start on %s:%s", cfg.ServerHost, cfg.ServerPort)

	application, err := app.New(cfg)
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
		log.SetOutput(os.Stdout)
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
		log.Printf("Server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
