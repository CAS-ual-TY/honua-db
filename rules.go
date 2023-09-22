package honuadb

import (
	"database/sql"
	"errors"
	"log"

	"github.com/JonasBordewick/honua-db/models"
)

func (hdb *HonuaDB) AddRule(rule *models.Rule, hasID bool) (int, error) {
	// First add Condition
	cID, err := hdb.AddCondition(rule.Condition, hasID)
	if err != nil {
		return -1, err
	}

	var periodicTrigger sql.NullInt32 = sql.NullInt32{
		Valid: !rule.EventBasedEvaluation,
		Int32: int32(rule.PeriodicTrigger),
	}

	const query = "INSERT INTO rules(id, identity, target_id, enabled, event_based_evaluation, periodic_trigger_type, condition_id) VALUES($1, $2, $3, $4, $5, $6, $7);"

	rID, err := hdb.get_rule_id(rule.Identity)
	if err != nil {
		hdb.DeleteCondition(cID, rule.Identity)
		return -1, err
	}

	if hasID {
		rID = rule.ID
	}

	_, err = hdb.psqlDB.Exec(query, rID, rule.Identity, rule.TargetID, rule.Enabled, rule.EventBasedEvaluation, periodicTrigger, cID)
	if err != nil {
		hdb.DeleteCondition(cID, rule.Identity)
		return -1, err
	}

	for _, action := range rule.ThenActions {
		err = hdb.AddAction(action, hasID, true, rID)
		if err != nil {
			return -1, err
		}
	}

	for _, action := range rule.ElseActions {
		err = hdb.AddAction(action, hasID, false, rID)
		if err != nil {
			return -1, err
		}
	}

	return rID, nil
}

func (hdb *HonuaDB) GetRules(identity string) ([]*models.Rule, error) {
	const query = "SELECT * FROM rules WHERE identity=$1;"
	rows, err := hdb.psqlDB.Query(query, identity)
	if err != nil {
		return nil, err
	}

	var rules []*models.Rule = []*models.Rule{}

	for rows.Next() {
		rule, err := hdb.make_rule(rows)
		if err != nil {
			rows.Close()
			return nil, err
		}

		rules = append(rules, rule)
	}

	rows.Close()

	return rules, nil
}

func (hdb *HonuaDB) GetRule(ruleID int, identity string) (*models.Rule, error) {
	const query = "SELECT * from rules WHERE id=$1 AND identity=$2;"
	rows, err := hdb.psqlDB.Query(query, ruleID, identity)
	if err != nil {
		return nil, err
	}

	var rule *models.Rule

	for rows.Next() {
		rule, err = hdb.make_rule(rows)
		if err != nil {
			return nil, err
		}
	}

	rows.Close()

	if rule == nil {
		return nil, errors.New("something went wrong")
	}

	return rule, nil
}

func (hdb *HonuaDB) GetRuleOfTarget(targetID int, identity string) (*models.Rule, error) {
	const query = "SELECT * from rules WHERE target_id=$1 AND identity=$2;"
	rows, err := hdb.psqlDB.Query(query, targetID, identity)
	if err != nil {
		return nil, err
	}

	var rule *models.Rule

	for rows.Next() {
		rule, err = hdb.make_rule(rows)
		if err != nil {
			rows.Close()
			return nil, err
		}
	}

	rows.Close()

	if rule == nil {
		return nil, errors.New("something went wrong")
	}

	return rule, nil
} 

func (hdb *HonuaDB) DeleteRule(delayID int, identity string) error {
	const query = "DELETE FROM delays WHERE id=$1 AND identity=$2;"
	_, err := hdb.psqlDB.Exec(query, delayID, identity)
	return err
}

func (hdb *HonuaDB) GetStateOfRule(targetID int, identity string) (bool, error) {
	const query = "SELECT enabled FROM rules WHERE target_id=$1 AND identity=$2;"
	rows, err := hdb.psqlDB.Query(query, targetID, identity)
	if err != nil {
		return false, err
	}

	var enabled bool = false

	for rows.Next() {
		err := rows.Scan(&enabled)
		if err != nil {
			rows.Close()
			return false, err
		}
	}

	rows.Close()

	return enabled, nil
}

func (hdb *HonuaDB) HasRule(targetID int32, identity string) bool {
	const query = "SELECT CASE WHEN EXISTS ( SELECT * FROM rules WHERE target_id=$1 AND identity=$2) THEN true ELSE false END;"
	rows, err := hdb.psqlDB.Query(query, identity)
	if err != nil {
		return false
	}

	var state bool = false

	for rows.Next() {
		err = rows.Scan(&state)
		if err != nil {
			rows.Close()
			return false
		}
	}

	rows.Close()

	return state
}

func (hdb *HonuaDB) ToggleRule(ruleID int32, state bool, identity string) error {
	const query = "UPDATE rules SET enabled=$1 WHERE id=$2 AND identity=$3;"
	_, err := hdb.psqlDB.Exec(query, state, ruleID, identity)
	return err
}

func (hdb *HonuaDB) get_rule_id(identity string) (int, error) {
	query := "SELECT CASE WHEN EXISTS ( SELECT * FROM rules WHERE identity = $1) THEN true ELSE false END"

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

	query = "SELECT MAX(id) FROM rules WHERE identity = $1;"

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

func (hdb *HonuaDB) make_rule(rows *sql.Rows) (*models.Rule, error) {
	var id int
	var identity string
	var targetID int32
	var enabled bool
	var eventBasedEvaluation bool
	var triggerType sql.NullInt32
	var conditionID int

	err := rows.Scan(&id, &identity, &targetID, &enabled, &eventBasedEvaluation, &triggerType, &conditionID)
	if err != nil {
		return nil, err
	}

	condition, err := hdb.GetCondition(conditionID, identity)
	if err != nil {
		return nil, err
	}

	thenActions, elseActions, err := hdb.GetActions(identity, id)
	if err != nil {
		return nil, err
	}

	log.Printf("PeriodicTriggerType %d\n", triggerType.Int32)

	if !eventBasedEvaluation && triggerType.Valid {
		return &models.Rule{
			ID:                   id,
			Identity:             identity,
			TargetID:             int(targetID),
			Enabled:              enabled,
			EventBasedEvaluation: false,
			PeriodicTrigger:      models.PeriodicTriggerType(triggerType.Int32),
			Condition:            condition,
			ThenActions:          thenActions,
			ElseActions:          elseActions,
		}, nil
	} else {
		return &models.Rule{
			ID:                   id,
			Identity:             identity,
			TargetID:             int(targetID),
			Enabled:              enabled,
			EventBasedEvaluation: true,
			Condition:            condition,
			ThenActions:          thenActions,
			ElseActions:          elseActions,
		}, nil
	}
}
