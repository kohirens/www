// Package mongodb
// Deprecated see github.com/kohirens/mongodb standalone library. This will be
// removed in the next major release.
package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

const (
	ConnectionEnvVar = "MONGODB_CONNECTION"
)

func Connection() (*mongo.Client, error) {
	dbConnStr, ok3 := os.LookupEnv(ConnectionEnvVar)
	if !ok3 {
		return nil, fmt.Errorf(stderr.EnvVarUnset, ConnectionEnvVar)
	}

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dbConnStr).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, e1 := mongo.Connect(context.TODO(), opts)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.Connect, e1.Error())
	}

	return client, nil
}

func InsertOne(doc interface{}, database, collection string, c *mongo.Client) (*mongo.InsertOneResult, error) {
	coll := c.Database(database).Collection(collection)

	result, e1 := coll.InsertOne(context.TODO(), doc)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.CannotInsertData, database, collection, e1.Error())
	}

	return result, nil
}

// UpsertOne Update an existing document or insert when it cannot be found.
func UpsertOne(query interface{}, doc interface{}, collection *mongo.Collection, hint ...interface{}) (*mongo.UpdateResult, error) {
	truthy := true // because they require a pointer instead of a copy.
	opts := &options.UpdateOptions{
		Upsert: &truthy,
	}

	if len(hint) > 0 {
		switch hint[0].(type) {
		case string:
			opts.Hint = hint[0]
		}
	}

	updateDoc := bson.M{
		"$set": doc,
	}

	result, e1 := collection.UpdateOne(context.TODO(), query, updateDoc, opts)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.CannotUpsertData, e1.Error())
	}

	return result, nil
}
