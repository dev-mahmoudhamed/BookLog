package handlers

import (
	"net/http"
	"userService/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

type UserRegisterDto struct {
	FulllName string `json:"full_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func NewUserHandler(usrSrv *services.UserService) *UserHandler {
	return &UserHandler{userService: usrSrv}
}

func (h UserHandler) Register(c *gin.Context) {

	var body UserRegisterDto

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	err := h.userService.Register(body.FulllName, body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (h UserHandler) Login(c *gin.Context) {

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	userToken, exp, err := h.userService.Login(body.Email, body.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": userToken,
		"expires_at":   exp,
	})
}
