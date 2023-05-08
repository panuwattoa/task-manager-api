package mongo

import (
	"context"
	"fmt"
	cfg "task-manager-api/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type MongoDB struct {
	client *mongo.Client
}

func NewMongoDB() *MongoDB {
	return new(MongoDB)
}

func (m *MongoDB) Open(ctx context.Context) error {
	clientOptions := options.Client()
	clientOptions.SetAppName(cfg.Conf.MongoDB.AppName)
	clientOptions.SetConnectTimeout(cfg.Conf.MongoDB.Timeout * time.Second)
	clientOptions.SetAuth(options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		AuthSource:    cfg.MongoDBName,
		Username:      cfg.MongoUser,
		Password:      cfg.MongoPassword,
	})

	ctx, cancel := context.WithTimeout(ctx, cfg.Conf.MongoDB.DefaultContextTimeout*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions.ApplyURI(cfg.MongoHost))
	if err != nil {
		fmt.Printf("cannot connect to MongoDB (%v), %v", cfg.MongoDBName, err)
		return err
	}

	m.client = client
	return nil
}

func (m *MongoDB) Close(ctx context.Context) error {
	err := m.Status(ctx)
	if err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Conf.MongoDB.DefaultContextTimeout*time.Second)
		defer cancel()

		err = m.client.Disconnect(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MongoDB) Status(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Conf.MongoDB.DefaultContextTimeout*time.Second)
	defer cancel()
	if err := m.client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	collOpts := []*options.CollectionOptions{
		options.Collection().SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
	}
	return m.client.Database(cfg.MongoDBName).Collection(name, collOpts...)
}
