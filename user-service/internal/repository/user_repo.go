package repository

import (
	"userService/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
}
