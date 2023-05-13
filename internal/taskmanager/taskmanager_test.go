package taskmanager

import (
	"context"
	"errors"
	"reflect"
	mock "task-manager-api/internal/mongo/mock"
	mock_taskmanager "task-manager-api/internal/taskmanager/mock"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskManagerTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	mockMongo    *mock_taskmanager.MockIMongo
	service      *TaskManager
	singleResult *mock.MockSingleResult
	cursor       *mock.MockCursor
}

func (t *TaskManagerTestSuite) SetupTest() {
	t.ctrl = gomock.NewController(t.T())
	t.mockMongo = mock_taskmanager.NewMockIMongo(t.ctrl)
	t.service = NewTaskManager(t.mockMongo)
	t.singleResult = mock.NewMockSingleResult(t.ctrl)
	t.cursor = mock.NewMockCursor(t.ctrl)
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
	t.cursor = nil
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

	t.Run("create task success but can not convert _id", func() {
		t.mockMongo.EXPECT().InsertOne(context.Background(), TaskDoc{
			Topic:       "topic",
			Description: "description",
			Status:      1,
			CreateDate:  t.service.now().Unix(),
			OwnerID:     "owner_id",
		}).Return(&mongo.InsertOneResult{
			InsertedID: "5ad9a913478c26d220afb681",
		}, nil)
		taskDoc, err := t.service.CreateTask(context.Background(), "topic", "description", "owner_id")
		t.NotNil(err)
		t.EqualError(err, "cannot convert inserted id to object id")
		t.Nil(taskDoc)
	})

	t.Run("create task success", func() {
		objId, _ := primitive.ObjectIDFromHex("5ad9a913478c26d220afb681")
		t.mockMongo.EXPECT().InsertOne(context.Background(), TaskDoc{
			Topic:       "topic",
			Description: "description",
			Status:      1,
			CreateDate:  t.service.now().Unix(),
			OwnerID:     "owner_id",
		}).Return(&mongo.InsertOneResult{
			InsertedID: objId,
		}, nil)
		taskDoc, err := t.service.CreateTask(context.Background(), "topic", "description", "owner_id")
		t.Nil(err)
		t.NotNil(taskDoc)
		t.Equal("5ad9a913478c26d220afb681", taskDoc.ID)
		t.Equal("topic", taskDoc.Topic)
		t.Equal("description", taskDoc.Description)
		t.Equal(1, taskDoc.Status)
		t.Equal(int64(1569130951), taskDoc.CreateDate)
		t.Equal("owner_id", taskDoc.OwnerID)
		t.NoError(err)
	})
}

func (t *TaskManagerTestSuite) TestGetTask() {
	t.Run("get task but find one got error should return error ", func() {
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
		t.singleResult.EXPECT().Decode(&TaskDoc{}).Return(errors.New("something wrong")).Times(1)
		taskDoc, err := t.service.GetTask(context.Background(), "6041c3a6cfcba2fb9c4a4fd2")
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
		t.singleResult.EXPECT().Decode(&TaskDoc{}).Return(mongo.ErrNoDocuments).Times(1)
		taskDoc, err := t.service.GetTask(context.Background(), "6041c3a6cfcba2fb9c4a4fd2")
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
		t.Equal("6041c3a6cfcba2fb9c4a4fd2", taskDoc.ID)
		t.Equal("topic", taskDoc.Topic)
		t.Equal("description", taskDoc.Description)
		t.Equal(1, taskDoc.Status)
		t.Equal(int64(1614962551), taskDoc.CreateDate)
	})
}

func (t *TaskManagerTestSuite) TestArchiveTask() {
	t.Run("archive task but update one got error should return error", func() {
		objectId, _ := primitive.ObjectIDFromHex("6041c3a6cfcba2fb9c4a4fd2")
		t.mockMongo.EXPECT().UpdateOne(context.Background(), bson.M{
			"_id": objectId,
		}, bson.M{
			"$set": bson.M{
				"archive_date": t.service.now().Unix(),
				"update_date":  t.service.now().Unix(),
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
				"archive_date": t.service.now().Unix(),
				"update_date":  t.service.now().Unix(),
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
				"status":      1,
				"update_date": t.service.now().Unix(),
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
				"status":      2,
				"update_date": t.service.now().Unix(),
			},
		}).Return(&mongo.UpdateResult{
			MatchedCount: 1,
		}, nil)
		err := t.service.UpdateTaskStatus(context.Background(), "6041c3a6cfcba2fb9c4a4fd2", 2)
		t.NoError(err)
	})
}

func (t *TaskManagerTestSuite) TestGetAllTask() {
	l := int64(10)
	skip := int64(1*10 - 10)
	fOpt := &options.FindOptions{Limit: &l, Skip: &skip}

	t.Run("get all task but find got error should return error", func() {
		t.mockMongo.EXPECT().Find(context.Background(), bson.M{
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
		}, fOpt).Return(t.cursor, errors.New("find error"))
		_, err := t.service.GetAllTask(context.Background(), 1, 10)
		t.Error(err)
		t.EqualError(err, "find error")
	})

	t.Run("get all task success but cursor decode error", func() {
		t.mockMongo.EXPECT().Find(context.Background(), bson.M{
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
		}, fOpt).Return(t.cursor, nil)
		taskDocs := make([]TaskDoc, 0)
		t.cursor.EXPECT().All(context.Background(), &taskDocs).DoAndReturn(func(ctx context.Context, result interface{}) error {
			return errors.New("cursor decode error")
		}).Times(1)
		tasks, err := t.service.GetAllTask(context.Background(), 1, 10)
		t.NotNil(err)
		t.EqualError(err, "cursor decode error")
		t.Nil(tasks)
	})

	t.Run("get all task success", func() {
		t.mockMongo.EXPECT().Find(context.Background(), bson.M{
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
		}, fOpt).Return(t.cursor, nil)
		taskDocs := make([]TaskDoc, 0)
		t.cursor.EXPECT().All(context.Background(), &taskDocs).DoAndReturn(func(ctx context.Context, result interface{}) error {
			taskDocs = append(taskDocs, TaskDoc{
				ID:          "6041c3a6cfcba2fb9c4a4fd2",
				Topic:       "topic",
				Description: "description",
				Status:      1,
				CreateDate:  1614962551,
			})
			taskDocs = append(taskDocs, TaskDoc{
				ID:          "6041c3a6cfcba2fb9c4a4fd3",
				Topic:       "topic2",
				Description: "description2",
				Status:      2,
				CreateDate:  1614962552,
			})
			reflect.ValueOf(result).Elem().Set(reflect.ValueOf(taskDocs))
			return nil
		}).Times(1)
		tasks, err := t.service.GetAllTask(context.Background(), 1, 10)
		t.NoError(err)
		t.NotNil(tasks)
		t.Equal(2, len(tasks))
		t.Equal("6041c3a6cfcba2fb9c4a4fd2", tasks[0].ID)
		t.Equal("topic", tasks[0].Topic)
	})
}
