package handler

import (
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"task-manager-api/config"
	"task-manager-api/internal/comment"
	mock "task-manager-api/internal/handler/mock"
	"task-manager-api/internal/profile"
	"task-manager-api/internal/taskmanager"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	handler        *Handler
	taskService    *mock.MockITasks
	commentService *mock.MockIComments
	profileService *mock.MockIProfile
}

func (t *HandlerTestSuite) SetupTest() {
	t.ctrl = gomock.NewController(t.T())
	t.taskService = mock.NewMockITasks(t.ctrl)
	t.commentService = mock.NewMockIComments(t.ctrl)
	t.profileService = mock.NewMockIProfile(t.ctrl)
	t.handler = NewHandler(t.taskService, t.commentService, t.profileService)

	config.Conf = &config.Config{}
	config.Conf.Pagination.MaxGetProfileLimit = 3
	config.Conf.Pagination.MaxLimit = 10
}

func (t *HandlerTestSuite) TearDownTest() {
	t.ctrl.Finish()
	t.handler = nil
	t.taskService = nil
	t.commentService = nil
	t.profileService = nil
}

func TestCHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (t HandlerTestSuite) TestCreateTask() {

	t.Run("create task but invalid topic should return 400", func() {
		// Define Fiber app.
		app := fiber.New()
		// Create route with POST method for test
		app.Post("/account/:ownerId/tasks", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.CreateTask(c)
		})
		req := httptest.NewRequest("POST", "/account/1234/tasks", strings.NewReader(`{"description":"mock_desv"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(400, resp.StatusCode)
	})

	t.Run("create task but invalid description should return 400", func() {
		// Define Fiber app.
		app := fiber.New()
		// Create route with POST method for test
		app.Post("/account/:ownerId/tasks", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.CreateTask(c)
		})
		req := httptest.NewRequest("POST", "/account/1234/tasks", strings.NewReader(`{"topic":"test_topic","description":""}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(400, resp.StatusCode)
	})
	t.Run("create task but service has error should return error", func() {
		t.profileService.EXPECT().GetProfile(gomock.Any(), "1234").Return(&profile.ProfileDoc{}, nil)
		t.taskService.EXPECT().CreateTask(gomock.Any(), "1234", "test_topic", "mock_desv").Return(&taskmanager.TaskDoc{}, errors.New("create task error"))
		// Define Fiber app.
		app := fiber.New()
		// Create route with POST method for test
		app.Post("/account/:ownerId/tasks", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.CreateTask(c)
		})
		req := httptest.NewRequest("POST", "/account/1234/tasks", strings.NewReader(`{"topic":"test_topic","description":"mock_desv"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(500, resp.StatusCode)
	})

	t.Run("create task but profile service has error should return error", func() {
		t.profileService.EXPECT().GetProfile(gomock.Any(), "1234").Return(&profile.ProfileDoc{}, errors.New("get profile error"))
		// Define Fiber app.
		app := fiber.New()
		// Create route with POST method for test
		app.Post("/account/:ownerId/tasks", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.CreateTask(c)
		})
		req := httptest.NewRequest("POST", "/account/1234/tasks", strings.NewReader(`{"topic":"test_topic","description":"mock_desv"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(500, resp.StatusCode)
	})

	t.Run("create task success return task", func() {
		t.profileService.EXPECT().GetProfile(gomock.Any(), "12345").Return(&profile.ProfileDoc{}, nil)
		t.taskService.EXPECT().CreateTask(gomock.Any(), "12345", "test_topic", "mock_desv").Return(&taskmanager.TaskDoc{
			ID:          "1234",
			OwnerID:     "12345",
			Topic:       "test_topic",
			Description: "mock_desv",
			Status:      1,
			CreateDate:  2131341,
		}, nil)
		// Define Fiber app.
		app := fiber.New()
		// Create route with POST method for test
		app.Post("/account/:ownerId/tasks", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.CreateTask(c)
		})
		req := httptest.NewRequest("POST", "/account/12345/tasks", strings.NewReader(`{"topic":"test_topic","description":"mock_desv"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(201, resp.StatusCode)
		b, _ := io.ReadAll(resp.Body)
		t.Equal(`{"data":{"id":"1234","topic":"test_topic","description":"mock_desv","status":1,"create_date":2131341,"owner_id":"12345","archive_date":null,"update_date":null}}`, string(b))
	})
}

func (t HandlerTestSuite) TestGetTask() {
	t.Run("get task but service has error should return error", func() {
		t.taskService.EXPECT().GetTask(gomock.Any(), "1234").Return(&taskmanager.TaskDoc{}, errors.New("get task error"))
		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/account/:ownerId/tasks/:taskId", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetTask(c)
		})
		req := httptest.NewRequest("GET", "/account/1234/tasks/1234", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(500, resp.StatusCode)
	})

	t.Run("get task success return task", func() {
		t.taskService.EXPECT().GetTask(gomock.Any(), "1234").Return(&taskmanager.TaskDoc{
			ID:          "1234",
			OwnerID:     "12345",
			Topic:       "test_topic",
			Description: "mock_desv",
			Status:      1,
			CreateDate:  2131341,
		}, nil)
		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/account/:ownerId/tasks/:taskId", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetTask(c)
		})
		req := httptest.NewRequest("GET", "/account/1234/tasks/1234", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(200, resp.StatusCode)
		b, _ := io.ReadAll(resp.Body)
		t.Equal(`{"data":{"id":"1234","topic":"test_topic","description":"mock_desv","status":1,"create_date":2131341,"owner_id":"12345","archive_date":null,"update_date":null}}`, string(b))
	})
}

func (t HandlerTestSuite) TestUpdateTask() {
	t.Run("update task but status is invalid should return error", func() {
		// Define Fiber app.
		app := fiber.New()
		// Create route with PATCH method for test
		app.Patch("/account/:ownerId/tasks/:taskId", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.UpdateTask(c)
		})
		req := httptest.NewRequest("PATCH", "/account/1234/tasks/1234", strings.NewReader(`{"status":5}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(400, resp.StatusCode)
	})

	t.Run("update task but service has error should return error", func() {
		t.taskService.EXPECT().UpdateTaskStatus(gomock.Any(), "1234", "1234", 1).Return(errors.New("update task error"))
		// Define Fiber app.
		app := fiber.New()
		// Create route with PATCH method for test
		app.Patch("/account/:ownerId/tasks/:taskId", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.UpdateTask(c)
		})
		req := httptest.NewRequest("PATCH", "/account/1234/tasks/1234", strings.NewReader(`{"status":1}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(500, resp.StatusCode)
	})

	t.Run("update task success return task", func() {
		t.taskService.EXPECT().UpdateTaskStatus(gomock.Any(), "1234", "1234", 1).Return(nil)
		// Define Fiber app.
		app := fiber.New()
		// Create route with PATCH method for test
		app.Patch("/account/:ownerId/tasks/:taskId", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.UpdateTask(c)
		})
		req := httptest.NewRequest("PATCH", "/account/1234/tasks/1234", strings.NewReader(`{"status":1}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(200, resp.StatusCode)
		b, _ := io.ReadAll(resp.Body)
		t.Equal(`{"data":"Task status updated successfully"}`, string(b))
	})

}

func (t HandlerTestSuite) TestGetAllTask() {
	t.Run("get all task but service has error should return error", func() {
		t.taskService.EXPECT().GetAllTask(gomock.Any(), 1, 10).Return(nil, errors.New("get all task error"))
		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/tasks", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetAllTask(c)
		})
		req := httptest.NewRequest("GET", "/tasks", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(500, resp.StatusCode)
	})

	t.Run("get all task but over limit should return error", func() {
		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/tasks", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetAllTask(c)
		})
		req := httptest.NewRequest("GET", "/tasks?page=1&limit=1000", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(400, resp.StatusCode)
	})

	t.Run("get all task but page not number should return error", func() {
		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/tasks", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetAllTask(c)
		})
		req := httptest.NewRequest("GET", "/tasks?page=xxx&limit=1000", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(400, resp.StatusCode)
	})

	t.Run("get all task but limit not number should return error", func() {
		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/tasks", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetAllTask(c)
		})
		req := httptest.NewRequest("GET", "/tasks?page=1&limit=xxx", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(400, resp.StatusCode)
	})

	t.Run("get all task success return task", func() {
		t.taskService.EXPECT().GetAllTask(gomock.Any(), 1, 10).Return([]taskmanager.TaskDoc{
			{
				ID:          "1234",
				OwnerID:     "12345",
				Topic:       "test_topic",
				Description: "mock_desv",
				Status:      1,
				CreateDate:  2131341,
			},
		}, nil)
		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/tasks", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetAllTask(c)
		})
		req := httptest.NewRequest("GET", "/tasks", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(200, resp.StatusCode)
		b, _ := io.ReadAll(resp.Body)
		t.Equal(`{"data":[{"id":"1234","topic":"test_topic","description":"mock_desv","status":1,"create_date":2131341,"owner_id":"12345","archive_date":null,"update_date":null}]}`, string(b))
	})
}

