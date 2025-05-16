package main

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type HandlerManager struct {
	repo     Repo
	logger   *zap.SugaredLogger
	validate *validator.Validate
}

type Repo interface {
	GetClassRoom(code string, language string) (GetClassRoomResponse, error)
	CreateClassRoom(*AddClassRoomRequest) error
	GetMostVisitedClassRoom() (GetMostVisitedClassRoomsResponse, error)
}

func NewHandlerManager(repo Repo, logger *zap.SugaredLogger, validate *validator.Validate) *HandlerManager {
	return &HandlerManager{
		repo:     repo,
		logger:   logger,
		validate: validate,
	}
}

// Healtcheck Handler
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("OK")
}

// Welcome Handler
func Welcome(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Welcome to the Classroom API. Go `https://github.com/muhammedkucukaslan/where-is-this-class` to check the codebase."})
}

type LoginAdminRequest struct {
	Password string `json:"password" validate:"required"`
}

// Login Admin Handler
func LoginAdmin(c *fiber.Ctx) error {
	var req LoginAdminRequest
	if err := c.BodyParser(&req); err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if req.Password != os.Getenv("ADMIN_PASSWORD") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid password"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "strict",
	})
	return c.SendStatus(fiber.StatusNoContent)
}

// Get Class Room Handler
type GetClassRoomResponse struct {
	Building    string `json:"building"`
	Floor       int    `json:"floor"`
	ImageUrl    string `json:"imageUrl"`
	Description string `json:"description"`
	Detail      string `json:"detail"`
}

func (h *HandlerManager) GetClassRoom(c *fiber.Ctx) error {
	code := c.Params("code")
	language := c.Query("language")

	if language == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "You must specify 'language' in the query"})
	}

	classroom, err := h.repo.GetClassRoom(strings.ToUpper(strings.ReplaceAll(code, "%20", "")), language)
	if err != nil {
		if errors.Is(err, ErrClassRoomNotFound) {
			h.logger.Errorw("Classroom not found", "code", code)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Classroom not found"})
		}
		h.logger.Errorw("Error retrieving classroom", "code", code, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
	return c.Status(fiber.StatusOK).JSON(classroom)
}

// Add Class Room Handler
type AddClassRoomRequest struct {
	Code         string        `json:"code" validate:"required,min=1,max=10"`
	Floor        *int          `json:"floor" validate:"required"`
	ImageUrl     string        `json:"imageUrl" validate:"required,url"`
	Translations []Translation `json:"translations" validate:"required,dive,required"`
}

type Translation struct {
	Language    string `json:"language" validate:"required,oneof=en tr ar"`
	Building    string `json:"building" validate:"required"`
	Description string `json:"description" validate:"required"`
	Detail      string `json:"detail"`
}

func (h *HandlerManager) AddClassRoom(c *fiber.Ctx) error {
	var req AddClassRoomRequest

	if err := c.BodyParser(&req); err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) {
		h.logger.Errorw("Error parsing request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.validate.Struct(req)

	if err != nil {
		h.logger.Errorw("Validation error", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = h.repo.CreateClassRoom(&req)
	if err != nil {
		if errors.Is(err, ErrClassRoomAlreadyExists) {
			h.logger.Errorw("Classroom already exists", "code", req.Code)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Classroom already exists"})
		}
		h.logger.Errorw("Error creating classroom", "code", req.Code, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Get Most Visited Class Room Handler

type GetMostVisitedClassRoomsResponse []MostVisitedClassRoom

type MostVisitedClassRoom struct {
	Code    string `json:"code"`
	Visited int    `json:"visited"`
}

func (h *HandlerManager) GetMostVisitedClassRoom(c *fiber.Ctx) error {
	mostVisitedClassRooms, err := h.repo.GetMostVisitedClassRoom()
	if err != nil {
		h.logger.Errorw("Error retrieving most visited classrooms", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
	return c.Status(fiber.StatusOK).JSON(mostVisitedClassRooms)
}
