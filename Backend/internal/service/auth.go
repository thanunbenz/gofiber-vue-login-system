package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/thanunbenz/login_with_go_vue/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserPassword = errors.New("invalid username or password")

type AuthService struct {
	repo   *repository.UserRepo
	secret []byte
}

func NewAuthService(repo *repository.UserRepo, secret string) *AuthService {
	return &AuthService{repo: repo, secret: []byte(secret)}
}

func (s *AuthService) Register(ctx context.Context, username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.CreateUser(ctx, username, string(hash))
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrUserPassword
		}
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrUserPassword
	}

	claims := jwt.MapClaims{
		"uid":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}
