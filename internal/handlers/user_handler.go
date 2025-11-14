package handlers

import (
	"backend-store/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Заглушки - можно реализовать позже
func (h *UserHandler) GetProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "User profile endpoint - to be implemented",
		"user_id": "123", // TODO: получить из контекста
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Update profile endpoint - to be implemented",
	})
}
