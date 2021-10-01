package drivers

import "context"

type Driver interface {
	//Insert into database
	Insert(ctx context.Context, entity interface{}) (err error)

	//Get entity from database
	Get(ctx context.Context, entity interface{}) (err error)

	//Update entity in database
	Update(ctx context.Context, entity interface{}) (err error)

	//Delete entity from database
	Delete(ctx context.Context, entity interface{}) (err error)
}
