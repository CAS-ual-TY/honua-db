package honuadb

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

func (hdb *HonuaDB) exists_metadata(filepath string) (bool, error) {
	rows, err := hdb.psqlDB.Query("SELECT CASE WHEN EXISTS ( SELECT * FROM metadata WHERE filepath = $1) THEN true ELSE false END;", filepath)
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

	rows, err := hdb.psqlDB.Query("SELECT filepath FROM metadata ORDER BY id ASC;")
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

func (hdb *HonuaDB) read_migrations() []string {
	var files []string

	files, err := filepath.Glob(fmt.Sprintf("%s/*.sql", hdb.pathToFiles))
	if err != nil {
		log.Printf("Error running readMigrations %s\n", err.Error())
	}
	return files
}

func (hdb *HonuaDB) write_metadata(migration string) {
	_, err := hdb.psqlDB.Exec("INSERT INTO metadata(filepath) VALUES ($1);", migration)
	if err != nil {
		log.Printf("Error running writeMetadata %s\n", err.Error())
	}
}

// Public Method um die Datenbank Migration zu starten
func (hdb *HonuaDB) Migrate() {
	var done []string = hdb.get_all_done_migrations() // Liste der bereits gemachten Migrations
	var migrations []string = hdb.read_migrations()   // Liste aller SQL-Files im Ordner inkl. create.sql
	var todo []string = []string{}                    // Leere Liste, um alle noch nicht ausgeführten Migrations zu speichern

	// Iteration durch alle SQL-Files, es werden jene zur todo liste hinzugefügt, die nicht bereits ausgeführt wurden
	for _, migration := range migrations {
		if !string_array_contains_string(migration, done) {
			todo = append(todo, migration)
		}
	}

	// Iteration durch alle TODOs
	for _, migration := range todo {
		// Überspringt create.sql
		if strings.Contains(migration, "create.sql") {
			log.Printf("Migrate Database skip file %s\n", migration)
			continue
		}
		log.Printf("Migrate Database with file %s\n", migration)
		stmts, err := readAndParseSqlFile(migration) // liest das SQL-File ein und returnt ein Array an Statements
		if err != nil {
			log.Printf("Error while Migrating with file %s: %s\n", migration, err.Error())
			continue
		}
		// Jedes Statement wird ausgeführt
		for _, stmt := range stmts {
			_, err := hdb.psqlDB.Exec(stmt)
			if err != nil {
				log.Printf("Error while Migrating with file %s: %s\n", migration, err.Error())
				continue
			}
		}
		// die Migration wird in der Datenbank gespeichert, sodass diese beim nächsten mal nicht erneut gemacht wird.
		hdb.write_metadata(migration)
	}
}
