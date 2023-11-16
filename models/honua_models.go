package models

import "time"

// Identity repräsentiert eine Identität mit einer eindeutigen ID und einem Namen.
type Identity struct {
	ID   string // Eindeutige ID der Identität
	Name string // Name der Identität
}

// VictronSensorType definiert verschiedene Arten von Victron-Sensoren.
type VictronSensorType int32

// ALL Victron Sensor Types
const (
	NONE VictronSensorType = iota     // Kein Victron-Sensor
	ACLOADS                          // AC-Lasten-Sensor
	TOTALPV                          // Gesamter PV-Sensor
	GRID                             // Netz-Sensor
	SOC                              // SOC-Sensor (State of Charge)
	BATTERYVALUE                     // Batteriewert-Sensor
	BATTERYSTATE                     // Batteriezustand-Sensor
)

// Entity repräsentiert eine Entität mit verschiedenen Eigenschaften.
type Entity struct {
	ID              int32             // Eindeutige ID der Entität
	Identity        string            // Identität der Entität
	EntityID        string            // ID der Entität
	Name            string            // Name der Entität
	IsDevice        bool              // Gibt an, ob die Entität ein Gerät ist
	AllowRules      bool              // Gibt an, ob Regeln für die Entität erlaubt sind
	HasAttribute    bool              // Gibt an, ob die Entität ein Attribut hat
	Attribute       string            // Attribut der Entität
	IsVictronSensor bool              // Gibt an, ob die Entität ein Victron-Sensor ist
	SensorType      VictronSensorType // Typ des Victron-Sensors
	HasNumericState bool              // Gibt an, ob die Entität einen numerischen Zustand hat
}

// State repräsentiert den Zustand einer Entität zu einem bestimmten Zeitpunkt.
type State struct {
	ID         int32         // Eindeutige ID des Zustands
	EntityID   int32         // ID der zugehörigen Entität
	Identity   string        // Identität der Entität
	State      string        // Zustand der Entität
	RecordTime *time.Time    // Zeitpunkt des Zustands
}

// HonuaService repräsentiert einen Honua-Service mit einer eindeutigen ID, Identität, Domäne und Namen.
type HonuaService struct {
	ID       int32  // Eindeutige ID des Honua-Services
	Identity string // Identität des Honua-Services
	Domain   string // Domäne des Honua-Services
	Name     string // Name des Honua-Services
}