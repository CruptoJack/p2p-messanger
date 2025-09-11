package postgresql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/p2p-messanger/internal/models"
	"github.com/p2p-messanger/internal/repository"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repository.UserRepository {
	return &UserRepo{db: db}
}

func (ur *UserRepo) Create(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users(login, password) VALUES($1,$2) RETURNING id, created_at"
	return ur.db.QueryRowxContext(ctx, query, user.Login, user.Password).
		Scan(&user.ID, &user.Created_at)

}

func (ur *UserRepo) FindByLogin(ctx context.Context, login string) (*models.User, error) {
	query := "SELECT id, login, password, created_at FROM users WHERE login=$1"
	user := &models.User{}
	err := ur.db.QueryRowxContext(ctx, query, login).Scan(&user)
	if err != nil {
		return nil, fmt.Errorf("error scan login req: %w", err)
	}
	return user, nil
}
