package handler

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestRefreshCookieSecurity(t *testing.T) {
	h := &AuthHandler{cookieSecure: true}
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		h.setRefreshCookie(c, "token", time.Now().Add(time.Hour))
		return c.SendStatus(fiber.StatusNoContent)
	})

	response, err := app.Test(httptest.NewRequest("GET", "/", nil))
	if err != nil {
		t.Fatal(err)
	}
	cookie := strings.ToLower(response.Header.Get("Set-Cookie"))
	for _, attribute := range []string{"__host-refresh_token=", "httponly", "secure", "samesite=lax", "path=/"} {
		if !strings.Contains(cookie, attribute) {
			t.Fatalf("refresh cookie is missing %s: %s", attribute, cookie)
		}
	}
}
