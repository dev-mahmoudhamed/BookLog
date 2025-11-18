package services

import (
	"authService/internal/models"
	"authService/internal/repository"
	"authService/util"
	"errors"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	repo   repository.UserRepository
	secret string
}

func NewAuthService(repo repository.UserRepository, secret string) *AuthService {
	return &AuthService{repo: repo, secret: secret}
}

func (s *AuthService) Register(fullname, email, password string) error {
	// check existing user by email
	if _, err := s.repo.GetByEmail(email); err == nil {
		return errors.New("user already exists")
	}

	hashed, err := util.HashPassword(password)
	if err != nil {
		return err
	}

	u := &models.User{
		ID:       uuid.New(),
		FullName: fullname,
		Email:    email,
		Password: hashed,
		Role:     "user",
	}

	if err := s.repo.Create(u); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(email, password string) (string, time.Time, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", time.Time{}, errors.New("invalid credentials")
	}

	if !util.CheckPasswordHash(password, user.Password) {
		return "", time.Time{}, errors.New("invalid credentials")
	}

	token, exp, err := util.GenerateJWT(user.ID, s.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return token, exp, err
}
