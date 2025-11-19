package main

import (
	"log"
	"userService/config"
	"userService/database"
	"userService/handlers"
	"userService/internal/repository"
	"userService/internal/services"
	"userService/middleware"

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

	userRepo := repository.NewUserRepositoryPostgres(db)
	userService := services.NewUserService(userRepo, cfg.JwtSecret)
	userHandler := handlers.NewUserHandler(userService)

	r := gin.Default()

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

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
