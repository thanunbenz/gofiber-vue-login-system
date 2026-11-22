package handler

import (
	"errors"
	"strings"
	"time"

	"github.com/thanunbenz/login_with_go_vue/backend/internal/repository"
	"github.com/thanunbenz/login_with_go_vue/backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	svc          *service.AuthService
	cookieSecure bool
}

func NewAuthHandler(svc *service.AuthService, cookieSecure bool) *AuthHandler {
	return &AuthHandler{svc: svc, cookieSecure: cookieSecure}
}

type registerReq struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req registerReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	req.Username = strings.TrimSpace(req.Username)

	if len(req.Username) < 4 || len(req.Username) > 20 {
		return fiber.NewError(fiber.StatusBadRequest, "username must be 4-20 characters")
	}
	if len(req.Password) < 8 {
		return fiber.NewError(fiber.StatusBadRequest, "password must be at least 8 characters")
	}
	if req.Password != req.ConfirmPassword {
		return fiber.NewError(fiber.StatusBadRequest, "passwords do not match")
	}

	if err := h.svc.Register(c.UserContext(), req.Username, req.Password); err != nil {
		if errors.Is(err, repository.ErrUser) {
			return fiber.NewError(fiber.StatusConflict, "username already exists")
		}
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "registered"})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	accessToken, refreshToken, err := h.svc.Login(c.UserContext(), strings.TrimSpace(req.Username), req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserPassword) {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		return err
	}
	h.setRefreshCookie(c, refreshToken, time.Now().Add(12*time.Hour))
	return c.JSON(fiber.Map{"token": accessToken})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	accessToken, refreshToken, err := h.svc.Refresh(c.UserContext(), c.Cookies(h.refreshCookieName()))
	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid refresh token")
		}
		return err
	}
	h.setRefreshCookie(c, refreshToken, time.Now().Add(12*time.Hour))
	return c.JSON(fiber.Map{"token": accessToken})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	if err := h.svc.Logout(c.UserContext(), c.Cookies(h.refreshCookieName())); err != nil {
		return err
	}
	h.setRefreshCookie(c, "", time.Unix(0, 0))
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AuthHandler) setRefreshCookie(c *fiber.Ctx, value string, expires time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     h.refreshCookieName(),
		Value:    value,
		Path:     "/",
		Expires:  expires,
		HTTPOnly: true,
		Secure:   h.cookieSecure,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
}

func (h *AuthHandler) refreshCookieName() string {
	if h.cookieSecure {
		return "__Host-refresh_token"
	}
	return "refresh_token"
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"username": c.Locals("username")})
}
