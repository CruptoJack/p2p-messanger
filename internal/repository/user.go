package repository

import (
	"context"
	"errors"

	"github.com/p2p-messanger/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByLogin(ctx context.Context, login string) (*models.User, error)
}