func (t HandlerTestSuite) TestArchiveTask() {
	t.Run("archive task but service has error should return error", func() {
		t.taskService.EXPECT().ArchiveTask(gomock.Any(), "1234", "134134134").Return(0, errors.New("archive task error"))
		// Define Fiber app.
		app := fiber.New()
		// Create route with PATCH method for test
		app.Patch("/account/:ownerId/tasks/:taskId/archive", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.ArchiveTask(c)
		})
		req := httptest.NewRequest("PATCH", "/account/1234/tasks/134134134/archive", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(500, resp.StatusCode)
	})

	t.Run("archive task success but modification count is 0 return 400", func() {
		t.taskService.EXPECT().ArchiveTask(gomock.Any(), "1234", "134134134").Return(0, nil)
		// Define Fiber app.
		app := fiber.New()
		// Create route with PATCH method for test
		app.Patch("/account/:ownerId/tasks/:taskId/archive", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.ArchiveTask(c)
		})
		req := httptest.NewRequest("PATCH", "/account/1234/tasks/134134134/archive", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(400, resp.StatusCode)
	})

	t.Run("archive task success return task", func() {
		t.taskService.EXPECT().ArchiveTask(gomock.Any(), "1234", "134134134").Return(1, nil)
		// Define Fiber app.
		app := fiber.New()
		// Create route with PATCH method for test
		app.Patch("/account/:ownerId/tasks/:taskId/archive", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.ArchiveTask(c)
		})
		req := httptest.NewRequest("PATCH", "/account/1234/tasks/134134134/archive", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(200, resp.StatusCode)
		b, _ := io.ReadAll(resp.Body)
		t.Equal(`{"data":"Task archived successfully"}`, string(b))
	})
}

