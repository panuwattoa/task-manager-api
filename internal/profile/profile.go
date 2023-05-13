package profile

import (
	"context"
	"errors"
	iMongo "task-manager-api/internal/mongo"
	m "task-manager-api/internal/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate mockgen -source=./profile.go -destination=./mock/profile.go
type IMongo interface {
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) iMongo.SingleResult
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (m.Cursor, error)
}

type ProfileDoc struct {
	OwnerId     string `json:"owner_id" bson:"owner_id"`
	DisplayName string `json:"display_name" bson:"display_name"`
	Email       string `json:"email" bson:"email"`
	DisplayPic  string `json:"display_pic" bson:"display_pic"`
	UpdateDate  int64  `json:"-" bson:"update_date"`
	CreateDate  int64  `json:"-" bson:"create_date"`
}

type Profile struct {
	mongo IMongo
	time  func() time.Time
}

func NewProfileService(mongo IMongo) *Profile {
	return &Profile{mongo: mongo}
}

func (p *Profile) GetProfile(ctx context.Context, ownerId string) (*ProfileDoc, error) {
	result := p.mongo.FindOne(ctx, bson.M{"owner_id": ownerId})
	profile := new(ProfileDoc)
	if err := result.Decode(profile); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return profile, nil
}

func (p *Profile) GetProfileList(ctx context.Context, ownerId []string) ([]ProfileDoc, error) {
	curr, err := p.mongo.Find(ctx, bson.M{"owner_id": bson.M{"$in": ownerId}})

	if err != nil {
		// TODO: log error
		return nil, err
	}
	var profiles = make([]ProfileDoc, 0)
	if err := curr.All(ctx, &profiles); err != nil {
		return nil, err
	}
	return profiles, nil
}

func (p *Profile) now() time.Time {
	if p.time == nil {
		return time.Now()
	}

	return p.time()
}
