package handler

import (
	"errors"
	"strings"

	"github.com/thanunbenz/login_with_go_vue/backend/internal/repository"
	"github.com/thanunbenz/login_with_go_vue/backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
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

	token, err := h.svc.Login(c.UserContext(), strings.TrimSpace(req.Username), req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserPassword) {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		return err
	}
	return c.JSON(fiber.Map{"token": token})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"username": c.Locals("username")})
}
