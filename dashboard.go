package honuadb

import (
	"context"

	"github.com/JonasBordewick/honua-db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddDashboard fügt zur Datenbank ein (neues) Dashboard hinzu. Falls bereits
// ein Dashboard mit der ID existiert, dann wird das bereits existierende
// Dashboard gelöscht und erst dann wird das neue hinzugefügt.
// Die Dashboard ID ist die identity des entsprechenden backends.
// Die Methode ist ADD + EDIT gleichzeitig
func (hdb *HonuaDB) AddDashboard(dashboard *models.Dashboard) error {
	exist, err := hdb.exists_dashboard(dashboard.ID)
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

// GetDashboard gibt das Dashboard mit der angegebenen ID zurück
func (hdb *HonuaDB) GetDashboard(id string) (*models.Dashboard, error) {
	filter := bson.M{"_id": id}
	var result *models.Dashboard
	err := hdb.mongoDB.Collection("dashboard").FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteDashboard löscht das Dashboard mit der angegebenen ID
func (hdb *HonuaDB) DeleteDashboard(id string) error {
	filter := bson.M{"_id": id}
	_, err := hdb.mongoDB.Collection("dashboard").DeleteOne(context.Background(), filter)
	return err
}

// Die Private Methode exists_dashboard checkt, ob bereits ein Dashboard unter
// der angegebenen ID existiert. Und gibt true zurück falls diese bereits
// existiert und false wenn es kein Dashboard mit der angegebenen ID gibt.
func (hdb *HonuaDB) exists_dashboard(id string) (bool, error) {
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
