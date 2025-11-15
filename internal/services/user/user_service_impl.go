package user

import (
	"bookLog/internal/models"
	"bookLog/internal/repository"

	"github.com/google/uuid"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// GetAllUsers implements UserService.
func (u *userService) GetAllUsers(search string) ([]models.User, int, error) {
	return u.repo.GetAll(search)
}

// DeleteUser implements UserService.
func (u *userService) DeleteUser(id uuid.UUID) error {
	return u.repo.Delete(id)
}

// GetUserByID implements UserService.
func (u *userService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return u.repo.GetByID(id)
}

// UpdateUser implements UserService.
func (u *userService) UpdateUser(user *models.User) error {
	return u.repo.Update(user)
}
