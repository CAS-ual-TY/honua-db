package honuadb

import "github.com/JonasBordewick/honua-db/models"

func (hdb *HonuaDB) AddIdentity(identity *models.Identity) error {
	const query = "INSERT INTO identities(id, name) VALUES($1, $2);"
	_, err := hdb.psqlDB.Exec(query, identity.ID, identity.Name)
	return err
}

func (hdb *HonuaDB) DeleteIdentity(id string) error {
	const query = "DELETE FROM identities WHERE id = $1"
	_, err := hdb.psqlDB.Exec(query, id)
	return err
}

func (hdb *HonuaDB) ExistIdentity(identifier string) (bool, error) {
	const query = "SELECT CASE WHEN EXISTS ( SELECT * FROM identities WHERE id = $1) THEN true ELSE false END;"
	rows, err := hdb.psqlDB.Query(query, identifier)
	if err != nil {
		return false, err
	}

	var state bool = false

	for rows.Next() {
		err = rows.Scan(&state)
		if err != nil {
			rows.Close()
			return false, err
		}
	}

	rows.Close()

	return state, nil
}