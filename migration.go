package honuadb

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

const (
	get_metadata    = "SELECT filepath FROM metadata ORDER BY id ASC;"
	add_metadata    = "INSERT INTO metadata(filepath) VALUES ($1);"
	exists_metadata = "SELECT CASE WHEN EXISTS ( SELECT * FROM metadata WHERE filepath = $1) THEN true ELSE false END"
)

func (hdb *HonuaDB) exists_metadata(filepath string) (bool, error) {
	rows, err := hdb.psqlDB.Query(exists_metadata, filepath)
	if err != nil {
		log.Printf("An error occured during checking if the metadata %s exists: %s\n", filepath, err.Error())
		return false, err
	}

	var state bool = false

	for rows.Next() {
		err = rows.Scan(&state)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during checking if the metadata %s exists: %s\n", filepath, err.Error())
			return false, err
		}
	}

	rows.Close()

	return state, nil
}

func (hdb *HonuaDB) get_all_done_migrations() []string {
	var migrations []string

	rows, err := hdb.psqlDB.Query(get_metadata)
	if err != nil {
		rows.Close()
		return migrations
	}

	for rows.Next() {
		var migration string
		err := rows.Scan(&migration)
		if err != nil {
			rows.Close()
			return migrations
		}
		migrations = append(migrations, migration)
	}

	rows.Close()

	return migrations
}

// reads all the migrations sql files which are in the folder /app/database/files
func (hdb *HonuaDB) read_migrations() []string {
	var files []string

	files, err := filepath.Glob(fmt.Sprintf("%s/*.sql", hdb.pathToFiles))
	if err != nil {
		log.Printf("Error running readMigrations %s\n", err.Error())
	}
	return files
}

// adds a migration to the migration table
func (hdb *HonuaDB) write_metadata(migration string) {
	_, err := hdb.psqlDB.Exec(add_metadata, migration)
	if err != nil {
		log.Printf("Error running writeMetadata %s\n", err.Error())
	}
}

// public Method to start the Migration
func (hdb *HonuaDB) Migrate() {
	// get all migrations that where done in the past
	var done []string = hdb.get_all_done_migrations()
	// get all migrations which are in the folder /app/database/files
	var migrations []string = hdb.read_migrations()

	// array to store all the migrations that were not already done in the past
	var todo []string = []string{}

	// Iterate through the migration list and add those migrations that have not already been done to the todo list
	for _, migration := range migrations {
		if !string_array_contains_string(migration, done) {
			todo = append(todo, migration)
		}
	}

	// Iterate through the todo list and parse the migrations to string statements. And Execute each statement
	// After that write the migration to the metadata table
	for _, migration := range todo {
		if strings.Contains(migration, "create.sql") {
			log.Printf("Migrate Database skip file %s\n", migration)
			continue
		}
		log.Printf("Migrate Database with file %s\n", migration)
		stmts, err := read_and_parse_sql_file(migration)
		if err != nil {
			log.Printf("Error while Migrating with file %s: %s\n", migration, err.Error())
			continue
		}
		for _, stmt := range stmts {
			_, err := hdb.psqlDB.Exec(stmt)
			if err != nil {
				log.Printf("Error while Migrating with file %s: %s\n", migration, err.Error())
				continue
			}
		}
		hdb.write_metadata(migration)
	}
}
