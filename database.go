package honuadb

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/lib/pq"
)

// HonuaDB repräsentiert eine Verbindung zu PostgreSQL und MongoDB Datenbanken.
type HonuaDB struct {
	mongoDB     *mongo.Database
	psqlDB      *sql.DB
	pathToFiles string
}

// instance ist eine Singleton-Instanz von HonuaDB.
var instance *HonuaDB

// GetInstance gibt eine Instanz von HonuaDB zurück. Wenn keine Instanz vorhanden ist, wird eine erstellt.
func GetInstance(dbname, psqlUser, psqlPwd, mongoUser, mongoPwd, psqlHost, psqlPort, mongoHost, mongoPort, pathToFiles string) *HonuaDB {
	if instance == nil {
		// Verbindung zur PSQL-Datenbank herstellen
		var connStr = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", psqlUser, psqlPwd, psqlHost, psqlPort, dbname)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}
		if err = db.Ping(); err != nil {
			log.Fatal(err)
		}
		log.Println("Die Verbindung zur PSQL-Datenbank wurde hergestellt")
		instance = &HonuaDB{
			psqlDB:      db,
			pathToFiles: pathToFiles,
		}
		// Verbindung zu MongoDB herstellen
		clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", mongoUser, mongoPwd, mongoHost, mongoPort))

		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}
		err = client.Ping(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Die Verbindung zur MongoDB wurde hergestellt")

		database := client.Database(dbname)

		instance.mongoDB = database

		// PSQL-Tabellen erstellen + alle Dateien migrieren
		err = instance.create_tables()
		if err != nil {
			panic(err) // Bei einem Fehler wird ein Panic ausgelöst
		}
		instance.Migrate()
	}

	return instance
}

// create_tables erstellt PSQL-Tabellen und migriert alle Dateien.
func (hdb *HonuaDB) create_tables() error {
	stmts, err := read_and_parse_sql_file(fmt.Sprintf("%s/create.sql", hdb.pathToFiles))
	if err != nil {
		return err
	}
	for _, stmt := range stmts {
		_, err := hdb.psqlDB.Exec(stmt)
		if err != nil {
			log.Printf("Fehler beim Ausführen des Statements %s: %s\n", stmt, err.Error())
			return err
		}
	}

	exists, err := hdb.exists_metadata(fmt.Sprintf("%s/create.sql", hdb.pathToFiles))
	if err != nil {
		return err
	}

	if !exists {
		hdb.write_metadata(fmt.Sprintf("%s/create.sql", hdb.pathToFiles))
	}

	return nil
}
