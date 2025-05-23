package main

import (
	"errors"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	if port == "" {
		port = "3000"
	}

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()
	app := fiber.New()

	app.Use(middlewareLogger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("CLIENT_URL"),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	repo, err := NewPostgreStore()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	sugar.Info("Connected to database")
	validate := validator.New()
	handler := NewHandlerManager(repo, sugar, validate)

	app.Get("/", Welcome)
	app.Get("/healthcheck", HealthCheck)
	app.Get("/classrooms/most-visited", handler.GetMostVisitedClassRoom)
	app.Get("/classrooms/:code", handler.GetClassRoom)
	app.Post("/classrooms", AuthMiddleware, handler.AddClassRoom)
	app.Post("/login", LoginAdmin)
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{"error": "The route you are looking for does not exist"})
	})

	log.Fatal(app.Listen(":" + port))
}
