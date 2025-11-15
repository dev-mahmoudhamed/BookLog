package main

import (
	"bookLog/database"
	"bookLog/handlers"
	"bookLog/internal/config"
	"bookLog/internal/middleware"
	"bookLog/internal/repository/postgres"
	"log"

	bookService "bookLog/internal/services/book"
	authService "bookLog/internal/services/user"
	userService "bookLog/internal/services/user"

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
	bookService := bookService.NewBookService(bookRepo)
	bookHandler := handlers.NewBookHandler(bookService)

	userRepo := postgres.NewUserRepositoryPostgres(db)
	userService := userService.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService, authService.NewAuthService(userRepo, cfg.JwtSecret))

	r := gin.Default()

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware(cfg.JwtSecret))

	// CRUD routes
	auth.GET("/books", bookHandler.GetBooks)
	auth.GET("/books/:id", bookHandler.GetBookByID)
	auth.POST("/books", bookHandler.CreateBook)
	auth.PUT("/books/:id", bookHandler.UpdateBook)
	auth.DELETE("/books/:id", bookHandler.DeleteBook)

	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
