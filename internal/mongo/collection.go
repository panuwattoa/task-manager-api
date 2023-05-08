package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollectionHelper struct {
	collection *mongo.Collection
}

func NewCollectionHelper(c *mongo.Collection) *CollectionHelper {
	return &CollectionHelper{collection: c}
}

func (c *CollectionHelper) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResult {
	return c.collection.FindOne(ctx, filter, opts...)
}

func (c *CollectionHelper) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.collection.UpdateOne(ctx, filter, update, opts...)
}

func (c *CollectionHelper) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return c.collection.InsertOne(ctx, document, opts...)
}
