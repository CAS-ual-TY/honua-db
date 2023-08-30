package honuadb

import (
	"database/sql"
	"errors"
	"log"

	"github.com/JonasBordewick/honua-db/models"
)

func (hdb *HonuaDB) AddEntity(entity *models.Entity, hasID bool, id int32) error {
	const query = `
INSERT INTO entities(
	id, identity, entity_id, name,
	is_device, allow_rules, has_attribute,
	attribute, is_victron_sensor, sensor_type, has_numeric_state
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
`
	var attributeString sql.NullString = sql.NullString{
		Valid:  entity.Attribute != "",
		String: entity.Attribute,
	}

	var err error

	if !hasID {
		id, err = hdb.get_entity_id(entity.Identity)
		if err != nil {
			return err
		}
	}

	_, err = hdb.psqlDB.Exec(
		query, id, entity.Identity, entity.EntityID, entity.Name,
		entity.IsDevice, entity.AllowRules, entity.HasAttribute,
		attributeString, entity.IsVictronSensor, entity.SensorType,
		entity.HasNumericState,
	)

	return err
}

func (hdb *HonuaDB) GetEntity(identity string, id int) (*models.Entity, error) {
	const query = "SELECT * FROM entities WHERE id=$1 AND identity=$2;"

	rows, err := hdb.psqlDB.Query(query, id, identity)
	if err != nil {
		return nil, err
	}

	var result *models.Entity

	for rows.Next() {
		entity, err := hdb.make_entity(rows)
		if err != nil {
			rows.Close()
			return nil, err
		}
		result = entity
	}

	rows.Close()

	return result, nil
}

func (hdb *HonuaDB) GetEntities(identity string) ([]*models.Entity, error) {
	const query = "SELECT * FROM entities WHERE identity = $1;"

	rows, err := hdb.psqlDB.Query(query, identity)
	if err != nil {
		return nil, err
	}

	var entities []*models.Entity = []*models.Entity{}

	for rows.Next() {
		entity, err := hdb.make_entity(rows)
		if err != nil {
			rows.Close()
			return nil, err
		}
		entities = append(entities, entity)
	}

	rows.Close()

	return entities, nil
}

func (hdb *HonuaDB) EditEntity(identity string, entity *models.Entity) error {
	const query = `
UPDATE entities
SET name = $1, is_device = $2, allow_rules = $3, has_attribute = $4, attribute = $5, is_victron_sensor = $6, sensor_type = $7, has_numeric_state = $8
WHERE identity = $9 AND id = $10;
	`

	var attributeString sql.NullString = sql.NullString{
		Valid:  entity.Attribute != "",
		String: entity.Attribute,
	}

	entity.HasAttribute = attributeString.Valid

	_, err := hdb.psqlDB.Exec(query, entity.Name, entity.IsDevice, entity.AllowRules, entity.HasAttribute, attributeString, entity.IsVictronSensor, entity.SensorType, entity.HasNumericState, entity.Identity, entity.ID)

	if err != nil {
		log.Printf("An error occured during editity entitiy: %s\n", err.Error())
	}
	return err
}

func (hdb *HonuaDB) DeleteEntity(identity string, id int32) error {
	const query = "DELETE FROM entities WHERE identity=$1 AND id = $2;"

	_, err := hdb.psqlDB.Exec(query, identity, id)
	if err != nil {
		log.Printf("An error occured during deleting the entity with id = %d: %s\n", id, err.Error())
	}
	return err
}

func (hdb *HonuaDB) ExistEntity(identity string, id int32) (bool, error) {

	const query = "SELECT CASE WHEN EXISTS ( SELECT * FROM entities WHERE identity = $1 AND id = $2) THEN true ELSE false END"

	rows, err := hdb.psqlDB.Query(query, identity, id)
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

func (hdb *HonuaDB) make_entity(rows *sql.Rows) (*models.Entity, error) {
	var id int
	var identity string
	var entityID string
	var name string
	var isDevice bool
	var allowRules bool
	var hasAttribute bool
	var attribute sql.NullString
	var isVictronSensor bool
	var sensorType models.VictronSensorType
	var hasNumericState bool

	err := rows.Scan(&id, &identity, &entityID, &name, &isDevice, &allowRules, &hasAttribute, &attribute, &isVictronSensor, &sensorType,  &hasNumericState)
	if err != nil {
		return nil, err
	}

	var result *models.Entity = &models.Entity{
		ID:              int32(id),
		Identity:      identity,
		Name:            name,
		EntityID:        entityID,
		IsDevice:        isDevice,
		AllowRules:      allowRules,
		HasAttribute:    hasAttribute,
		IsVictronSensor: isVictronSensor,
		SensorType: sensorType,
		HasNumericState: hasNumericState,
	}

	if hasAttribute && attribute.Valid {
		result.Attribute = attribute.String
	} else {
		result.HasAttribute = false
		result.Attribute = ""
	}
	return result, nil
}

func (hdb *HonuaDB) get_entity_id(identifier string) (int32, error) {
	query := "SELECT CASE WHEN EXISTS ( SELECT * FROM entities WHERE identity = $1) THEN true ELSE false END"

	rows, err := hdb.psqlDB.Query(query, identifier)
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

	query = "SELECT MAX(id) FROM entities WHERE identity = $1;"

	rows, err = hdb.psqlDB.Query(query, identifier)
	if err != nil {
		return -1, err
	}

	var id int32 = -1

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
