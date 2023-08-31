package main

import (
	hdb "github.com/JonasBordewick/honua-db"
	"github.com/JonasBordewick/honua-db/models"
)

func main() {
	db := hdb.GetInstance("honua_private", "admin", "pwd4Database", "admin", "pwd4Database", "192.168.0.138", "5432", "192.168.0.138", "27017", "./files")
	db.AddServices("test", []*models.Service{&models.Service{Domain: "Test"}})
}