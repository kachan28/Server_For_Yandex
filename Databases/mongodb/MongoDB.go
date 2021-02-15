package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func setConnection() (context.Context, *mongo.Client, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	return ctx, client, cancel, err
}

//GetItems function to get multiple items from mongo database
func GetItems(filter bson.M, options *options.FindOptions, databases []string) (*mongo.Client, context.Context, []*mongo.Cursor, context.CancelFunc, error) {
	ctx, client, cancel, err := setConnection()
	var result []*mongo.Cursor
	for _, database := range databases {
		collection := client.Database(database).Collection("series")
		cur, err := collection.Find(ctx, filter, options)
		if err != nil {
			panic(err)
		}
		result = append(result, cur)
	}
	return client, ctx, result, cancel, err
}
