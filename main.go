package main

import (
	"errors"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	middlewareLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/subosito/gotenv"
	"go.uber.org/zap"
)

var (
	ErrClassRoomNotFound      = errors.New("classroom not found")
	ErrClassRoomAlreadyExists = errors.New("classroom already exists")
)

func main() {
	gotenv.Load()

	port := os.Getenv("PORT")

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	app := fiber.New()

	app.Use(middlewareLogger.New())

	repo, err := NewPostgreStore()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	sugar.Info("Connected to database")
	handler := NewHandlerManager(repo, sugar)

	app.Get("/", Welcome)
	app.Get("/healthcheck", HealthCheck)
	app.Get("/codes/:code", handler.GetClassRoom)
	app.Post("/codes", handler.AddClassRoom)
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{})
	})

	log.Fatal(app.Listen(":" + port))
}
