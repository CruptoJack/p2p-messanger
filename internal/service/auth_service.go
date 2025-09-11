package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/p2p-messanger/internal/models"
	"github.com/p2p-messanger/internal/repository"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (as *AuthService) RegistUser(login, password string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashed password:%w", err)
	}

	user := &models.User{
		Login:    login,
		Password: string(hashedPassword),
	}

	ctx := context.Background()
	err = as.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error of create user in psql:%w", err)
	}

	return user, nil
}

func (s *AuthService) LoginUser(login, password string) (*models.User, error) {
	ctx := context.Background()
	user, err := s.userRepo.FindByLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("user not found:%w", err)
	}
	return user, nil
}
