package models

import "time"

type Identity struct {
	ID   string
	Name string
}

type VictronSensorType int32

// ALL Victron Sensor Types
const (
	NONE VictronSensorType = iota
	ACLOADS
	TOTALPV
	GRID
	SOC
	BATTERYVALUE
	BATTERYSTATE
)

type Entity struct {
	ID              int32
	Identity        string
	EntityID        string
	Name            string
	IsDevice        bool
	AllowRules      bool
	HasAttribute    bool
	Attribute       string
	IsVictronSensor bool
	SensorType      VictronSensorType
	HasNumericState bool
}

type State struct {
	ID         int32
	EntityID   int32
	Identity   string
	State      string
	RecordTime *time.Time
}

type HonuaService struct {
	ID       int32
	Identity string
	Domain   string
	Name     string
}
