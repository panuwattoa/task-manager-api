package profile

import (
	"context"
	iMongo "task-manager-api/internal/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate mockgen -source=./profile.go -destination=./mock/profile_mock.go
type IMongo interface {
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) iMongo.SingleResult
}

type ProfileDoc struct {
	ID          string `json:"id" bson:"_id"`
	OwnerID     string `json:"owner_id" bson:"owner_id"`
	DisplayName string `json:"display_name" bson:"display_name"`
	Email       string `json:"email" bson:"email"`
	UpdateDate  int64  `bson:"update_date"`
	CreateDate  int64  `bson:"create_date"`
}

type Profile struct {
	mongo IMongo
}

func NewProfileService(mongo IMongo) *Profile {
	return &Profile{mongo: mongo}
}

func (p *Profile) GetProfile(ctx context.Context, id string) (*ProfileDoc, error) {
	return nil, nil
}
