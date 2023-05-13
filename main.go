package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"task-manager-api/config"
	"task-manager-api/internal/handler"
	"task-manager-api/internal/mongo"
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
	taskService := taskmanager.NewTaskManager(mongoTaskCollection)
	handler := handler.NewHandler(taskService)
	app := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: errorInterceptor,
	})

	app.Post("account/:ownerId/tasks", handler.CreateTask)
	app.Get("/tasks", handler.GetAllTask)

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

func errorInterceptor(ctx *fiber.Ctx, err error) error {
	// Override default error handler
	// Status code defaults to 500
	code := fiber.StatusInternalServerError
	// Retrieve the custom status code if it's an fiber.*Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

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