func (t *HandlerTestSuite) TestGetTopicComments() {
	t.Run("get topic comments but service has error should return error", func() {
		t.commentService.EXPECT().GetTopicComments(gomock.Any(), "134134134", 1, 10).Return(nil, errors.New("get topic comments error"))
		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/account/:ownerId/tasks/:taskId/comments", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetTopicComments(c)
		})
		req := httptest.NewRequest("GET", "/account/1234/tasks/134134134/comments?page=1&limit=10", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(500, resp.StatusCode)
	})

	t.Run("get topic comments but page not number should return error", func() {
		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/account/:ownerId/tasks/:taskId/comments", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetTopicComments(c)
		})
		req := httptest.NewRequest("GET", "/account/1234/tasks/134134134/comments?page=xxx&limit=10", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(400, resp.StatusCode)
	})

	t.Run("get topic comments but limit not number should return error", func() {
		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/account/:ownerId/tasks/:taskId/comments", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetTopicComments(c)
		})
		req := httptest.NewRequest("GET", "/account/1234/tasks/134134134/comments?page=1&limit=xxx", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(400, resp.StatusCode)
	})

	t.Run("get topic comments but over limit return erro 400 r", func() {
		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/account/:ownerId/tasks/:taskId/comments", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetTopicComments(c)
		})
		req := httptest.NewRequest("GET", "/account/1234/tasks/134134134/comments?page=1&limit=500", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(400, resp.StatusCode)
	})

	t.Run("get topic comments success return task", func() {
		t.commentService.EXPECT().GetTopicComments(gomock.Any(), "134134134", 1, 10).Return([]comment.CommentDoc{
			{
				ID:      "1234",
				TaskId:  "134134134",
				Content: "test_comment",
				OwnerId: "12345",
			},
		}, nil)
		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/account/:ownerId/tasks/:taskId/comments", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetTopicComments(c)
		})

		req := httptest.NewRequest("GET", "/account/1234/tasks/134134134/comments?page=1&limit=10", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(200, resp.StatusCode)
		b, _ := io.ReadAll(resp.Body)
		t.Equal(`{"data":[{"id":"1234","owner_id":"12345","task_id":"134134134","content":"test_comment","create_date":0,"update_date":null}]}`, string(b))
	})
}

