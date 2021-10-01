package filtersAndOptions

import (
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
)

type Filter struct {
	f bson.M
}

func NewFilter() *Filter {
	return &Filter{f: bson.M{}}
}

func (f *Filter) And(values ...bson.M) {
	_, exist := f.f[operator.And]
	if exist == false {
		f.f[operator.And] = []bson.M{}
	}

	for _, value := range values {
		f.f[operator.And] = append(f.f[operator.And].([]bson.M), value)
	}
}

func (f *Filter) Or(values ...bson.M) {
	_, exist := f.f[operator.Or]
	if exist == false {
		f.f[operator.Or] = []bson.M{}
	}

	for _, value := range values {
		f.f[operator.Or] = append(f.f[operator.Or].([]bson.M), value)
	}
}

func (f *Filter) Greater(key string, value interface{}) bson.M {
	return bson.M{key: bson.M{operator.Gte: value}}
}

func (f *Filter) Less(key string, value interface{}) bson.M {
	return bson.M{key: bson.M{operator.Lte: value}}
}

func (f *Filter) Equal(key string, value interface{}) bson.M {
	return bson.M{key: bson.M{operator.Eq: value}}
}

func (f *Filter) Sort(_ bson.D) {}

func (f *Filter) Page(_ *int64) {}

func (f *Filter) ToBson() bson.M {
	return f.f
}
