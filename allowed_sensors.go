package honuadb

func (hdb *HonuaDB) AllowSensor(identity string, deviceId, sensorId int32) error {
	const query = "INSERT INTO allowed_sensors(identity, device_id, sensor_id) VALUES ($1, $2, $3);"
	_, err := hdb.psqlDB.Exec(query, identity, deviceId, sensorId)
	return err
}

func (hdb *HonuaDB) ForbidSensor(identity string, deviceId, sensorId int32) error {
		const query = "DELETE FROM allowed_sensors WHERE identity=$1 AND device_id=$2 AND sensor_id=$3;"

	_, err := hdb.psqlDB.Exec(query, identity, deviceId, sensorId)
	return err
}

func (hdb *HonuaDB) IsSensorAllowed(identity string, deviceId, sensorId int32) (bool, error) {
	const query = "SELECT CASE WHEN EXISTS ( SELECT * FROM allowed_sensors WHERE identity=$1 AND device_id = $2 AND sensor_id = $3) THEN true ELSE false END"

	rows, err := hdb.psqlDB.Query(query, identity, deviceId, sensorId)
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