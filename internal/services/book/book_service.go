package book

import "bookLog/internal/models"

type BookService interface {
	GetAllBooks(page, limit int, search string) ([]models.Book, int, error)
	GetBookByID(id int) (*models.Book, error)
	CreateBook(book *models.Book) error
	UpdateBook(book *models.Book) error
	DeleteBook(id int) error
}
