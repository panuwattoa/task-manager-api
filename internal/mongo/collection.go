package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoPaginate struct {
	limit int64
	page  int64
}

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

func (c *CollectionHelper) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (Cursor, error) {
	return c.collection.Find(ctx, filter, opts...)
}

func NewMongoPaginate(limit, page int) *mongoPaginate {
	return &mongoPaginate{
		limit: int64(limit),
		page:  int64(page),
	}
}

func (mp *mongoPaginate) GetPaginatedOpts() *options.FindOptions {
	l := mp.limit
	skip := mp.page*mp.limit - mp.limit
	fOpt := options.FindOptions{Limit: &l, Skip: &skip}

	return &fOpt
}
