package queryModels

import (
	"discountDealer/repository/productRepository/filtersAndOptions"
	"discountDealer/repository/searchRepository"
	"go.uber.org/zap"
	"strings"
)

type FindQuery struct {
	Root
	Title *string `query:"title"`
}

func (f *FindQuery) parse(value interface{}, tagName string, res filtersAndOptions.ProductOpt) {
	articuls, err := searchRepository.Init().FindProduct(value.(*string))
	if err != nil {
		f.Log().Error("Error while searching by title", zap.Error(err))
	}
	for _, articul := range articuls {
		res.Or(filtersAndOptions.NewFilter().Equal("articul", articul))
	}
}

type FindQueryEqFilter struct {
	Root
	Category    *int    `query:"category"`
	SubCategory *int    `query:"subcategory"`
	Sex         *string `query:"sex"`
}

func (f FindQueryEqFilter) parse(value interface{}, tagName string, res filtersAndOptions.ProductOpt) {
	res.And(filtersAndOptions.NewFilter().Equal(tagName, value))
}

type FindQueryCompFilter struct {
	Root
	PriceFrom    *int     `query:"new_price_from"`
	PriceTo      *int     `query:"new_price_to"`
	DiscountFrom *float64 `query:"discount_from"`
	DiscountTo   *float64 `query:"discount_to"`
}

func (f FindQueryCompFilter) parse(value interface{}, tagName string, res filtersAndOptions.ProductOpt) {
	filterColumn := tagName[:strings.LastIndex(tagName, "_")]
	filterDir := tagName[strings.LastIndex(tagName, "_")+1:]
	if filterDir == "to" {
		res.And(filtersAndOptions.NewFilter().Less(filterColumn, value))
	} else {
		res.And(filtersAndOptions.NewFilter().Greater(filterColumn, value))
	}
}
