package handlers

import (
	"authService/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

func NewUserHandler(authSvc *services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authSvc}
}

func (h AuthHandler) Register(c *gin.Context) {

	var body struct {
		FulllName string `json:"full_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	err := h.AuthService.Register(body.FulllName, body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (h AuthHandler) Login(c *gin.Context) {

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	userToken, exp, err := h.AuthService.Login(body.Email, body.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": userToken,
		"expires_at":   exp,
	})
}
