package comment

import (
	"context"
	"errors"
	"reflect"
	mock "task-manager-api/internal/mongo/mock"
	"testing"
	"time"

	mock_comment "task-manager-api/internal/comment/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommentTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	mockMongo    *mock_comment.MockIMongo
	service      *Comment
	singleResult *mock.MockSingleResult
	cursor       *mock.MockCursor
}

func (t *CommentTestSuite) SetupTest() {
	t.ctrl = gomock.NewController(t.T())
	t.mockMongo = mock_comment.NewMockIMongo(t.ctrl)
	t.service = NewCommentService(t.mockMongo)
	t.singleResult = mock.NewMockSingleResult(t.ctrl)
	t.cursor = mock.NewMockCursor(t.ctrl)
	t.service.time = func() time.Time {
		loc, _ := time.LoadLocation("Asia/Bangkok")
		return time.Date(2019, 9, 22, 12, 42, 31, 0, loc)
	}
}

func (t *CommentTestSuite) TearDownTest() {
	t.ctrl.Finish()
	t.mockMongo = nil
	t.service = nil
	t.singleResult = nil
	t.cursor = nil
}

func TestCommentTestSuite(t *testing.T) {
	suite.Run(t, new(CommentTestSuite))
}

func (t *CommentTestSuite) TestGetTopicComments() {
	l := int64(10)
	skip := int64(1*10 - 10)
	fOpt := &options.FindOptions{Limit: &l, Skip: &skip}

	t.Run("get topic comments but find has error should return error", func() {
		t.mockMongo.EXPECT().Find(context.Background(), bson.M{
			"task_id": "topic_id",
		}, fOpt).Return(nil, errors.New("find error"))
		comments, err := t.service.GetTopicComments(context.Background(), "topic_id", 1, 10)
		t.Error(err)
		t.Nil(comments)
		t.EqualError(err, "find error")
	})
	t.Run("get topic comments but decode has error should return error", func() {
		t.mockMongo.EXPECT().Find(context.Background(), bson.M{
			"task_id": "topic_id",
		}, fOpt).Return(t.cursor, nil)
		t.cursor.EXPECT().All(context.Background(), gomock.Any()).DoAndReturn(func(ctx context.Context, result interface{}) error {
			return errors.New("cursor decode error")
		}).Times(1)
		comments, err := t.service.GetTopicComments(context.Background(), "topic_id", 1, 10)
		t.Error(err)
		t.Nil(comments)
		t.EqualError(err, "cursor decode error")
	})

	t.Run("get topic comments should return comments", func() {
		t.mockMongo.EXPECT().Find(context.Background(), bson.M{
			"task_id": "topic_id",
		}, fOpt).Return(t.cursor, nil)
		var comments = make([]CommentDoc, 0)
		t.cursor.EXPECT().All(context.Background(), &comments).DoAndReturn(func(ctx context.Context, result interface{}) error {
			comments = append(comments, CommentDoc{
				ID:         "comment_id",
				TaskId:     "topic_id",
				Content:    "content",
				CreateDate: 1569151351,
				OwnerId:    "owner_id",
			})
			reflect.ValueOf(result).Elem().Set(reflect.ValueOf(comments))
			return nil
		}).Times(1)
		comments, err := t.service.GetTopicComments(context.Background(), "topic_id", 1, 10)
		t.NoError(err)
		t.NotNil(comments)
		t.Equal(1, len(comments))
		t.Equal("comment_id", comments[0].ID)
		t.Equal("topic_id", comments[0].TaskId)
		t.Equal("content", comments[0].Content)
		t.Equal(int64(1569151351), comments[0].CreateDate)
		t.Equal("owner_id", comments[0].OwnerId)
	})
}

func (t *CommentTestSuite) TestCreateComment() {
	t.Run("create comment but insert has error should return error", func() {
		t.mockMongo.EXPECT().InsertOne(context.Background(), CommentDoc{
			TaskId:     "topic_id",
			Content:    "content",
			OwnerId:    "owner_id",
			CreateDate: int64(1569130951),
		}).Return(nil, errors.New("insert error"))
		comment, err := t.service.CreateComment(context.Background(), "owner_id", "topic_id", "content")
		t.Error(err)
		t.Nil(comment)
		t.EqualError(err, "insert error")
	})

	t.Run("create comment but can not convert _id", func() {
		t.mockMongo.EXPECT().InsertOne(context.Background(), CommentDoc{
			TaskId:     "topic_id",
			Content:    "content",
			OwnerId:    "owner_id",
			CreateDate: int64(1569130951),
		}).Return(&mongo.InsertOneResult{
			InsertedID: "objId",
		}, nil)
		comment, err := t.service.CreateComment(context.Background(), "owner_id", "topic_id", "content")
		t.Error(err)
		t.Nil(comment)
		t.EqualError(err, "cannot convert inserted id to object id")
	})

	t.Run("create comment should return comment", func() {
		objId, _ := primitive.ObjectIDFromHex("5ad9a913478c26d220afb681")
		t.mockMongo.EXPECT().InsertOne(context.Background(), CommentDoc{
			TaskId:     "topic_id",
			Content:    "content",
			OwnerId:    "owner_id",
			CreateDate: int64(1569130951),
		}).Return(&mongo.InsertOneResult{
			InsertedID: objId,
		}, nil)
		comment, err := t.service.CreateComment(context.Background(), "owner_id", "topic_id", "content")
		t.NoError(err)
		t.NotNil(comment)
		t.Equal("5ad9a913478c26d220afb681", comment.ID)
		t.Equal("topic_id", comment.TaskId)
		t.Equal("content", comment.Content)
		t.Equal(int64(1569130951), comment.CreateDate)
		t.Equal("owner_id", comment.OwnerId)
	})
}
