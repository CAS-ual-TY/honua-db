package honuadb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/JonasBordewick/honua-db/models"
)

// AddCondition fügt eine Bedingung zur Datenbank hinzu. Die Funktion gibt die ID der hinzugefügten Bedingung zurück.
// Wenn hasID true ist, wird die mit der Bedingung verbundene ID verwendet, andernfalls wird eine neue ID generiert.
func (hdb *HonuaDB) AddCondition(condition *models.Condition, hasID bool) (int, error) {
	if condition.Type >= models.NUMERICSTATE {
		return -1, fmt.Errorf("%d is not valid for parent condition", condition.Type)
	}

	id, err := hdb.get_condition_id(condition.Identity)
	if err != nil {
		return -1, err
	}

	if hasID {
		id = condition.ID
	}

	const query = "INSERT INTO conditions(id, identity, condition_type) VALUES($1, $2, $3);"
	_, err = hdb.psqlDB.Exec(query, id, condition.Identity, condition.Type)
	if err != nil {
		return -1, err
	}

	for _, sub := range condition.SubConditions {
		err = hdb.add_subcondition(sub, id, hasID)
		if err != nil {
			return -1, err
		}
	}

	return id, nil
}

// add_subcondition fügt eine Subbedingung zur Datenbank hinzu.
func (hdb *HonuaDB) add_subcondition(condition *models.Condition, parentID int, hasID bool) error {
	if condition.Type < models.NUMERICSTATE {
		return fmt.Errorf("%d is not valid for subcondition", condition.Type)
	}

	id, err := hdb.get_condition_id(condition.Identity)
	if err != nil {
		return err
	}

	if hasID {
		id = condition.ID
	}

	if condition.Type == models.NUMERICSTATE {
		const query = "INSERT INTO conditions(id, identity, condition_type, sensor_id, above, below, parent_id) VALUES($1, $2, $3, $4, $5, $6, $7);"
		var below sql.NullInt32 = sql.NullInt32{}
		var above sql.NullInt32 = sql.NullInt32{}
		if condition.Below != nil {
			below = sql.NullInt32{Valid: condition.Below.Valid, Int32: int32(condition.Below.Value)}
		}
		if condition.Above != nil {
			above = sql.NullInt32{Valid: condition.Above.Valid, Int32: int32(condition.Above.Value)}
		}

		_, err = hdb.psqlDB.Exec(query, id, condition.Identity, condition.Type, condition.SensorID, above, below, parentID)
		return err
	} else if condition.Type == models.STATE {
		const query = "INSERT INTO conditions(id, identity, condition_type, sensor_id, comparison_state, parent_id) VALUES($1, $2, $3, $4, $5, $6);"
		_, err = hdb.psqlDB.Exec(query, id, condition.Identity, condition.Type, condition.SensorID, condition.ComparisonState, parentID)
		return err
	} else if condition.Type == models.TIME {
		var before sql.NullString = sql.NullString{
			Valid:  len(condition.Before) > 0,
			String: condition.Before,
		}
		var after sql.NullString = sql.NullString{
			Valid:  len(condition.After) > 0,
			String: condition.After,
		}
		const query = "INSERT INTO conditions(id, identity, condition_type, after, before, parent_id) VALUES($1, $2, $3, $4, $5, $6);"
		_, err = hdb.psqlDB.Exec(query, id, condition.Identity, condition.Type, after, before, parentID)
		return err
	}
	return fmt.Errorf("%d is not valid for subcondition", condition.Type)
}

// GetCondition gibt eine Bedingung anhand der Bedingungs-ID und Identität zurück.
func (hdb *HonuaDB) GetCondition(conditionID int, identity string) (*models.Condition, error) {
	const query = "SELECT id, identity, condition_type FROM conditions WHERE id=$1 AND identity=$2;"
	rows, err := hdb.psqlDB.Query(query, conditionID, identity)
	if err != nil {
		return nil, err
	}

	var result *models.Condition

	for rows.Next() {
		result, err = hdb.make_condition(rows)
		if err != nil {
			rows.Close()
			return nil, err
		}
	}

	rows.Close()

	return result, nil
}

