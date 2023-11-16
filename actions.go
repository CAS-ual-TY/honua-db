package honuadb

import (
	"database/sql"
	"errors"

	"github.com/JonasBordewick/honua-db/models"
)


// AddAction fügt eine Aktion zur Datenbank hinzu, basierend auf den übergebenen Parametern.
// Wenn hasID true ist, wird die mitgegebene ID verwendet, andernfalls wird eine neue ID generiert.
// isThenAction gibt an, ob es sich um eine "then"-Aktion handelt, und ruleID ist die ID der zugehörigen Regel.
func (hdb *HonuaDB) AddAction(action *models.Action, hasID, isThenAction bool, ruleID int) error {
	id, err := hdb.get_action_id(action.Identity)
	if err != nil {
		return err
	}

	if hasID {
		id = action.ID
	}

	if action.Type == models.SERVICE {
		const query = "INSERT INTO actions(id, identity, is_then_Action, action_type, service_id, rule_id) VALUES($1, $2, $3, $4, $5, $6);"
		_, err = hdb.psqlDB.Exec(query, id, action.Identity, isThenAction, models.SERVICE, action.ServiceID, ruleID)
		if err != nil {
			return err
		}
	} else {
		delayID, err := hdb.AddDelay(action.Delay, hasID)
		if err != nil {
			return err
		}
		const query = "INSERT INTO actions(id, identity, is_then_Action, action_type, delay_id, rule_id) VALUES($1, $2, $3, $4, $5, $6);"
		_, err = hdb.psqlDB.Exec(query, id, action.Identity, isThenAction, models.DELAY, delayID, ruleID)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetActions gibt alle Aktionen für eine gegebene Identität und Regel-ID zurück.
// Es werden zwei Slices von Aktionen zurückgegeben: "thenActions" und "elseActions".
func (hdb *HonuaDB) GetActions(identity string, ruleID int) ([]*models.Action, []*models.Action, error) {
	const query = "SELECT * FROM actions WHERE identity=$1 AND rule_id=$2"
	rows, err := hdb.psqlDB.Query(query, identity, ruleID)
	if err != nil {
		return nil, nil, err
	}

	var thenActions []*models.Action = []*models.Action{}
	var elseActions []*models.Action = []*models.Action{}

	for rows.Next() {
		action, err := hdb.make_action(rows)
		if err != nil {
			rows.Close()
			return nil, nil, err
		}
		if action.IsThenAction {
			thenActions = append(thenActions, action)
		} else {
			elseActions = append(elseActions, action)
		}
	}

	rows.Close()

	return thenActions, elseActions, nil
}

// make_action erstellt ein Action-Objekt aus den Daten einer SQL-Abfrage.
func (hdb *HonuaDB) make_action(rows *sql.Rows) (*models.Action, error) {
	var id int
	var identity string
	var isThenAction bool
	var actionType models.ActionType
	var serviceID sql.NullInt32
	var delayID sql.NullInt32
	var ruleID int

	err := rows.Scan(&id, &identity, &isThenAction, &actionType, &serviceID, &delayID, &ruleID)
	if err != nil {
		return nil, err
	}

	if actionType == models.SERVICE {
		return &models.Action{ID: id, Identity: identity, Type: actionType, IsThenAction: isThenAction, ServiceID: serviceID.Int32}, nil
	} else {
		delay, err := hdb.GetDelay(int(delayID.Int32), identity)
		if err != nil {
			return nil, err
		}
		return &models.Action{ID: id, Identity: identity, Type: actionType, IsThenAction: isThenAction, Delay: delay}, nil
	}
}

// get_action_id gibt die nächste verfügbare ID für eine Aktion zurück.
func (hdb *HonuaDB) get_action_id(identity string) (int, error) {
	query := "SELECT CASE WHEN EXISTS ( SELECT * FROM actions WHERE identity = $1) THEN true ELSE false END"

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

	query = "SELECT MAX(id) FROM actions WHERE identity = $1;"

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
		return -1, errors.New("something went wrong during getting id of action")
	}

	id = id + 1

	return id, nil
}
