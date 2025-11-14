package middleware

import (
	"backend-store/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authClient *auth.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Пропускаем публичные routes
		publicPaths := []string{"/health", "/api/auth", "/openapi.yaml", "/swagger"}
		for _, path := range publicPaths {
			if strings.HasPrefix(c.Request.URL.Path, path) {
				c.Next()
				return
			}
		}

		// GET запросы к products тоже публичные
		if c.Request.Method == "GET" && strings.HasPrefix(c.Request.URL.Path, "/api/product") {
			c.Next()
			return
		}

		// Проверяем токен для защищенных routes
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		resp, err := authClient.ValidateToken(c.Request.Context(), token)
		if err != nil || !resp.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Сохраняем user_id в контекст для использования в handlers
		c.Set("user_id", resp.UserId)
		c.Set("user_email", resp.Email)
		c.Next()
	}
}
