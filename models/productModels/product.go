package productModels

import "github.com/kamva/mgm/v3"

type Product struct {
	mgm.DefaultModel `json:",omitempty"  bson:",inline"`
	Class            string   `json:",omitempty"  bson:"_cls"`
	Shop             string   `json:"shop"        bson:"shop"`
	Articul          string   `json:"articul"     bson:"articul"`
	Url              string   `json:"url"         bson:"url"`
	Title            string   `json:"title"       bson:"title"`
	Category         int32    `json:"category"    bson:"category"`
	Subcategory      int32    `json:"subcategory" bson:"subcategory"`
	Sex              string   `json:"sex"         bson:"sex"`
	Images           []string `json:"images"      bson:"images"`
	Old_price        int      `json:"old_price"   bson:"old_price"`
	New_price        int      `json:"new_price"   bson:"new_price"`
	Discount         float64  `json:"discount"    bson:"discount"`
	Sizes            []string `json:"sizes"       bson:"sizes"`
	Views            int32    `json:"views"       bson:"views"`
}
