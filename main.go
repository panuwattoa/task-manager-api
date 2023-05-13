package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"task-manager-api/config"
	"task-manager-api/internal/comment"
	"task-manager-api/internal/handler"
	"task-manager-api/internal/mongo"
	"task-manager-api/internal/profile"
	"task-manager-api/internal/taskmanager"

	"github.com/gofiber/fiber/v2"
)

func main() {
	err := config.Read(config.Conf)
	if err != nil {
		log.Fatalf("failed to read configuration: %v", err)
	}

	// Initialize mongoDB
	mongoDB := mongo.NewMongoDB()
	defer mongoDB.Close(context.Background())
	if err := mongoDB.Open(context.Background()); err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	if err := mongoDB.Status(context.Background()); err != nil {
		log.Fatalf("MongoDB health check failed: %v", err)
	}

	// Initialize collections
	mongoTaskCollection := mongoDB.GetCollection(config.Conf.MongoDB.Collections.Tasks)
	profileCollection := mongoDB.GetCollection(config.Conf.MongoDB.Collections.Profiles)
	commentCollection := mongoDB.GetCollection(config.Conf.MongoDB.Collections.Comments)

	// Initialize services and handlers
	taskService := taskmanager.NewTaskManager(mongo.NewCollectionHelper(mongoTaskCollection))
	pfService := profile.NewProfileService(mongo.NewCollectionHelper(profileCollection))
	commentService := comment.NewCommentService(mongo.NewCollectionHelper(commentCollection))
	handler := handler.NewHandler(taskService, commentService, pfService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	// Define routes
	app.Get("/tasks", handler.GetAllTask)
	app.Get("/tasks/:taskId", handler.GetTask)
	app.Get("/profiles/:ownerId", handler.GetProfile)
	app.Get("/profiles", handler.GetProfileList)
	app.Get("/tasks/:taskId/comments", handler.GetTopicComments)

	customerGroup := app.Group("/account")
	customerGroup.Use(authInterceptor)
	customerGroup.Post(":ownerId/tasks", handler.CreateTask)
	customerGroup.Post(":ownerId/tasks/:taskId/comments", handler.CreateComment)
	customerGroup.Patch(":ownerId/tasks/:taskId", handler.UpdateTask)
	customerGroup.Patch(":ownerId/tasks/:taskId/archive", handler.ArchiveTask)

	// Start HTTP server
	go func() {
		if err := app.Listen(":" + config.Conf.Server.Port); err != nil {
			log.Fatalf("failed to start HTTP server: %v", err)
		}
	}()

	// Wait for SIGTERM or SIGINT signal
	gracefully(app, mongoDB)
}

func authInterceptor(ctx *fiber.Ctx) error {
	// TODO: validate Authorization header
	return ctx.Next()
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	log.Printf("error: %v", err)
	msg := map[string]interface{}{
		"status":    code,
		"error_msg": err.Error(),
	}
	return ctx.Status(code).JSON(msg)
}

func gracefully(app *fiber.App, mongoDB *mongo.MongoDB) {
	// Make SIGINT send context cancel for graceful stop
	gfs := make(chan os.Signal, 1)
	signal.Notify(gfs, syscall.SIGTERM, syscall.SIGINT)
	<-gfs

	// Shutdown server
	if err := app.Shutdown(); err != nil {
		log.Fatal(err)
	}

	// Stop mongo db
	if err := mongoDB.Close(context.Background()); err != nil {
		log.Fatal(err)
	}
}
