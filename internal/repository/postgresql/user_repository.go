package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/p2p-messanger/internal/models"
	"github.com/p2p-messanger/internal/repository"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &UserRepo{db: db}
}

func (ur *UserRepo) Create(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users(login, password) VALUES($1,$2) RETURNING id, created_at"

	return ur.db.QueryRow(ctx, query, user.Login, user.Password).
		Scan(&user.ID, &user.CreatedAt)

}

func (ur *UserRepo) FindByLogin(ctx context.Context, login string) (*models.User, error) {
	query := "SELECT id, login, password, created_at FROM users WHERE login=$1"
	user := &models.User{}
	err := ur.db.QueryRow(ctx, query, login).Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", repository.ErrUserNotFound)
		}
		return nil, fmt.Errorf("error finding user by login: %w", err)
	}
	return user, nil
}
