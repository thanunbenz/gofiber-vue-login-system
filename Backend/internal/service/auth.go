package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/thanunbenz/login_with_go_vue/backend/internal/model"
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
	users         userRepository
	refreshTokens refreshTokenRepository
	secret        []byte
}

type userRepository interface {
	CreateUser(context.Context, string, string) error
	GetUserByUsername(context.Context, string) (*model.User, error)
}

type refreshTokenRepository interface {
	Create(context.Context, int, []byte, time.Time) error
	Rotate(context.Context, []byte, []byte, time.Time) (int, string, error)
	Delete(context.Context, []byte) error
}

func NewAuthService(users userRepository, refreshTokens refreshTokenRepository, secret string) *AuthService {
	return &AuthService{users: users, refreshTokens: refreshTokens, secret: []byte(secret)}
}

func (s *AuthService) Register(ctx context.Context, username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.users.CreateUser(ctx, username, string(hash))
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, string, error) {
	user, err := s.users.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", "", ErrUserPassword
		}
		return "", "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", ErrUserPassword
	}

	accessToken, err := s.newToken(user.ID, user.Username, accessTokenLife)
	if err != nil {
		return "", "", err
	}
	refreshToken, refreshHash, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}
	if err := s.refreshTokens.Create(ctx, user.ID, refreshHash, time.Now().Add(refreshTokenLife)); err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	newToken, newHash, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}
	userID, username, err := s.refreshTokens.Rotate(
		ctx,
		refreshTokenHash(refreshToken),
		newHash,
		time.Now().Add(refreshTokenLife),
	)
	if errors.Is(err, repository.ErrInvalidRefreshToken) {
		return "", "", ErrInvalidRefreshToken
	}
	if err != nil {
		return "", "", err
	}

	accessToken, err := s.newToken(userID, username, accessTokenLife)
	return accessToken, newToken, err
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	return s.refreshTokens.Delete(ctx, refreshTokenHash(refreshToken))
}

func (s *AuthService) newToken(userID int, username string, lifetime time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"uid":      userID,
		"username": username,
		"type":     "access",
		"exp":      time.Now().Add(lifetime).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

func newRefreshToken() (string, []byte, error) {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return "", nil, err
	}
	token := base64.RawURLEncoding.EncodeToString(value)
	return token, refreshTokenHash(token), nil
}

func refreshTokenHash(token string) []byte {
	hash := sha256.Sum256([]byte(token))
	return hash[:]
}
