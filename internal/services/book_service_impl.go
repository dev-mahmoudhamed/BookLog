package service

import (
	"bookLog/internal/models"
	"bookLog/internal/repository"
	"fmt"
)

type bookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) BookService {
	return &bookService{repo: repo}
}

func (s *bookService) GetAllBooks() ([]models.Book, error) {
	return s.repo.GetAll()
}

func (s *bookService) GetBookByID(id int) (*models.Book, error) {
	return s.repo.GetByID(id)
}

func (s *bookService) CreateBook(book *models.Book) error {
	// Business rule: title must not be empty
	if book.Title == "" {
		return fmt.Errorf("book title cannot be empty")
	}
	return s.repo.Create(book)
}

func (s *bookService) UpdateBook(book *models.Book) error {
	return s.repo.Update(book)
}

func (s *bookService) DeleteBook(id int) error {
	return s.repo.Delete(id)
}
