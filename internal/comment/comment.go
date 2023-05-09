package comment

import (
	"context"
	iMongo "task-manager-api/internal/mongo"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate mockgen -source=./comment.go -destination=./mock/comment_mock.go
type IMongo interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) iMongo.SingleResult
}

type CommentDoc struct {
	ID         string `json:"id" bson:"_id"`
	OwnerID    string `json:"owner_id" bson:"owner_id"`
	TopicID    string `json:"topic_id" bson:"topic_id"`
	Content    string `json:"content" bson:"content"`
	CreateDate int64  `bson:"create_date"`
	UpdateDate int64  `bson:"update_date"`
}

type Comment struct {
	mongo IMongo
}

func NewCommentService(mongo IMongo) *Comment {
	return &Comment{mongo: mongo}
}

func (c *Comment) CreateComment(ctx context.Context, ownerID string, topicID string, content string) error {
	return nil
}

func (c *Comment) GetComment(ctx context.Context, id string) (*CommentDoc, error) {
	return nil, nil
}
