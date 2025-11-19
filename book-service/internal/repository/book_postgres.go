package repository

import (
	"book-service/internal/models"

	"github.com/google/uuid"
)

func (r *bookGorm) Create(book models.Book) (uint, error) {
	result := r.db.Create(&book)
	return book.ID, result.Error
}

func (r *bookGorm) GetAll(userID uuid.UUID) ([]models.Book, error) {
	var books []models.Book
	result := r.db.Where("user_id = ?", userID).Find(&books)
	return books, result.Error
}

func (r *bookGorm) GetByID(id uint) (models.Book, error) {
	var book models.Book
	result := r.db.First(&book, id)
	return book, result.Error
}

func (r *bookGorm) Update(id uint, book models.Book) error {
	result := r.db.Model(&models.Book{}).Where("id = ?", id).Updates(book)
	return result.Error
}

func (r *bookGorm) Delete(id uint) error {
	result := r.db.Delete(&models.Book{}, id)
	return result.Error
}
