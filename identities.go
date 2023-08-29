package honuadb

import "github.com/JonasBordewick/honua-db/models"

func (hdb *HonuaDB) AddIdentity(identity *models.Identity) error {
	const query = "INSERT INTO identities(id, name) VALUES($1, $2);"
	_, err := hdb.psqlDB.Exec(query, identity.ID, identity.Name)
	return err
}

func (hdb *HonuaDB) DeleteIdentity(id string) error {
	const query = "DELETE FROM identities WHERE identifier = $1"
	_, err := hdb.psqlDB.Exec(query, id)
	return err
}