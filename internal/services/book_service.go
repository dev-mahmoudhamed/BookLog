package service

import "bookLog/internal/models"

type BookService interface {
	GetAllBooks() ([]models.Book, error)
	GetBookByID(id int) (*models.Book, error)
	CreateBook(book *models.Book) error
	UpdateBook(book *models.Book) error
	DeleteBook(id int) error
}
