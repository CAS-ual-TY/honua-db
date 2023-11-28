package honuadb

import (
	"context"

	"github.com/JonasBordewick/honua-db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (hdb *HonuaDB) AddDashboard(dashboard *models.Dashboard) error {
	exist, err := hdb.ExistDashboard(dashboard.ID)
	if err != nil {
		return err
	}

	if exist {
		err = hdb.DeleteDashboard(dashboard.ID)
		if err != nil {
			return err
		}
	}
	_, err = hdb.mongoDB.Collection("dashboard").InsertOne(context.Background(), dashboard)
	return err
}

func (hdb *HonuaDB) GetDashboard(id string) (*models.Dashboard, error) {
	filter := bson.M{"_id": id}
	var result *models.Dashboard
	err := hdb.mongoDB.Collection("dashboard").FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (hdb *HonuaDB) DeleteDashboard(id string) error {
	filter := bson.M{"_id": id}
	_, err := hdb.mongoDB.Collection("dashboard").DeleteOne(context.Background(), filter)
	return err
}

func (hdb *HonuaDB) ExistDashboard(id string) (bool, error) {
	filter := bson.M{"_id": id}
	var result *models.Dashboard
	err := hdb.mongoDB.Collection("dashboard").FindOne(context.Background(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
