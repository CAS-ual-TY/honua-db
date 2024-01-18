package honuadb

import (
	"context"
	"github.com/JonasBordewick/honua-db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (hdb *HonuaDB) AddHistory(history *models.EntityHistory) error {
	_, err := hdb.mongoDB.Collection("history").InsertOne(context.Background(), history)
	return err
}

func (hdb *HonuaDB) GetHistory(id string) (*models.EntityHistory, error) {
	filter := bson.M{"_id": id}
	var result *models.EntityHistory
	err := hdb.mongoDB.Collection("history").FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (hdb *HonuaDB) DeleteHistory(id string) error {
	filter := bson.M{"_id": id}
	_, err := hdb.mongoDB.Collection("history").DeleteOne(context.Background(), filter)
	return err
}

func (hdb *HonuaDB) ExistHistory(id string) (bool, error) {
	filter := bson.M{"_id": id}
	var result *models.EntityHistory
	err := hdb.mongoDB.Collection("history").FindOne(context.Background(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
