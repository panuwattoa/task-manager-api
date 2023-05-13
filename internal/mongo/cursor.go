package mongo

import "context"

//go:generate mockgen -source=./cursor.go -destination=./mock/cursor.go
// Cursor is an interface for `mongo.Cursor` structure
// Documentation: https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Cursor
type Cursor interface {
	All(ctx context.Context, results interface{}) error
	Close(ctx context.Context) error
	Decode(val interface{}) error
	Err() error
	ID() int64
	Next(ctx context.Context) bool
	RemainingBatchLength() int
	TryNext(ctx context.Context) bool
}
