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
var ErrInvalidRefreshToken = errors.New("invalid refresh token")

const (
	accessTokenLife  = 20 * time.Minute
	refreshTokenLife = 12 * time.Hour
)

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

func (s *AuthService) Login(ctx context.Context, username, password string) (string, string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", "", ErrUserPassword
		}
		return "", "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", ErrUserPassword
	}

	accessToken, err := s.newToken(user.ID, user.Username, "access", accessTokenLife)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := s.newToken(user.ID, user.Username, "refresh", refreshTokenLife)
	return accessToken, refreshToken, err
}

func (s *AuthService) Refresh(refreshToken string) (string, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (any, error) {
		return s.secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || !token.Valid {
		return "", ErrInvalidRefreshToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return "", ErrInvalidRefreshToken
	}
	uid, uidOK := claims["uid"].(float64)
	username, usernameOK := claims["username"].(string)
	if !uidOK || !usernameOK {
		return "", ErrInvalidRefreshToken
	}

	return s.newToken(int(uid), username, "access", accessTokenLife)
}

func (s *AuthService) newToken(userID int, username, tokenType string, lifetime time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"uid":      userID,
		"username": username,
		"type":     tokenType,
		"exp":      time.Now().Add(lifetime).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}
