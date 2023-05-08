package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"task-manager-api/internal/mongo"
)

func main() {
	// Mongo
	mongoDB := mongo.NewMongoDB()

	if err := mongoDB.Open(context.Background()); err != nil {
		// TODO: log fatal error
	}
	if err := mongoDB.Status(context.Background()); err != nil {
		// TODO: log fatal error
	}

	// make SIGINT send context cancel for graceful stop
	gfs := make(chan os.Signal, 1)
	signal.Notify(gfs, syscall.SIGTERM, syscall.SIGINT)
	<-gfs
	// stop mongo db
	err := mongoDB.Close(context.Background())
	if err != nil {
		// TODO: log error
	}
}
