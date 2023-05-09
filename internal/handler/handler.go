package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type ITasks interface {
	CreateTask(ctx context.Context) error
}
type Handler struct {
	ts ITasks
}

func NewHandler(tasksService ITasks) *Handler {
	return &Handler{
		ts: tasksService,
	}
}

func (h *Handler) CreateTask(c *fiber.Ctx) error {
	return fiber.NewError(fiber.StatusNotImplemented, "not implemented")
}
