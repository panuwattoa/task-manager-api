package taskmanager

import (
	"context"
	iMongo "task-manager-api/internal/mongo"
	"time"

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
	time  func() time.Time
}

func NewTaskManager(mongo IMongo) *TaskManager {
	return &TaskManager{mongo: mongo}
}

const (
	TaskStatusOpen = iota + 1
	TaskStatusInProgress
	TaskStatusDone
)

type TaskDoc struct {
	ID          string `bson:"_id,omitempty"`
	Topic       string `bson:"topic"`
	Description string `bson:"description"`
	Status      int    `bson:"status"`
	CreateDate  int64  `bson:"create_date"`
	OwnerID     string `bson:"owner_id"`
	ArchiveDate *int64 `bson:"archive_date"`
	UpdateDate  *int64 `bson:"update_date"`
}

func (t *TaskManager) CreateTask(ctx context.Context, topic string, desc string, ownerId string) (*TaskDoc, error) {
	return nil, nil
}

func (t *TaskManager) GetTask(ctx context.Context, id string) (*TaskDoc, error) {
	return nil, nil
}

func (t *TaskManager) UpdateTaskStatus(ctx context.Context, id string, status int) error {
	return nil
}

func (t *TaskManager) ArchiveTask(ctx context.Context, id string) error {
	return nil
}

func (t *TaskManager) now() time.Time {
	if t.time == nil {
		return time.Now()
	}

	return t.time()
}
