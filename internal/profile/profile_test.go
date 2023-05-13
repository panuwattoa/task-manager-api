package profile

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	mock "task-manager-api/internal/mongo/mock"
	mock_profile "task-manager-api/internal/profile/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProfileTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	mockMongo    *mock_profile.MockIMongo
	service      *Profile
	singleResult *mock.MockSingleResult
	cursor       *mock.MockCursor
}

func (t *ProfileTestSuite) SetupTest() {
	t.ctrl = gomock.NewController(t.T())
	t.mockMongo = mock_profile.NewMockIMongo(t.ctrl)
	t.service = NewProfileService(t.mockMongo)
	t.cursor = mock.NewMockCursor(t.ctrl)
	t.singleResult = mock.NewMockSingleResult(t.ctrl)
	t.service.time = func() time.Time {
		loc, _ := time.LoadLocation("Asia/Bangkok")
		return time.Date(2019, 9, 22, 12, 42, 31, 0, loc)
	}
}

func (t *ProfileTestSuite) TearDownTest() {
	t.ctrl.Finish()
	t.mockMongo = nil
	t.service = nil
	t.singleResult = nil
	t.cursor = nil
}

func TestProfileTestSuite(t *testing.T) {
	suite.Run(t, new(ProfileTestSuite))
}

func (t *ProfileTestSuite) TestGetProfile() {
	t.Run("get profile but decode has error should return error", func() {
		t.singleResult.EXPECT().Decode(gomock.Any()).Return(errors.New("decode error"))
		t.mockMongo.EXPECT().FindOne(context.Background(), bson.M{
			"owner_id": "user_id",
		}).Return(t.singleResult)
		profile, err := t.service.GetProfile(context.Background(), "user_id")
		t.Error(err)
		t.Nil(profile)
		t.EqualError(err, "decode error")
	})

	t.Run("get profile but error no document should return nil", func() {
		t.singleResult.EXPECT().Decode(gomock.Any()).Return(mongo.ErrNoDocuments)
		t.mockMongo.EXPECT().FindOne(context.Background(), bson.M{
			"owner_id": "user_id",
		}).Return(t.singleResult)
		profile, err := t.service.GetProfile(context.Background(), "user_id")
		t.Nil(err)
		t.Nil(profile)
	})
	t.Run("get profile success", func() {
		t.singleResult.EXPECT().Decode(&ProfileDoc{}).DoAndReturn(func(doc *ProfileDoc) error {
			doc.OwnerID = "user_id"
			doc.DisplayName = "display_name"
			doc.Email = "email"
			doc.UpdateDate = 1569152551
			doc.CreateDate = 1569152551
			doc.DisplayPic = "url"
			return nil
		}).Times(1)
		t.mockMongo.EXPECT().FindOne(context.Background(), bson.M{
			"owner_id": "user_id",
		}).Return(t.singleResult)
		profile, err := t.service.GetProfile(context.Background(), "user_id")
		t.NoError(err)
		t.NotNil(profile)
		t.Equal("user_id", profile.OwnerID)
		t.Equal("display_name", profile.DisplayName)
		t.Equal("email", profile.Email)
		t.Equal(int64(1569152551), profile.UpdateDate)
		t.Equal(int64(1569152551), profile.CreateDate)
		t.Equal("url", profile.DisplayPic)
	})

}

func (t *ProfileTestSuite) TestGetProgileList() {
	t.Run("get profile list but decode has error should return error", func() {
		t.mockMongo.EXPECT().Find(context.Background(), bson.M{
			"owner_id": bson.M{
				"$in": []string{
					"1", "2", "3",
				},
			}}).Return(t.cursor, errors.New("find error"))
		profile, err := t.service.GetProfileList(context.Background(), []string{
			"1", "2", "3",
		})
		t.Error(err)
		t.Nil(profile)
		t.EqualError(err, "find error")
	})

	t.Run("get profile list success but cursor error", func() {
		var profiles = make([]ProfileDoc, 0)
		t.cursor.EXPECT().All(context.Background(), &profiles).DoAndReturn(func(ctx context.Context, result interface{}) error {
			return errors.New("cursor error")
		}).Times(1)
		t.mockMongo.EXPECT().Find(context.Background(), bson.M{
			"owner_id": bson.M{
				"$in": []string{
					"1", "2", "3", "user_id",
				}},
		}).Return(t.cursor, nil)
		profile, err := t.service.GetProfileList(context.Background(), []string{
			"1", "2", "3", "user_id"})
		t.Error(err)
		t.Nil(profile)
		t.EqualError(err, "cursor error")
	})

	t.Run("get profile list success", func() {
		var profiles = make([]ProfileDoc, 0)
		t.cursor.EXPECT().All(context.Background(), &profiles).DoAndReturn(func(ctx context.Context, result interface{}) error {
			profiles = append(profiles, ProfileDoc{
				OwnerID:     "user_id",
				DisplayName: "display_name",
				Email:       "email",
				UpdateDate:  1569152551,
				CreateDate:  1569152551,
				DisplayPic:  "url",
			})
			reflect.ValueOf(result).Elem().Set(reflect.ValueOf(profiles))
			return nil
		}).Times(1)
		t.mockMongo.EXPECT().Find(context.Background(), bson.M{
			"owner_id": bson.M{
				"$in": []string{
					"1", "2", "3", "user_id",
				}},
		}).Return(t.cursor, nil)
		profile, err := t.service.GetProfileList(context.Background(), []string{
			"1", "2", "3", "user_id"})
		t.NoError(err)
		t.NotNil(profile)
		t.Equal("user_id", profile[0].OwnerID)
		t.Equal("display_name", profile[0].DisplayName)
		t.Equal("email", profile[0].Email)
		t.Equal(int64(1569152551), profile[0].UpdateDate)
		t.Equal(int64(1569152551), profile[0].CreateDate)
		t.Equal("url", profile[0].DisplayPic)
	})

}
