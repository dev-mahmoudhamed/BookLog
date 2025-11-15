package user

import (
	"bookLog/internal/models"
	"bookLog/internal/repository"
	"bookLog/internal/util"
	"errors"

	"github.com/google/uuid"
)

type AuthService struct {
	repo   repository.UserRepository
	secret string
}

func NewAuthService(repo repository.UserRepository, secret string) *AuthService {
	return &AuthService{repo: repo, secret: secret}
}

func (s *AuthService) Register(fullname, email, password string) (*models.User, error) {
	// check existing user by email
	if _, err := s.repo.GetByEmail(email); err == nil {
		return nil, errors.New("user already exists")
	}

	hashed, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}

	u := &models.User{
		ID:       uuid.New(),
		FullName: fullname,
		Email:    email,
		Password: hashed,
		Role:     "user",
	}

	if err := s.repo.Create(u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !util.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := util.GenerateJWT(user.ID, s.secret)
	if err != nil {
		return "", err
	}

	return token, nil
}
