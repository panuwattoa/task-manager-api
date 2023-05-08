package mongo

//go:generate mockgen -source=./single_result.go -destination=./mock/single_result.go
type SingleResult interface {
	Decode(v interface{}) error
}
