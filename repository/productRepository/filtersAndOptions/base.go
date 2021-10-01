package filtersAndOptions

import "go.mongodb.org/mongo-driver/bson"

type ProductOpt interface {
	And(values ...bson.M)
	Sort(sort bson.D)
	Page(page *int64)
	Or(values ...bson.M)
}
