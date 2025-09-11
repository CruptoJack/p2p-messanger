package repository

import (
	"context"

	"github.com/p2p-messanger/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByLogin(ctx context.Context, login string) (*models.User, error)
}
