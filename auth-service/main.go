package main

import (
	"authService/config"
	"authService/database"
	"authService/handlers"
	"authService/internal/repository"
	"authService/internal/services"
	"authService/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("❌ Failed to load config:", err)
	}

	// Connect to DB using config
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("❌ Database connection error:", err)
	}
	defer db.Close()

	authRepo := repository.NewUserRepositoryPostgres(db)
	authService := services.NewAuthService(authRepo, cfg.JwtSecret)
	authHandler := handlers.NewUserHandler(authService)

	r := gin.Default()

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware(cfg.JwtSecret))

	auth.GET("/secret", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		c.JSON(200, gin.H{
			"message": "This is a secret endpoint!",
			"user_id": userID,
		})
	})

	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
