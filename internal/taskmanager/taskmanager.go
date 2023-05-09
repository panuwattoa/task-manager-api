package taskmanager

import (
	"context"
	iMongo "task-manager-api/internal/mongo"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate mockgen -source=./taskmanager.go -destination=./mock/taskmanager_mock.go
type IMongo interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) iMongo.SingleResult
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

type TaskManager struct {
	mongo IMongo
}

func NewTaskManager(mongo IMongo) *TaskManager {
	return &TaskManager{mongo: mongo}
}

func (t *TaskManager) CreateTask(ctx context.Context) error {
	return nil
}
