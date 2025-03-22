package mongodb

import (
	"context"
	"fmt"
	"github.com/kohirens/www/session"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StorageDocument struct {
	collection *mongo.Collection
}

func NewStorageMongoDB(c *mongo.Client, database, collection string) *StorageDocument {
	return &StorageDocument{
		collection: c.Database(database).Collection(collection),
	}
}

func (sd *StorageDocument) Save(data *session.Data) error {
	query := map[string][]byte{"session_id": []byte(data.Id)}

	_, e1 := UpsertOne(
		query,
		data,
		sd.collection,
	)
	if e1 != nil {
		return fmt.Errorf(stderr.CannotSaveSession, e1.Error())
	}

	return nil
}

func (sd *StorageDocument) Load(id string) ([]byte, error) {
	query := bson.M{"session_id": id}
	result := sd.collection.FindOne(context.TODO(), query)

	bsonDoc, e1 := result.Raw()
	if e1 != nil {
		return nil, fmt.Errorf(stderr.CannotLoadSession, e1.Error())
	}

	byteAry, e2 := bson.Marshal(bsonDoc)
	if e2 != nil {
		return nil, fmt.Errorf(stderr.CannotLoadSession, e2.Error())
	}
	return byteAry, nil
}
