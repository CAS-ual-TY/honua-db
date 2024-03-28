package honuadb

import (
	"context"
	"github.com/JonasBordewick/honua-db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (hdb *HonuaDB) AddDetailedDashboard(dashboard *models.DetailedDashboard) error {
	exist, err := hdb.ExistDetailedDashboard(dashboard.ID)
	if err != nil {
		return err
	}

	if exist {
		err = hdb.DeleteDetailedDashboard(dashboard.ID)
		if err != nil {
			return err
		}
	}
	_, err = hdb.mongoDB.Collection("detailed_dashboard").InsertOne(context.Background(), dashboard)
	return err
}

func (hdb *HonuaDB) GetDetailedDashboard(id string) (*models.DetailedDashboard, error) {
	filter := bson.M{"_id": id}
	var result *models.DetailedDashboard
	err := hdb.mongoDB.Collection("detailed_dashboard").FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (hdb *HonuaDB) DeleteDetailedDashboard(id string) error {
	filter := bson.M{"_id": id}
	_, err := hdb.mongoDB.Collection("detailed_dashboard").DeleteOne(context.Background(), filter)
	return err
}

func (hdb *HonuaDB) ExistDetailedDashboard(id string) (bool, error) {
	filter := bson.M{"_id": id}
	var result *models.DetailedDashboard
	err := hdb.mongoDB.Collection("detailed_dashboard").FindOne(context.Background(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
