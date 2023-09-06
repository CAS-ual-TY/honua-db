package honuadb

import (
	"errors"

	"github.com/JonasBordewick/honua-db/models"
)

func (hdb *HonuaDB) AddDelay(delay *models.Delay, hasID bool) (int, error) {
	const query = "INSERT INTO delays(id, identity, hours, minutes, seconds) VALUES($1, $2, $3, $4, $5);"
	id, err := hdb.get_delay_id(delay.Identity)
	if err != nil {
		return -1, err
	}
	if hasID {
		id = int(delay.ID)
	}

	_, err = hdb.psqlDB.Exec(query, id, delay.Identity, delay.Hours, delay.Minutes, delay.Seconds)
	return id, err
}

func (hdb *HonuaDB) GetDelay(delayID int, identity string) (*models.Delay, error) {
	const query = "SELECT * from delays WHERE id=$1 AND identity=$2;"
	rows, err := hdb.psqlDB.Query(query, delayID, identity)
	if err != nil {
		return nil, err
	}

	var delay *models.Delay;

	for rows.Next() {
		var id int
		var identity string
		var hours int32
		var minutes int32
		var seconds int32

		err := rows.Scan(&id, &identity, &hours, &minutes, &seconds)
		if err != nil {
			return nil, err
		}
		delay = &models.Delay{
			ID: id,
			Identity: identity,
			Hours: hours,
			Minutes: minutes,
			Seconds: seconds,
		}
	}

	rows.Close()

	if delay == nil {
		return nil, errors.New("something went wrong")
	}

	return delay, nil
} 

func (hdb *HonuaDB) EditDelay(delay *models.Delay) error {
	const query = "UPDATE delays SET hours=$1, minutes=$2, seconds=$3 WHERE id=$4 AND identity=$5;"
	_, err := hdb.psqlDB.Exec(query, delay.Hours, delay.Minutes, delay.Seconds, delay.ID, delay.Identity)
	return err
}

func (hdb *HonuaDB) DeleteDelay(delayID int, identity string) error {
	const query = "DELETE FROM delays WHERE id=$1 AND identity=$2;"
	_, err := hdb.psqlDB.Exec(query, delayID, identity)
	return err
}

func (hdb *HonuaDB) get_delay_id(identity string) (int, error) {
	query := "SELECT CASE WHEN EXISTS ( SELECT * FROM delays WHERE identity = $1) THEN true ELSE false END"

	rows, err := hdb.psqlDB.Query(query, identity)
	if err != nil {
		return -1, err
	}

	var exist_identity bool = false

	for rows.Next() {
		err = rows.Scan(&exist_identity)
		if err != nil {
			rows.Close()
			return -1, err
		}
	}

	rows.Close()

	if !exist_identity {
		return 0, nil
	}

	query = "SELECT MAX(id) FROM delays WHERE identity = $1;"

	rows, err = hdb.psqlDB.Query(query, identity)
	if err != nil {
		return -1, err
	}

	var id int = -1

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			rows.Close()
			return -1, err
		}
	}
	rows.Close()

	if id == -1 {
		return -1, errors.New("something went wrong during getting id of entity")
	}

	id = id + 1

	return id, nil
}
