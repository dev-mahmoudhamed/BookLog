package repository

import "bookLog/internal/models"

type BookRepository interface {
	GetAll(limit, offset int, search string) ([]models.Book, int, error)
	GetByID(id int) (*models.Book, error)
	Create(book *models.Book) error
	Update(book *models.Book) error
	Delete(id int) error
}
