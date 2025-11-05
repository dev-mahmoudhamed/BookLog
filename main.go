package main

import (
	"bookLog/database"
	"bookLog/handlers"
	"bookLog/internal/config"
	"bookLog/internal/repository/postgres"
	service "bookLog/internal/services"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
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

	bookRepo := postgres.NewBookRepositoryPostgres(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handlers.NewBookHandler(bookService)

	r := gin.Default()

	// CRUD routes
	r.GET("/books", bookHandler.GetBooks)
	r.GET("/books/:id", bookHandler.GetBookByID)
	r.POST("/books", bookHandler.CreateBook)
	r.PUT("/books/:id", bookHandler.UpdateBook)
	r.DELETE("/books/:id", bookHandler.DeleteBook)

	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
