package handler

import (
	"context"
	"task-manager-api/internal/comment"
	"task-manager-api/internal/profile"
	"task-manager-api/internal/taskmanager"

	"github.com/gofiber/fiber/v2"
)

//go:generate mockgen -source=./handler.go -destination=./mock/handler_mock.go
type ITasks interface {
	CreateTask(ctx context.Context, topic string, desc string, ownerId string) (*taskmanager.TaskDoc, error)
}
type IComments interface {
	CreateComment(ctx context.Context, ownerId string, topicID string, content string) error
	GetComment(ctx context.Context, id string) (*comment.CommentDoc, error)
}

type IProfile interface {
	GetProfile(ctx context.Context, ownerId string) (*profile.ProfileDoc, error)
}
type Handler struct {
	task    ITasks
	comment IComments
	profile IProfile
}

func NewHandler(tasksService ITasks) *Handler {
	return &Handler{}
}

func (h *Handler) CreateTask(c *fiber.Ctx) error {
	return fiber.NewError(fiber.StatusNotImplemented, "not implemented")
}
