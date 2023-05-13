package taskmanager

import (
	"context"
	"errors"
	m "task-manager-api/internal/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate mockgen -source=./taskmanager.go -destination=./mock/taskmanager.go
type IMongo interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) m.SingleResult
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (m.Cursor, error)
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
	ID          string `json:"id" bson:"_id,omitempty"`
	Topic       string `json:"topic" bson:"topic"`
	Description string `json:"description" bson:"description"`
	Status      int    `json:"status" bson:"status"`
	CreateDate  int64  `json:"create_date" bson:"create_date"`
	OwnerID     string `json:"owner_id" bson:"owner_id"`
	ArchiveDate *int64 `json:"archive_date" bson:"archive_date"`
	UpdateDate  *int64 `json:"update_date" bson:"update_date"`
}

func (t *TaskManager) CreateTask(ctx context.Context, ownerId string, topic string, desc string) (*TaskDoc, error) {
	// create new task
	// TODO: some other business logic here
	now := t.now().Unix()
	result, err := t.mongo.InsertOne(ctx, TaskDoc{
		Topic:       topic,
		Description: desc,
		Status:      TaskStatusOpen,
		CreateDate:  now,
		OwnerID:     ownerId,
	})
	if err != nil {
		// TODO: log error
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return &TaskDoc{
			ID:          oid.Hex(),
			Topic:       topic,
			Description: desc,
			Status:      TaskStatusOpen,
			CreateDate:  now,
			OwnerID:     ownerId,
		}, nil
	} else {
		// TODO: log error
		return nil, errors.New("cannot convert inserted id to object id")
	}
}

func (t *TaskManager) GetAllTask(ctx context.Context, page int, limit int) ([]TaskDoc, error) {
	// find all task with pagination
	curr, err := t.mongo.Find(ctx, bson.M{
		"$or": []bson.M{
			{
				"archive_date": bson.M{
					"$exists": false,
				},
			},
			{
				"archive_date": nil,
			},
		},
	}, m.NewMongoPaginate(limit, page).GetPaginatedOpts())

	if err != nil {
		// TODO: log error
		return nil, err
	}
	var tasks = make([]TaskDoc, 0)
	if err := curr.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t *TaskManager) GetTask(ctx context.Context, id string) (*TaskDoc, error) {
	// find task by id
	objectId, _ := primitive.ObjectIDFromHex(id)
	curr := t.mongo.FindOne(ctx, bson.M{
		"_id": objectId,
		"$or": []bson.M{
			{
				"archive_date": bson.M{
					"$exists": false,
				},
			},
			{
				"archive_date": nil,
			},
		},
	}, &options.FindOneOptions{
		Sort: bson.M{
			"_id": 1,
		},
	})
	var task TaskDoc
	if err := curr.Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (t *TaskManager) UpdateTaskStatus(ctx context.Context, ownerId string, id string, status int) error {
	// update task status
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := t.mongo.UpdateOne(ctx, bson.M{
		"_id":      objectId,
		"owner_id": ownerId,
	}, bson.M{
		"$set": bson.M{
			"status":      status,
			"update_date": t.now().Unix(),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskManager) ArchiveTask(ctx context.Context, ownerId string, id string) (int, error) {
	// archive task
	objectId, _ := primitive.ObjectIDFromHex(id)
	results, err := t.mongo.UpdateOne(ctx, bson.M{
		"_id":      objectId,
		"owner_id": ownerId,
	}, bson.M{
		"$set": bson.M{
			"archive_date": t.now().Unix(),
			"update_date":  t.now().Unix(),
		},
	})
	if err != nil {
		return 0, err
	}
	return int(results.MatchedCount), nil
}

func (t *TaskManager) now() time.Time {
	if t.time == nil {
		return time.Now()
	}

	return t.time()
}
