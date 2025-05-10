package main

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type HandlerManager struct {
	repo   Repo
	logger *zap.SugaredLogger
}

type Repo interface {
	GetClassRoom(code string) (GetClassRoomResponse, error)
	CreateClassRoom(*AddClassRoomRequest) error
}

func NewHandlerManager(repo Repo, logger *zap.SugaredLogger) *HandlerManager {
	return &HandlerManager{
		repo:   repo,
		logger: logger,
	}
}

// Healtcheck Handler
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

// Welcome Handler
func Welcome(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Welcome to the Classroom API. Go `https://github.com/muhammedkucukaslan/where-is-this-class` to check the codebase."})
}

// Get Class Room Handler
type GetClassRoomResponse struct {
	Code       string `json:"code"`
	Building   string `json:"building"`
	Floor      int    `json:"floor"`
	FloorName  string `json:"floorName"`
	Directions string `json:"directions"`
}

func (h *HandlerManager) GetClassRoom(c *fiber.Ctx) error {
	code := c.Params("code")

	classroom, err := h.repo.GetClassRoom(code)
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

type AddClassRoomRequest struct {
	Code       string `json:"code" validate:"required"`
	Building   string `json:"building" validate:"required"`
	Floor      int    `json:"floor" validate:"required"`
	FloorName  string `json:"floorName" validate:"required"`
	Directions string `json:"directions" validate:"required"`
}

func (h *HandlerManager) AddClassRoom(c *fiber.Ctx) error {

	var req AddClassRoomRequest
	if err := c.BodyParser(req); err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) {
		h.logger.Errorw("Error parsing request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.repo.CreateClassRoom(&req)
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
