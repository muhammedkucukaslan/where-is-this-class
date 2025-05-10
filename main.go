package main

import (
	"errors"
	"log"

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

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	app := fiber.New()

	app.Use(middlewareLogger.New())

	repo, err := NewPostgreStore()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	handler := NewHandlerManager(repo, sugar)

	app.Get("/:code", handler.GetClassRoom)
	app.Post("/", handler.AddClassRoom)
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{})
	})

	log.Fatal(app.Listen(":3000"))
}
