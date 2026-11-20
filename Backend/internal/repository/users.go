package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thanunbenz/login_with_go_vue/backend/internal/model"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

var ErrUser = errors.New("username already exists")
var ErrNotFound = errors.New("username not found")

// create a new user
func (r *UserRepo) CreateUser(ctx context.Context, username, password string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (username, password_hash) VALUES ($1, $2)`,
		username, password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "pgerrcode.UniqueViolation" {
			return ErrUser
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx,
		`SELECT id, username, password_hash, created_at FROM users WHERE username = $1`,
		username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	return &user, err
}
