package models

type PeriodicTriggerType int

const (
	OneMin PeriodicTriggerType = iota
	TwoMin
	FiveMin
	TenMin
	FifteenMin
	TwentyMin
	TwentyFiveMin
	FortyFiveMin
	OneH
	TwoH
	SixH
)

type Rule struct {
	ID                   int
	Identity             string
	TargetID             int
	Enabled              bool
	EventBasedEvaluation bool
	PeriodicTrigger      PeriodicTriggerType
	Condition            *Condition
	ThenActions          []*Action
	ElseActions          []*Action
}

type ConditionType int

const (
	OR ConditionType = iota
	AND
	NOR
	NAND
	NUMERICSTATE
	STATE
	TIME
)

type OptionalValue struct {
	Valid bool
	Value int
}

type Condition struct {
	ID              int
	Identity        string
	Type            ConditionType
	SensorID        int
	ComparisonState string
	After           string
	Before          string
	Above           *OptionalValue
	Below           *OptionalValue
	SubConditions   []*Condition
}

type ActionType int

const (
	SERVICE ActionType = iota
	DELAY
)

type Action struct {
	ID           int
	Identity     string
	Type         ActionType
	IsThenAction bool
	ServiceID    int32
	Delay        *Delay
}

type Delay struct {
	ID       int
	Identity string
	Hours    int32
	Minutes  int32
	Seconds  int32
}
