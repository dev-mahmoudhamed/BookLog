package handlers

import (
	service "bookLog/internal/services/user"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc service.UserService
	authSvc *service.AuthService
}

type registerRequest struct {
	FullName string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func NewUserHandler(userSvc service.UserService, authSvc *service.AuthService) *UserHandler {
	return &UserHandler{userSvc: userSvc, authSvc: authSvc}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	search := c.Query("search")
	users, total, err := h.userSvc.GetAllUsers(search)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  users,
		"total": total,
	})
}

func (h *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.authSvc.Register(req.FullName, req.Email, req.Password)
	if err != nil {
		le := strings.ToLower(err.Error())
		if strings.Contains(le, "exists") || strings.Contains(le, "duplicate") {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return user without password
	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id":         u.ID,
			"full_name":  u.FullName,
			"email":      u.Email,
			"role":       u.Role,
			"created_at": u.CreatedAt,
			"updated_at": u.UpdatedAt,
		},
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authSvc.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
