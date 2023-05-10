package taskmanager

import (
	"context"
	"errors"
	mock_sigleResult "task-manager-api/internal/mongo/mock"
	mock_taskmanager "task-manager-api/internal/taskmanager/mock"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskManagerTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	mockMongo    *mock_taskmanager.MockIMongo
	service      *TaskManager
	singleResult *mock_sigleResult.MockSingleResult
}

func (t *TaskManagerTestSuite) SetupTest() {
	t.ctrl = gomock.NewController(t.T())
	t.mockMongo = mock_taskmanager.NewMockIMongo(t.ctrl)
	t.service = NewTaskManager(t.mockMongo)
	t.singleResult = mock_sigleResult.NewMockSingleResult(t.ctrl)
	t.service.time = func() time.Time {
		loc, _ := time.LoadLocation("Asia/Bangkok")
		return time.Date(2019, 9, 22, 12, 42, 31, 0, loc)
	}
}

func (t *TaskManagerTestSuite) TearDownTest() {
	t.ctrl.Finish()
	t.mockMongo = nil
	t.service = nil
	t.singleResult = nil
}

func TestTaskManagerTestSuite(t *testing.T) {
	suite.Run(t, new(TaskManagerTestSuite))
}

func (t *TaskManagerTestSuite) TestCreateTask() {
	t.Run("create task but insert one has error should return error", func() {
		t.mockMongo.EXPECT().InsertOne(context.Background(), TaskDoc{
			Topic:       "topic",
			Description: "description",
			Status:      1,
			CreateDate:  t.service.now().Unix(),
			OwnerID:     "owner_id",
		}).Return(nil, errors.New("insert one error"))
		taskDoc, err := t.service.CreateTask(context.Background(), "topic", "description", "owner_id")
		t.Error(err)
		t.Nil(taskDoc)
		t.EqualError(err, "insert one error")
	})

	t.Run("create task success", func() {
		t.mockMongo.EXPECT().InsertOne(context.Background(), TaskDoc{
			Topic:       "topic",
			Description: "description",
			Status:      1,
			CreateDate:  t.service.now().Unix(),
			OwnerID:     "owner_id",
		}).Return(&mongo.InsertOneResult{
			InsertedID: 100,
		}, nil)
		taskDoc, err := t.service.CreateTask(context.Background(), "topic", "description", "owner_id")
		t.NotNil(taskDoc)
		t.ElementsMatch(
			*taskDoc,
			TaskDoc{
				ID:          "100",
				Topic:       "topic",
				Description: "description",
				Status:      1,
				CreateDate:  t.service.now().Unix(),
				OwnerID:     "owner_id",
			},
		)
		t.NoError(err)
	})
}

func (t *TaskManagerTestSuite) TestGetTask() {
	t.Run("get task but find one got error should return error ", func() {
		objectId, _ := primitive.ObjectIDFromHex("6041c3a6cfcba2fb9c4a4fd2")
		t.mockMongo.EXPECT().FindOne(context.Background(), bson.M{
			"_id": objectId,
		}).Return(t.singleResult)
		t.singleResult.EXPECT().Decode(nil).Return(errors.New("something wrong")).Times(1)
		taskDoc, err := t.service.GetTask(context.Background(), "id")
		t.Error(err)
		t.Nil(taskDoc)
		t.EqualError(err, "something wrong")
	})

	t.Run("get task but find one has not found should return error no documents", func() {
		objectId, _ := primitive.ObjectIDFromHex("6041c3a6cfcba2fb9c4a4fd2")
		t.mockMongo.EXPECT().FindOne(context.Background(), bson.M{
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
		}).Return(t.singleResult)
		t.singleResult.EXPECT().Decode(nil).Return(mongo.ErrNoDocuments).Times(1)
		taskDoc, err := t.service.GetTask(context.Background(), "id")
		t.Error(err)
		t.Nil(taskDoc)
		t.EqualError(err, "mongo: no documents in result")
	})

	t.Run("get task but decode error", func() {
		objectId, _ := primitive.ObjectIDFromHex("6041c3a6cfcba2fb9c4a4fd2")
		t.mockMongo.EXPECT().FindOne(context.Background(), bson.M{
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
		}).Return(t.singleResult)
		t.singleResult.EXPECT().Decode(&TaskDoc{}).Return(errors.New("decode error")).Times(1)
		taskDoc, err := t.service.GetTask(context.Background(), "6041c3a6cfcba2fb9c4a4fd2")
		t.Nil(taskDoc)
		t.Error(err)
		t.EqualError(err, "decode error")
	})

	t.Run("get task success", func() {
		objectId, _ := primitive.ObjectIDFromHex("6041c3a6cfcba2fb9c4a4fd2")
		// filter archive_date is null or not exist
		t.mockMongo.EXPECT().FindOne(context.Background(), bson.M{
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
		}).Return(t.singleResult)
		t.singleResult.EXPECT().Decode(&TaskDoc{}).DoAndReturn(func(doc *TaskDoc) error {
			doc.ID = "6041c3a6cfcba2fb9c4a4fd2"
			doc.Topic = "topic"
			doc.Description = "description"
			doc.Status = 1
			doc.CreateDate = 1614962551
			return nil
		}).Times(1)
		taskDoc, err := t.service.GetTask(context.Background(), "6041c3a6cfcba2fb9c4a4fd2")
		t.NotNil(taskDoc)
		t.NoError(err)
		t.ElementsMatch(
			*taskDoc,
			TaskDoc{
				ID:          "6041c3a6cfcba2fb9c4a4fd2",
				Topic:       "topic",
				Description: "description",
				Status:      1,
				CreateDate:  1614962551,
			},
		)
	})
}

