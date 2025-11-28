package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandlers struct {
	authClient *AuthClient
}

func NewAuthHandlers(authClient *AuthClient) *AuthHandlers {
	return &AuthHandlers{
		authClient: authClient,
	}
}

func (h *AuthHandlers) RegisterHandler(c *gin.Context) {
	var req struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		Surname   string `json:"surname"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	registerResp, err := h.authClient.Register(c.Request.Context(), req.Email, req.Password, req.FirstName, req.Surname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	loginResp, err := h.authClient.Login(c.Request.Context(), req.Email, req.Password)
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
		"expires_at":    loginResp.ExpiresAt,
	})
}

func (h *AuthHandlers) LoginHandler(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := h.authClient.Login(c.Request.Context(), req.Email, req.Password)
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
