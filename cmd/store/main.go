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

	"backend-store/internal/auth"
	"backend-store/internal/middleware"
)

var log logger.Log

func main() {
	cfg := config.Load()
	log = logger.New(cfg.LogLevel)

	setupLogging(cfg)

	log.Info("Starting application", "mode", cfg.Environment)
	log.Info("Server will start on", "address", cfg.ServerHost+":"+cfg.ServerPort)

	// Подключаемся к auth-service
	authClient, err := auth.NewAuthClient(cfg.AuthServiceURL)
	if err != nil {
		log.Fatal("Failed to connect to auth service:", err)
	}
	defer authClient.Close()
	log.Info("✅ Successfully connected to auth service")

	// Инициализируем приложение (без authClient)
	application, err := app.New(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}
	defer application.Close()

	// Настраиваем роутер с middleware аутентификации
	router := setupRouter(application.Handlers, authClient)
	startServer(cfg, router)
}

func setupLogging(cfg *config.Config) {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

func setupRouter(handlers *app.Handlers, authClient *auth.AuthClient) *gin.Engine {
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	setupSwagger(router)

	// Public routes
	router.GET("/health", healthCheck)

	// Auth routes (простые handlers напрямую)
	router.POST("/api/auth/register", createRegisterHandler(authClient))
	router.POST("/api/auth/login", createLoginHandler(authClient))

	api := router.Group("/api")
	{
		// Product routes - публичные для GET, защищенные для остального
		product := api.Group("/product")
		{
			product.GET("/", handlers.ProductHandler.GetAllProducts)
			product.GET("/:id", handlers.ProductHandler.GetProductByID)

			// Защищенные routes
			product.Use(middleware.AuthMiddleware(authClient))
			{
				product.POST("/", handlers.ProductHandler.CreateProduct)
				product.PUT("/:id", handlers.ProductHandler.UpdateProduct)
				product.DELETE("/:id", handlers.ProductHandler.DeleteProduct)
			}
		}

		// Order routes - полностью защищенные
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

	return router
}

func createRegisterHandler(authClient *auth.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email     string `json:"email"`
			Password  string `json:"password"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Сначала регистрируем
		registerResp, err := authClient.Register(c.Request.Context(), req.Email, req.Password, req.FirstName, req.LastName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Затем логинимся чтобы получить токен
		loginResp, err := authClient.Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "User registered but login failed",
				"user_id": registerResp.Id,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":       "User registered successfully",
			"user_id":       registerResp.Id,
			"access_token":  loginResp.AccessToken,
			"refresh_token": loginResp.RefreshToken,
		})
	}
}

func createLoginHandler(authClient *auth.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		resp, err := authClient.Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  resp.AccessToken,
			"refresh_token": resp.RefreshToken,
			"expires_at":    resp.ExpiresAt,
		})
	}
}

// Остальные функции без изменений...
func setupSwagger(router *gin.Engine) {
	router.GET("/openapi.yaml", func(c *gin.Context) {
		openAPIPath := filepath.Join("docs", "openapi.yaml")
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
		"service":   "backend-store",
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