func (t *TaskManagerTestSuite) TestArchiveTask() {
	t.Run("archive task but update one got error should return error", func() {
		objectId, _ := primitive.ObjectIDFromHex("6041c3a6cfcba2fb9c4a4fd2")
		t.mockMongo.EXPECT().UpdateOne(context.Background(), bson.M{
			"_id": objectId,
		}, bson.M{
			"$set": bson.M{
				"ArchiveDate": t.service.now().Unix(),
				"UpdateDate":  t.service.now().Unix(),
			},
		}).Return(nil, errors.New("update one error"))
		err := t.service.ArchiveTask(context.Background(), "6041c3a6cfcba2fb9c4a4fd2")
		t.Error(err)
		t.EqualError(err, "update one error")
	})

	t.Run("archive task success", func() {
		objectId, _ := primitive.ObjectIDFromHex("6041c3a6cfcba2fb9c4a4fd2")
		t.mockMongo.EXPECT().UpdateOne(context.Background(), bson.M{
			"_id": objectId,
		}, bson.M{
			"$set": bson.M{
				"ArchiveDate": t.service.now().Unix(),
				"UpdateDate":  t.service.now().Unix(),
			},
		}).Return(&mongo.UpdateResult{
			MatchedCount: 1,
		}, nil)
		err := t.service.ArchiveTask(context.Background(), "6041c3a6cfcba2fb9c4a4fd2")
		t.NoError(err)
	})
}

func (t *TaskManagerTestSuite) TestUpdateTaskStatus() {
	t.Run("update task status but update one got error should return error", func() {
		objectId, _ := primitive.ObjectIDFromHex("6041c3a6cfcba2fb9c4a4fd2")
		t.mockMongo.EXPECT().UpdateOne(context.Background(), bson.M{
			"_id": objectId,
		}, bson.M{
			"$set": bson.M{
				"Status":     1,
				"UpdateDate": t.service.now().Unix(),
			},
		}).Return(nil, errors.New("update one error"))
		err := t.service.UpdateTaskStatus(context.Background(), "6041c3a6cfcba2fb9c4a4fd2", 1)
		t.Error(err)
		t.EqualError(err, "update one error")
	})

	t.Run("update task status success", func() {
		objectId, _ := primitive.ObjectIDFromHex("6041c3a6cfcba2fb9c4a4fd2")
		t.mockMongo.EXPECT().UpdateOne(context.Background(), bson.M{
			"_id": objectId,
		}, bson.M{
			"$set": bson.M{
				"Status":     2,
				"UpdateDate": t.service.now().Unix(),
			},
		}).Return(&mongo.UpdateResult{
			MatchedCount: 1,
		}, nil)
		err := t.service.UpdateTaskStatus(context.Background(), "6041c3a6cfcba2fb9c4a4fd2", 2)
		t.NoError(err)
	})
}
