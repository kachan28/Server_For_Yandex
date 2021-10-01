package filtersAndOptions

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	Asc  = 1
	Desc = -1
)

var (
	defaultOffset int64 = 0
	defaultLimit  int64 = 25
	defaultSort         = bson.D{{"_id", Desc}}
)

type Opts struct {
	o *options.FindOptions
}

func NewOpts() *Opts {
	opt := new(options.FindOptions)
	opt.SetSkip(defaultOffset)
	opt.SetLimit(defaultLimit)
	opt.SetSort(defaultSort)
	return &Opts{o: opt}
}

func (o *Opts) Offset(offset int64) {
	o.o.SetSkip(offset)
}

func (o *Opts) Limit(limit int64) {
	o.o.SetLimit(limit)
}

func (o *Opts) Sort(sort bson.D) {
	o.o.SetSort(sort)
}

func (o *Opts) Page(page *int64) {
	o.Offset(*page * defaultLimit)
}

func (o *Opts) Get() *options.FindOptions {
	return o.o
}

func (o *Opts) And(_ ...bson.M) {}

func (o *Opts) Or(_ ...bson.M) {}

func SortColumn(column string, value int) bson.D {
	return bson.D{{column, value}}
}

func (o *Opts) ToBson() *options.FindOptions {
	return o.o
}
