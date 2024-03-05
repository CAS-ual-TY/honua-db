package honuadb

import (
	"context"
	"github.com/JonasBordewick/honua-db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (hdb *HonuaDB) AddEntityInformation(information *models.EntityInformation) error {
	_, err := hdb.mongoDB.Collection("information").InsertOne(context.Background(), information)
	return err
}

func (hdb *HonuaDB) GetEntityInformation(id string) (*models.EntityInformation, error) {
	filter := bson.M{"_id": id}
	var result *models.EntityInformation
	err := hdb.mongoDB.Collection("information").FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (hdb *HonuaDB) DeleteEntityInformation(id string) error {
	filter := bson.M{"_id": id}
	_, err := hdb.mongoDB.Collection("information").DeleteOne(context.Background(), filter)
	return err
}

func (hdb *HonuaDB) ExistEntityInformation(id string) (bool, error) {
	filter := bson.M{"_id": id}
	var result *models.EntityHistory
	err := hdb.mongoDB.Collection("information").FindOne(context.Background(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
