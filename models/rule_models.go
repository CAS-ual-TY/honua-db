package models

// PeriodicTriggerType definiert verschiedene Perioden für periodische Auslöser.
type PeriodicTriggerType int

const (
	OneMin       PeriodicTriggerType = iota
	TwoMin                            // ...
	FiveMin
	TenMin
	FifteenMin
	TwentyMin
	TwentyFiveMin
	ThirtyMin
	FortyFiveMin
	OneH
	TwoH
	SixH
)

// Rule repräsentiert eine Regel mit verschiedenen Eigenschaften.
type Rule struct {
	ID                   int                  // Eindeutige ID der Regel
	Identity             string               // Identität der Regel
	TargetID             int                  // Ziel-ID der Regel
	Enabled              bool                 // Gibt an, ob die Regel aktiviert ist
	EventBasedEvaluation bool                 // Gibt an, ob die Regel auf Ereignissen basiert
	PeriodicTrigger      PeriodicTriggerType  // Art des periodischen Auslösers
	Condition            *Condition           // Bedingung der Regel
	ThenActions          []*Action            // Aktionen, die bei Erfüllung der Bedingung ausgeführt werden
	ElseActions          []*Action            // Aktionen, die bei Nichterfüllung der Bedingung ausgeführt werden
}

// ConditionType definiert verschiedene Arten von Bedingungen.
type ConditionType int

const (
	OR          ConditionType = iota  // Oder-Bedingung
	AND                               // Und-Bedingung
	NOR                               // Nicht Oder-Bedingung
	NAND                              // Nicht Und-Bedingung
	NUMERICSTATE                      // Numerischer Zustand
	STATE                             // Zustand
	TIME                              // Zeit
)

// OptionalValue repräsentiert einen optionalen numerischen Wert mit einem Gültigkeitsstatus.
type OptionalValue struct {
	Valid bool // Gibt an, ob der Wert gültig ist
	Value int  // Numerischer Wert (falls gültig)
}

// Condition repräsentiert eine Bedingung mit verschiedenen Eigenschaften.
type Condition struct {
	ID              int             // Eindeutige ID der Bedingung
	Identity        string          // Identität der Bedingung
	Type            ConditionType   // Typ der Bedingung
	SensorID        int             // ID des Sensors in der Bedingung (falls anwendbar)
	ComparisonState string          // Zustand für Vergleichsbedingungen (falls anwendbar)
	After           string          // Zeitpunkt nach dem die Bedingung gültig ist (falls anwendbar)
	Before          string          // Zeitpunkt vor dem die Bedingung gültig ist (falls anwendbar)
	Above           *OptionalValue  // Optionaler Wert für Über-Bedingungen (falls anwendbar)
	Below           *OptionalValue  // Optionaler Wert für Unter-Bedingungen (falls anwendbar)
	SubConditions   []*Condition    // Unterbedingungen
}

// ActionType definiert verschiedene Arten von Aktionen.
type ActionType int

const (
	SERVICE ActionType = iota  // Service-Aktion
	DELAY                     // Verzögerungsaktion
)

// Action repräsentiert eine Aktion mit verschiedenen Eigenschaften.
type Action struct {
	ID           int       // Eindeutige ID der Aktion
	Identity     string    // Identität der Aktion
	Type         ActionType // Typ der Aktion
	IsThenAction bool      // Gibt an, ob die Aktion eine "then"-Aktion ist
	ServiceID    int32      // ID des Dienstes in der Aktion (falls anwendbar)
	Delay        *Delay     // Verzögerung für Verzögerungsaktionen (falls anwendbar)
}

// Delay repräsentiert eine Verzögerung mit Stunden, Minuten und Sekunden.
type Delay struct {
	ID       int    // Eindeutige ID der Verzögerung
	Identity string // Identität der Verzögerung
	Hours    int32  // Stunden der Verzögerung
	Minutes  int32  // Minuten der Verzögerung
	Seconds  int32  // Sekunden der Verzögerung
}
