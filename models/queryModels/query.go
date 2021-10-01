package queryModels

import (
	"discountDealer/repository/productRepository/filtersAndOptions"
	"go.uber.org/zap"
	"reflect"
)

const structIterationStart = 1

var sortDirs = map[string]int{
	"asc":  filtersAndOptions.Asc,
	"desc": filtersAndOptions.Desc,
}

type root interface {
	Parse(queryOpts interface{})
}

type Root struct {
	log *zap.Logger
}

func NewRoot(logger *zap.Logger) Root {
	return Root{log: logger}
}

func (r Root) Log() *zap.Logger {
	return r.log
}

func (r Root) Parse(queryOpts QueryParse, res filtersAndOptions.ProductOpt) {
	filterValues := reflect.ValueOf(queryOpts).Elem()
	filterTags := reflect.TypeOf(queryOpts).Elem()

	for i := structIterationStart; i < filterValues.NumField(); i++ {
		field := filterValues.Field(i)
		if field.IsNil() {
			continue
		}

		value := field.Interface()
		tagName := filterTags.Field(i).Tag.Get("query")
		queryOpts.parse(value, tagName, res)
	}
}

type QueryParse interface {
	parse(value interface{}, tagName string, res filtersAndOptions.ProductOpt)
}
