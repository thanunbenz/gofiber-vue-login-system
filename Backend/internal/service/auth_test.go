package service

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/thanunbenz/login_with_go_vue/backend/internal/model"
	"github.com/thanunbenz/login_with_go_vue/backend/internal/repository"
)

type fakeUserRepository struct{}

func (fakeUserRepository) CreateUser(context.Context, string, string) error {
	return nil
}

func (fakeUserRepository) GetUserByUsername(context.Context, string) (*model.User, error) {
	return nil, repository.ErrNotFound
}

type fakeRefreshTokenRepository struct {
	oldHash []byte
	err     error
}

func (f *fakeRefreshTokenRepository) Create(context.Context, int, []byte, time.Time) error {
	return nil
}

func (f *fakeRefreshTokenRepository) Rotate(_ context.Context, oldHash, _ []byte, _ time.Time) (int, string, error) {
	f.oldHash = oldHash
	return 1, "thanunbenz", f.err
}

func (f *fakeRefreshTokenRepository) Delete(context.Context, []byte) error {
	return nil
}

func TestRefreshReturnsNewTokens(t *testing.T) {
	repo := &fakeRefreshTokenRepository{}
	auth := NewAuthService(fakeUserRepository{}, repo, "secret")

	accessToken, refreshToken, err := auth.Refresh(context.Background(), "old-token")
	if err != nil {
		t.Fatal(err)
	}
	if accessToken == "" {
		t.Fatal("access token is empty")
	}
	if refreshToken == "" || refreshToken == "old-token" {
		t.Fatal("refresh token was not rotated")
	}
	if len(repo.oldHash) != 32 || bytes.Equal(repo.oldHash, []byte("old-token")) {
		t.Fatal("refresh token was not hashed")
	}

	parsed, err := jwt.Parse(accessToken, func(*jwt.Token) (any, error) {
		return []byte("secret"), nil
	})
	if err != nil || !parsed.Valid {
		t.Fatal("access token is invalid")
	}
}

func TestRefreshRejectsInvalidToken(t *testing.T) {
	repo := &fakeRefreshTokenRepository{err: repository.ErrInvalidRefreshToken}
	auth := NewAuthService(fakeUserRepository{}, repo, "secret")

	_, _, err := auth.Refresh(context.Background(), "wrong-token")
	if !errors.Is(err, ErrInvalidRefreshToken) {
		t.Fatalf("got %v, want %v", err, ErrInvalidRefreshToken)
	}
}
