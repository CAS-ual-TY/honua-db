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

type HonuaDB struct {
	mongoDB     *mongo.Database
	psqlDB      *sql.DB
	pathToFiles string
}

var instance *HonuaDB

func GetInstance(dbname, psqlUser, psqlPwd, mongoUser, mongoPwd, psqlHost, psqlPort, mongoHost, mongoPort, pathToFiles string) *HonuaDB {
	if instance == nil {
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

		err = instance.create_tables()
		if err != nil {
			panic(err) // Bei einem Fehler wird ein Panic ausgelöst
		}
		instance.Migrate()
	}

	return instance
}

func (hdb *HonuaDB) create_tables() error {
	stmts, err := readAndParseSqlFile(fmt.Sprintf("%s/create.sql", hdb.pathToFiles))
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
