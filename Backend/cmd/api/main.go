package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/thanunbenz/login_with_go_vue/backend/internal/config"
	"github.com/thanunbenz/login_with_go_vue/backend/internal/database"
	"github.com/thanunbenz/login_with_go_vue/backend/internal/handler"
	"github.com/thanunbenz/login_with_go_vue/backend/internal/middleware"
	"github.com/thanunbenz/login_with_go_vue/backend/internal/repository"
	"github.com/thanunbenz/login_with_go_vue/backend/internal/service"
)

func main() {
	cfg := config.Load()
	app := fiber.New(fiber.Config{
		ErrorHandler: handler.ErrorHandler,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.FrontendURL,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))
	ctx := context.Background()
	pool, err := database.Connect(ctx, cfg.DBURL)
	if err != nil {
		log.Fatal("db connect: ", err)
	}
	defer pool.Close()

	users := repository.NewUserRepo(pool)
	refreshTokens := repository.NewRefreshTokenRepo(pool)
	svc := service.NewAuthService(users, refreshTokens, cfg.JWT_SECRET)
	h := handler.NewAuthHandler(svc, cfg.CookieSecure)

	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	v1.Post("/register", h.Register)
	v1.Post("/login", h.Login)
	v1.Post("/refresh", h.Refresh)
	v1.Post("/logout", h.Logout)
	v1.Get("/me", middleware.JWT(cfg.JWT_SECRET), h.Me)

	app.Listen(":" + cfg.Port)
}
