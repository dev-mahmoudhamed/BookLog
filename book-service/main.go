package main

import (
	"book-service/config"
	"book-service/database"
	"book-service/handlers"
	"book-service/internal/models"
	"book-service/internal/repository"
	"book-service/internal/services"
	"book-service/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg)

	if err := db.AutoMigrate(&models.Book{}); err != nil {
		log.Fatal("failed to migrate DB:", err)
	}
	log.Println("✅ Book service DB migrated successfully")

	bookRepo := repository.NewBookRepository(db)
	bookService := services.NewBookService(bookRepo)
	bookHandler := handlers.NewBookHandler(bookService)

	r := gin.Default()
	r.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})

	auth := r.Group("/")

	auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		auth.POST("/books", bookHandler.CreateBook)
		auth.PUT("/books/:id", bookHandler.UpdateBook)
		auth.DELETE("/books/:id", bookHandler.DeleteBook)
		auth.GET("/books", bookHandler.GetBooks)
		auth.GET("/books/:id", bookHandler.GetBook)
	}

	log.Printf("Book service running on %s", cfg.ServerAddress)
	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	r.Run(":8081")
}
