package user

import (
	"bookLog/internal/models"

	"github.com/google/uuid"
)

type UserService interface {
	GetAllUsers(search string) ([]models.User, int, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uuid.UUID) error
}
