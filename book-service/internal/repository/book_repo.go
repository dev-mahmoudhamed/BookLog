package repository

import (
	"book-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookRepository interface {
	Create(book models.Book) (uint, error)
	GetAll(userID uuid.UUID) ([]models.Book, error)
	GetByID(id uint) (models.Book, error)
	Update(id uint, book models.Book) error
	Delete(id uint) error
}

type bookGorm struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookGorm{db: db}
}
