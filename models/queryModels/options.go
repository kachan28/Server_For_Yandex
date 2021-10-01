package queryModels

import (
	"discountDealer/repository/productRepository/filtersAndOptions"
	"strings"
)

type FindQuerySort struct {
	Root
	Price    *string `query:"price_sort"`
	Views    *string `query:"views_sort"`
	Title    *string `query:"title_sort"`
	Discount *string `query:"discount_sort"`
}

func (f *FindQuerySort) parse(value interface{}, tagName string, res filtersAndOptions.ProductOpt) {
	sortColumn := tagName[:strings.LastIndex(tagName, "_")]
	res.Sort(filtersAndOptions.SortColumn(sortColumn, sortDirs[*value.(*string)]))
}

type FindQueryPage struct {
	Root
	Page *int64 `query:"page"`
}

func (f *FindQueryPage) parse(value interface{}, _ string, res filtersAndOptions.ProductOpt) {
	res.Page(value.(*int64))
}