func (t *HandlerTestSuite) TestCreateComment() {
	t.Run("create comment but service has error should return error", func() {
		t.commentService.EXPECT().CreateComment(gomock.Any(), "1234", "134134134", "test_comment").Return(nil, errors.New("create comment error"))
		t.profileService.EXPECT().GetProfile(gomock.Any(), "1234").Return(&profile.ProfileDoc{}, nil)

		// Define Fiber app.
		app := fiber.New()
		// Create route with POST method for test
		app.Post("/account/:ownerId/tasks/:taskId/comments", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.CreateComment(c)
		})
		req := httptest.NewRequest("POST", "/account/1234/tasks/134134134/comments", strings.NewReader(`{"content":"test_comment"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(500, resp.StatusCode)
	})

	t.Run("create comment but content is empty should return error", func() {
		// Define Fiber app.
		app := fiber.New()
		// Create route with POST method for test
		app.Post("/account/:ownerId/tasks/:taskId/comments", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.CreateComment(c)
		})
		req := httptest.NewRequest("POST", "/account/1234/tasks/134134134/comments", strings.NewReader(`{"content":""}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(400, resp.StatusCode)
	})

	t.Run("create comment success should return comment", func() {
		t.commentService.EXPECT().CreateComment(gomock.Any(), "1234", "134134134", "test_comment").Return(&comment.CommentDoc{
			ID:      "1234",
			TaskId:  "134134134",
			Content: "test_comment",
			OwnerId: "12345",
		}, nil)
		t.profileService.EXPECT().GetProfile(gomock.Any(), "1234").Return(&profile.ProfileDoc{}, nil)

		// Define Fiber app.
		app := fiber.New()
		// Create route with POST method for test
		app.Post("/account/:ownerId/tasks/:taskId/comments", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.CreateComment(c)
		})
		req := httptest.NewRequest("POST", "/account/1234/tasks/134134134/comments", strings.NewReader(`{"content":"test_comment"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(201, resp.StatusCode)
		b, _ := io.ReadAll(resp.Body)
		t.Equal(`{"data":{"id":"1234","owner_id":"12345","task_id":"134134134","content":"test_comment","create_date":0,"update_date":null}}`, string(b))
	})
}

func (t *HandlerTestSuite) TestGetProfile() {
	t.Run("get profile but service has error should return error", func() {
		t.profileService.EXPECT().GetProfile(gomock.Any(), "1234").Return(nil, errors.New("get profile error"))

		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/account/:ownerId/profile", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetProfile(c)
		})
		req := httptest.NewRequest("GET", "/account/1234/profile", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(500, resp.StatusCode)
	})

	t.Run("get profile but got nil", func() {
		t.profileService.EXPECT().GetProfile(gomock.Any(), "1234").Return(nil, nil)

		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/account/:ownerId/profile", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetProfile(c)
		})
		req := httptest.NewRequest("GET", "/account/1234/profile", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(400, resp.StatusCode)
	})

	t.Run("get profile success should return profile", func() {
		t.profileService.EXPECT().GetProfile(gomock.Any(), "1234").Return(&profile.ProfileDoc{
			OwnerId:     "user_id",
			DisplayName: "display_name",
			Email:       "email",
			UpdateDate:  1569152551,
			CreateDate:  1569152551,
			DisplayPic:  "url",
		}, nil)

		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/account/:ownerId/profile", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetProfile(c)
		})
		req := httptest.NewRequest("GET", "/account/1234/profile", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(200, resp.StatusCode)
		b, _ := io.ReadAll(resp.Body)
		t.Equal(`{"data":{"owner_id":"user_id","display_name":"display_name","email":"email","display_pic":"url"}}`, string(b))
	})
}

func (t *HandlerTestSuite) TestGetProfileList() {

	t.Run("get profile list but service has error should return error", func() {
		t.profileService.EXPECT().GetProfileList(gomock.Any(), []string{"1234"}).Return(nil, errors.New("get profile error"))

		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/account/profiles", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetProfileList(c)
		})
		req := httptest.NewRequest("GET", "/account/profiles?owner_id=1234", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(500, resp.StatusCode)
	})

	t.Run("get profile list but got nil", func() {
		t.profileService.EXPECT().GetProfileList(gomock.Any(), []string{"1234"}).Return(nil, nil)

		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/account/profiles", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetProfileList(c)
		})
		req := httptest.NewRequest("GET", "/account/profiles?owner_id=1234", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(200, resp.StatusCode)
	})

	t.Run("get profile list but over limit", func() {
		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/account/profiles", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetProfileList(c)
		})
		req := httptest.NewRequest("GET", "/account/profiles?owner_id=1234,5454,542525,1111,333", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(400, resp.StatusCode)
	})

	t.Run("get profile list but empty query param owner_id", func() {
		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/account/profiles", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetProfileList(c)
		})
		req := httptest.NewRequest("GET", "/account/profiles", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(400, resp.StatusCode)
	})

	t.Run("get profile list success should return profile", func() {
		t.profileService.EXPECT().GetProfileList(gomock.Any(), []string{"1234", "5454"}).Return([]profile.ProfileDoc{
			{
				OwnerId:     "user_id",
				DisplayName: "display_name",
				Email:       "email",
				UpdateDate:  1569152551,
				CreateDate:  1569152551,
				DisplayPic:  "url",
			},
		}, nil)

		// Define Fiber app.
		app := fiber.New()
		// Create route with GET method for test
		app.Get("/account/profiles", func(c *fiber.Ctx) error {
			// Return simple string as response
			return t.handler.GetProfileList(c)
		})
		req := httptest.NewRequest("GET", "/account/profiles?owner_id=1234,5454", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, 20)
		t.Equal(200, resp.StatusCode)
		b, _ := io.ReadAll(resp.Body)
		t.Equal(`{"data":[{"owner_id":"user_id","display_name":"display_name","email":"email","display_pic":"url"}]}`, string(b))
	})
}
