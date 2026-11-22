package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInvalidRefreshToken = errors.New("invalid refresh token")

type RefreshTokenRepo struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepo(db *pgxpool.Pool) *RefreshTokenRepo {
	return &RefreshTokenRepo{db: db}
}

func (r *RefreshTokenRepo) Create(ctx context.Context, userID int, hash []byte, expiresAt time.Time) error {
	if _, err := r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE expires_at <= now()`); err != nil {
		return err
	}
	_, err := r.db.Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, hash, expiresAt)
	return err
}

func (r *RefreshTokenRepo) Rotate(ctx context.Context, oldHash, newHash []byte, expiresAt time.Time) (int, string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, "", err
	}
	defer tx.Rollback(ctx)

	var userID int
	var username string
	err = tx.QueryRow(ctx, `
		SELECT rt.user_id, u.username
		FROM refresh_tokens rt
		JOIN users u ON u.id = rt.user_id
		WHERE rt.token_hash = $1 AND rt.expires_at > now()
		FOR UPDATE OF rt`, oldHash).Scan(&userID, &username)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, "", ErrInvalidRefreshToken
	}
	if err != nil {
		return 0, "", err
	}

	if _, err = tx.Exec(ctx, `DELETE FROM refresh_tokens WHERE token_hash = $1`, oldHash); err != nil {
		return 0, "", err
	}
	if _, err = tx.Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, newHash, expiresAt); err != nil {
		return 0, "", err
	}

	return userID, username, tx.Commit(ctx)
}

func (r *RefreshTokenRepo) Delete(ctx context.Context, hash []byte) error {
	_, err := r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE token_hash = $1`, hash)
	return err
}
