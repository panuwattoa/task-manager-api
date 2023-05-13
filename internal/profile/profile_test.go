package profile

import (
	"context"
	"errors"
	"testing"
	"time"

	mock_sigleResult "task-manager-api/internal/mongo/mock"
	mock_profile "task-manager-api/internal/profile/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

type ProfileTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	mockMongo    *mock_profile.MockIMongo
	service      *Profile
	singleResult *mock_sigleResult.MockSingleResult
}

func (t *ProfileTestSuite) SetupTest() {
	t.ctrl = gomock.NewController(t.T())
	t.mockMongo = mock_profile.NewMockIMongo(t.ctrl)
	t.service = NewProfileService(t.mockMongo)
	t.singleResult = mock_sigleResult.NewMockSingleResult(t.ctrl)
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
	t.Run("get profile success", func() {
		t.singleResult.EXPECT().Decode(&ProfileDoc{}).DoAndReturn(func(doc *ProfileDoc) error {
			doc.OwnerID = "user_id"
			doc.DisplayName = "display_name"
			doc.Email = "email"
			doc.UpdateDate = 1569152551
			doc.CreateDate = 1569152551
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
	})

}
