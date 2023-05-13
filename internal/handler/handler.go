package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"task-manager-api/config"
	"task-manager-api/internal/comment"
	"task-manager-api/internal/profile"
	"task-manager-api/internal/taskmanager"

	"github.com/gofiber/fiber/v2"
)

type response struct {
	Data interface{} `json:"data"`
}

//go:generate mockgen -source=./handler.go -destination=./mock/handler_mock.go
type ITasks interface {
	CreateTask(ctx context.Context, ownerId string, topic string, desc string) (*taskmanager.TaskDoc, error)
	GetAllTask(ctx context.Context, page int, limit int) ([]taskmanager.TaskDoc, error)
	ArchiveTask(ctx context.Context, ownerId string, id string) (int, error)
	UpdateTaskStatus(ctx context.Context, ownerId string, id string, status int) error
	GetTask(ctx context.Context, id string) (*taskmanager.TaskDoc, error)
}
type IComments interface {
	CreateComment(ctx context.Context, ownerId string, taskId string, content string) (*comment.CommentDoc, error)
	GetTopicComments(ctx context.Context, taskId string, page int, limit int) ([]comment.CommentDoc, error)
}

type IProfile interface {
	GetProfile(ctx context.Context, ownerId string) (*profile.ProfileDoc, error)
	GetProfileList(ctx context.Context, ownerId []string) ([]profile.ProfileDoc, error)
}
type Handler struct {
	task    ITasks
	comment IComments
	profile IProfile
}

func NewHandler(tasksService ITasks, commentService IComments, profileService IProfile) *Handler {
	return &Handler{
		task:    tasksService,
		comment: commentService,
		profile: profileService,
	}
}

func (h *Handler) CreateTask(c *fiber.Ctx) error {
	payload := struct {
		Topic       string `json:"topic"`
		Description string `json:"description"`
	}{}

	ownerId := c.Params("ownerId")
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	if err := h.validateOwnerId(c, ownerId); err != nil {
		return err
	}
	topic := strings.TrimSpace(payload.Topic)
	description := strings.TrimSpace(payload.Description)
	if topic == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Topic is required")
	}

	if description == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Description is required")
	}

	task, err := h.task.CreateTask(c.Context(), ownerId, topic, description)
	if err != nil {
		// TODO: convert interl sensitive error to code [1213413]
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(http.StatusCreated).JSON(response{
		Data: task,
	})
}

func (h *Handler) GetAllTask(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		fmt.Println("Error during conversion")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid page number")
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		fmt.Println("Error during conversion 'limit'")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid limit number")
	}

	if limitInt > config.Conf.Pagination.MaxLimit {
		return fiber.NewError(fiber.StatusBadRequest, "Limit cannot be more than 100")
	}

	tasks, err := h.task.GetAllTask(c.Context(), pageInt, limitInt)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(response{
		Data: tasks,
	})
}

func (h *Handler) GetTask(c *fiber.Ctx) error {
	taskId := c.Params("taskId")
	task, err := h.task.GetTask(c.Context(), taskId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(response{
		Data: task,
	})
}

func (h *Handler) ArchiveTask(c *fiber.Ctx) error {
	taskId := c.Params("taskId")
	ownerId := c.Params("ownerId")
	modifiedCount, err := h.task.ArchiveTask(c.Context(), ownerId, taskId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if modifiedCount == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Task or account not found")
	}
	return c.JSON(response{
		Data: "Task archived successfully",
	})
}

func (h *Handler) UpdateTask(c *fiber.Ctx) error {
	taskId := c.Params("taskId")
	ownerId := c.Params("ownerId")
	payload := struct {
		Status *int `json:"status"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	// currently only status can be updated
	if payload.Status != nil {
		if *payload.Status < taskmanager.TaskStatusOpen || *payload.Status > taskmanager.TaskStatusDone {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid status")
		}
		err := h.task.UpdateTaskStatus(c.Context(), ownerId, taskId, *payload.Status)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(response{
			Data: "Task status updated successfully",
		})
	}

	return fiber.NewError(fiber.StatusNotImplemented, "Not implemented")
}

func (h *Handler) GetTopicComments(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		fmt.Println("Error during conversion")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid page number")
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		fmt.Println("Error during conversion 'limit'")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid limit number")
	}

	if limitInt > config.Conf.Pagination.MaxLimit {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Limit cannot be more than %v", config.Conf.Pagination.MaxLimit))
	}

	taskId := c.Params("taskId")
	comments, err := h.comment.GetTopicComments(c.Context(), taskId, pageInt, limitInt)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(response{
		Data: comments,
	})
}

func (h *Handler) CreateComment(c *fiber.Ctx) error {
	payload := struct {
		Content string `json:"content"`
	}{}

	ownerId := c.Params("ownerId")
	taskId := c.Params("taskId")
	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	if err := h.validateOwnerId(c, ownerId); err != nil {
		return err
	}

	content := strings.TrimSpace(payload.Content)
	if content == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Content is required")
	}

	comment, err := h.comment.CreateComment(c.Context(), ownerId, taskId, content)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(http.StatusCreated).JSON(response{
		Data: comment,
	})
}

func (h *Handler) GetProfile(c *fiber.Ctx) error {
	ownerId := c.Params("ownerId")
	profile, err := h.profile.GetProfile(c.Context(), ownerId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if profile == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid owner id")
	}

	return c.JSON(response{
		Data: profile,
	})
}

func (h *Handler) GetProfileList(c *fiber.Ctx) error {
	ownerId := c.Query("owner_id", "")
	ownerIds := strings.Split(ownerId, ",")
	if len(ownerIds) == 0 || ownerId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid owner id")
	}
	if len(ownerIds) > config.Conf.Pagination.MaxGetProfileLimit {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Limit cannot be more than %v", config.Conf.Pagination.MaxGetProfileLimit))
	}
	profiles, err := h.profile.GetProfileList(c.Context(), ownerIds)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response{
		Data: profiles,
	})
}

func (h *Handler) validateOwnerId(c *fiber.Ctx, ownerId string) error {
	if ownerId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid owner id")
	}
	profile, err := h.profile.GetProfile(c.Context(), ownerId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if profile == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid owner id")
	}

	return nil
}
