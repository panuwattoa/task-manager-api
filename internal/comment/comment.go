package comment

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

//go:generate mockgen -source=./comment.go -destination=./mock/comment.go
type IMongo interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (m.Cursor, error)
}

type CommentDoc struct {
	ID         string `json:"id" bson:"_id,omitempty"`
	OwnerId    string `json:"owner_id" bson:"owner_id"`
	TaskId     string `json:"task_id" bson:"task_id"`
	Content    string `json:"content" bson:"content"`
	CreateDate int64  `json:"create_date" bson:"create_date"`
	UpdateDate *int64 `json:"update_date" bson:"update_date"`
}

type Comment struct {
	mongo IMongo
	time  func() time.Time
}

func NewCommentService(mongo IMongo) *Comment {
	return &Comment{mongo: mongo}
}

func (c *Comment) CreateComment(ctx context.Context, ownerId string, TaskId string, content string) (*CommentDoc, error) {
	// TODO: validate topicID
	now := c.now().Unix()
	result, err := c.mongo.InsertOne(ctx, CommentDoc{
		OwnerId:    ownerId,
		TaskId:     TaskId,
		Content:    content,
		CreateDate: now,
	})
	if err != nil {
		return nil, err
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return &CommentDoc{
			ID:         oid.Hex(),
			TaskId:     TaskId,
			Content:    content,
			CreateDate: now,
			OwnerId:    ownerId,
		}, nil
	} else {
		// TODO: log error
		return nil, errors.New("cannot convert inserted id to object id")
	}
}

func (c *Comment) GetTopicComments(ctx context.Context, TaskId string, page int, limit int) ([]CommentDoc, error) {
	// find all comment in topic with pagination
	curr, err := c.mongo.Find(ctx, bson.M{
		"task_id": TaskId,
	}, m.NewMongoPaginate(limit, page).GetPaginatedOpts())
	if err != nil {
		return nil, err
	}

	var comments = make([]CommentDoc, 0)
	if err := curr.All(ctx, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *Comment) now() time.Time {
	if c.time == nil {
		return time.Now()
	}

	return c.time()
}
