package main

import (
	"context"
	"fmt"
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
	readConfiguration()
	// Mongo
	mongoDB := mongo.NewMongoDB()

	if err := mongoDB.Open(context.Background()); err != nil {
		log.Fatal(err)
	}
	if err := mongoDB.Status(context.Background()); err != nil {
		log.Fatal(err)
	}

	mongoTaskCollection := mongo.NewCollectionHelper(mongoDB.GetCollection(config.Conf.MongoDB.Collections.Tasks))
	profileCollection := mongo.NewCollectionHelper(mongoDB.GetCollection(config.Conf.MongoDB.Collections.Profiles))
	commentCollection := mongo.NewCollectionHelper(mongoDB.GetCollection(config.Conf.MongoDB.Collections.Comments))

	taskService := taskmanager.NewTaskManager(mongoTaskCollection)
	pfService := profile.NewProfileService(profileCollection)
	commentService := comment.NewCommentService(commentCollection)
	handler := handler.NewHandler(taskService, commentService, pfService)

	app := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: errorInterceptor,
	})

	app.Get("/tasks", handler.GetAllTask)
	app.Get("/tasks/:taskId", handler.GetTask)
	app.Get("/profiles/:ownerId", handler.GetProfile)
	app.Get("/profiles", handler.GetProfileList)
	app.Get("/tasks/:taskId/comments", handler.GetTopicComments)

	customerGroup := app.Group("/account")
	customerGroup.Use(
		func(c *fiber.Ctx) error {
			return authInterceptor(c)
		},
	)
	customerGroup.Post(":ownerId/tasks", handler.CreateTask)
	customerGroup.Post(":ownerId/tasks/:taskId/comments", handler.CreateComment)
	customerGroup.Patch(":ownerId/tasks/:taskId", handler.UpdateTask)
	customerGroup.Delete(":ownerId/tasks/:taskId", handler.ArchiveTask)

	go func() {
		if err := app.Listen(":" + config.Conf.Server.Port); err != nil {
			fmt.Print("can not start http server")
			log.Fatal(err)
		}
	}()
	// make SIGINT send context cancel for graceful stop
	gfs := make(chan os.Signal, 1)
	signal.Notify(gfs, syscall.SIGTERM, syscall.SIGINT)
	<-gfs

	_ = app.Shutdown()
	// stop mongo db
	err := mongoDB.Close(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func authInterceptor(ctx *fiber.Ctx) error {
	// TODO: validate Authorization header

	return ctx.Next()
}

func errorInterceptor(ctx *fiber.Ctx, err error) error {
	// Override default error handler
	// Status code defaults to 500
	code := fiber.StatusInternalServerError
	// Retrieve the custom status code if it's an fiber.*Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	fmt.Printf("Error: %v\n", err)
	msg := map[string]interface{}{
		"status":    code,
		"error_msg": err.Error(),
	}
	ctx.Set("Content-Type", "application/json")
	return ctx.Status(code).JSON(msg)
}

func readConfiguration() {
	if err := config.Read(config.Conf); err != nil {
		os.Exit(78) // 78 - Configuration error
	}

	return
}
