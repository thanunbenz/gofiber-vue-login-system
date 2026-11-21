package main

import (
	"context"
	"log"
	"os"

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
	app.Use(cors.New())
	ctx := context.Background()
	pool, err := database.Connect(ctx, cfg.DBURL)
	if err != nil {
		log.Fatal("db connect: ", err)
	}
	defer pool.Close()

	repo := repository.NewUserRepo(pool)
	svc := service.NewAuthService(repo, cfg.JWT_SECRET)
	h := handler.NewAuthHandler(svc)

	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	v1.Post("/register", h.Register)
	v1.Post("/login", h.Login)
	v1.Get("/me", middleware.JWT(cfg.JWT_SECRET), h.Me)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Listen(":" + port)
}
