package honuadb

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HonuaDB struct {
	mongoDB     *mongo.Database
	psqlDB      *sql.DB
	pathToFiles string
}

var instance *HonuaDB

func GetInstance(dbname, psqlUser, psqlPwd, mongoUser, mongoPwd, psqlHost, psqlPort, mongoHost, mongoPort, pathToFiles string) *HonuaDB {
	if instance == nil {
		// Connect to PSQL Database first
		var connStr = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", psqlUser, psqlPwd, psqlHost, psqlPort, dbname)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}
		if err = db.Ping(); err != nil {
			log.Fatal(err)
		}
		log.Println("The PSQL Database connection is established")
		instance = &HonuaDB{
			psqlDB:      db,
			pathToFiles: pathToFiles,
		}
		// Connect to MongoDB
		clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", mongoUser, mongoPwd, mongoHost, mongoPort))

		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}
		err = client.Ping(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("The Mongo Database connection is established")

		database := client.Database(dbname)
		
		instance.mongoDB = database

		// Create PSQL Tables + Migrate all files
		err = instance.create_tables()
		if err != nil {
			panic(err) // If any error occure Panic
		}
		instance.Migrate()
	}

	return instance
}

func (hdb *HonuaDB) create_tables() error {
	stmts, err := read_and_parse_sql_file(fmt.Sprintf("%s/create.sql", hdb.pathToFiles))
	if err != nil {
		return err
	}
	for _, stmt := range stmts {
		_, err := hdb.psqlDB.Exec(stmt)
		if err != nil {
			log.Printf("Error while executing statement %s: %s\n", stmt, err.Error())
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