// get_subconditions gibt alle Subbedingungen einer bestimmten Identität und übergeordneten Bedingungs-ID zurück.
func (hdb *HonuaDB) get_subconditions(identity string, parentID int) ([]*models.Condition, error) {
	const query = "SELECT * FROM conditions WHERE identity=$1 AND parent_id=$2;"
	rows, err := hdb.psqlDB.Query(query, identity, parentID)
	if err != nil {
		return nil, err
	}

	var result []*models.Condition = []*models.Condition{}

	for rows.Next() {
		sub, err := hdb.make_sub_condition(rows)
		if err != nil {
			rows.Close()
			return nil, err
		}
		result = append(result, sub)
	}

	rows.Close()

	return result, nil
}

// DeleteCondition löscht eine Bedingung anhand der Bedingungs-ID und Identität.
func (hdb *HonuaDB) DeleteCondition(conditionID int, identity string) error {
	const query = "DELETE FROM conditions WHERE id=$1 AND identity=$2;"
	_, err := hdb.psqlDB.Exec(query, conditionID, identity)
	return err
}

// get_condition_id gibt die ID einer Bedingung anhand der Identität zurück.
func (hdb *HonuaDB) get_condition_id(identity string) (int, error) {
	query := "SELECT CASE WHEN EXISTS ( SELECT * FROM conditions WHERE identity = $1) THEN true ELSE false END"

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

	query = "SELECT MAX(id) FROM conditions WHERE identity = $1;"

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

// make_condition erstellt ein Condition-Objekt basierend auf den Ergebnissen der Abfrage.
func (hdb *HonuaDB) make_condition(rows *sql.Rows) (*models.Condition, error) {
	var id int
	var identity string
	var cType models.ConditionType

	err := rows.Scan(&id, &identity, &cType)
	if err != nil {
		return nil, err
	}

	if cType >= models.NUMERICSTATE {
		return nil, fmt.Errorf("condition type %d not supported as parent condition", cType)
	}

	subs, err := hdb.get_subconditions(identity, id)
	if err != nil {
		return nil, err
	}

	return &models.Condition{ID: id, Identity: identity, Type: cType, SubConditions: subs}, nil
}

// make_sub_condition erstellt ein Subcondition-Objekt basierend auf den Ergebnissen der Abfrage.
func (hdb *HonuaDB) make_sub_condition(rows *sql.Rows) (*models.Condition, error) {
	var id int
	var identity string
	var cType models.ConditionType
	var sID sql.NullInt32
	var comparsionState sql.NullString
	var after sql.NullString
	var before sql.NullString
	var above sql.NullInt32
	var below sql.NullInt32
	var parentID int

	err := rows.Scan(&id, &identity, &cType, &sID, &comparsionState, &after, &before, &above, &below, &parentID)
	if err != nil {
		return nil, err
	}

	if cType < models.NUMERICSTATE {
		return nil, fmt.Errorf("condition type %d not supported as subcondition", cType)
	}

	if cType == models.NUMERICSTATE {
		if !sID.Valid || !(above.Valid || below.Valid) {
			return nil, errors.New("numeric_state condition is not valid")
		}

		return &models.Condition{
			ID:       id,
			Identity: identity,
			Type:     cType,
			SensorID: int(sID.Int32),
			Above:    &models.OptionalValue{Valid: above.Valid, Value: int(above.Int32)},
			Below:    &models.OptionalValue{Valid: below.Valid, Value: int(below.Int32)},
		}, nil

	} else if cType == models.STATE {
		if !sID.Valid || !comparsionState.Valid {
			return nil, errors.New("state condition is not valid")
		}

		return &models.Condition{
			ID:              id,
			Identity:        identity,
			Type:            cType,
			SensorID:        int(sID.Int32),
			ComparisonState: comparsionState.String,
		}, nil
	} else if cType == models.TIME {
		if !(after.Valid || before.Valid) {
			return nil, errors.New("time condition is not valid")
		}
		return &models.Condition{
			ID:       id,
			Identity: identity,
			Type:     cType,
			After:    after.String,
			Before:   before.String,
		}, nil
	}

	return nil, fmt.Errorf("%d is not valid for subcondition", cType)
}