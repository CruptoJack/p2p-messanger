package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/p2p-messanger/internal/models"
	"github.com/p2p-messanger/internal/repository"
	"github.com/p2p-messanger/pkg/config"
)

type AuthService struct {
	userRepo repository.UserRepository
	jwtCfg   config.JWT
}

func NewAuthService(userRepo repository.UserRepository, jwtCfg config.JWT) *AuthService {
	return &AuthService{userRepo: userRepo, jwtCfg: jwtCfg}
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

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password:%w", err)
	}

	return user, nil
}

func (as *AuthService) GenerateToken(user *models.User) (string, error) {

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"login":   user.Login,
		"exp":     time.Now().Add(as.jwtCfg.Expire).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(as.jwtCfg.SecretKey))
}

func (as *AuthService) ValidateToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(as.jwtCfg.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
