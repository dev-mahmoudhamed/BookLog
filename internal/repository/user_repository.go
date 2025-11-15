package repository

import (
	"bookLog/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetAll(search string) ([]models.User, int, error)
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uuid.UUID) error
}
