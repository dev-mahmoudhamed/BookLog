package services

import (
	"book-service/internal/models"
	"book-service/internal/repository"

	"github.com/google/uuid"
)

type BookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) *BookService {
	return &BookService{repo}
}

func (s *BookService) CreateBook(book models.Book) (uint, error) {
	return s.repo.Create(book)
}

func (s *BookService) GetBooks(userID uuid.UUID) ([]models.Book, error) {
	return s.repo.GetAll(userID)
}

func (s *BookService) GetBook(id uint) (models.Book, error) {
	return s.repo.GetByID(id)
}

func (s *BookService) UpdateBook(id uint, book models.Book) error {
	return s.repo.Update(id, book)
}

func (s *BookService) DeleteBook(id uint) error {
	return s.repo.Delete(id)
}
