package handler

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "internal server error"

	var fe *fiber.Error
	if errors.As(err, &fe) {
		code, msg = fe.Code, fe.Message
	} else {
		log.Println("error:", err)
	}
	return c.Status(code).JSON(fiber.Map{"error": msg})
}
