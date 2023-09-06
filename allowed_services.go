package honuadb

func (hdb *HonuaDB) AllowService(identity string, deviceId, serviceId int32) error {

	const query = "INSERT INTO allowed_services(identity, device_id, service_id) VALUES ($1, $2, $3)"

	_, err := hdb.psqlDB.Exec(query, identity, deviceId, serviceId)
	return err
}

func (hdb *HonuaDB) ForbidService(identity string, deviceId, serviceId int32) error {
	const query = "DELETE FROM allowed_services WHERE identity=$1 AND device_id=$2 AND service_id=$3;"

	_, err := hdb.psqlDB.Exec(query, identity, deviceId, serviceId)

	return err
}

func (hdb *HonuaDB) IsServiceAllowed(identity string, deviceId, serviceId int32) (bool, error) {
	const query = "SELECT CASE WHEN EXISTS ( SELECT * FROM allowed_services WHERE identity = $1 AND device_id = $2 AND service_id = $3) THEN true ELSE false END"

	rows, err := hdb.psqlDB.Query(query, identity, deviceId, serviceId)
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

func (hdb *HonuaDB) GetIDsOfAllowedServices(identity string, deviceId int32) ([]int, error) {
	const query = "SELECT service_id FROM allowed_services WHERE identity=$1 AND device_id=$2 ORDER BY id;"
	rows, err := hdb.psqlDB.Query(query, identity, deviceId)
	if err != nil {
		return nil, err
	}

	var ids []int = []int{}

	for rows.Next() {
		var sID int
		err = rows.Scan(&sID)
		if err != nil {
			rows.Close()
			return nil, err
		}

		ids = append(ids, sID)
	}

	rows.Close()

	return ids, nil
}