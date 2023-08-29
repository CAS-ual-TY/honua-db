package models

type Dashboard struct {
	ID      string    `bson:"_id, omitempty"`
	Widgets []*Widget `bson:"widgets, omitempty"`
}

type WidgetType int32

const (
	ENTITY WidgetType = iota
	DEVICE
	WEATHER
	GROUP
	DEFAULT
	HEATMODE
)

type Widget struct {
	WidgetType        WidgetType `bson:"type" json:"type"`
	Title             string     `bson:"title,omitempty" json:"title,omitempty"`
	Unit              string     `bson:"unit,omitempty" json:"unit,omitempty"`
	EntityID          int32      `bson:"entity_id,omitempty" json:"entity_id,omitempty"`
	SecondaryEntityID int32      `bson:"secondary_entity_id,omitempty" json:"secondary_entity_id,omitempty"`
	SecondTitle       string     `bson:"title_2,omitempty" json:"title_2,omitempty"`
	ThirdEntityID     int32      `bson:"third_entity_id,omitempty" json:"third_entity_id,omitempty"`
	ThirdTitle        string     `bson:"title_3,omitempty" json:"title_3,omitempty"`
	FourthEntityID    int32      `bson:"fourth_entity_id,omitempty" json:"fourth_entity_id,omitempty"`
	FourthTitle       string     `bson:"title_4,omitempty" json:"title_4,omitempty"`
	FifthEntityID     int32      `bson:"fifth_entity_id,omitempty" json:"fifth_entity_id,omitempty"`
	FifthTile         string     `bson:"title_5,omitempty" json:"title_5,omitempty"`
	Subtitle          string     `bson:"subtitle,omitempty" json:"subtitle,omitempty"`
	SwitchRules       bool       `bson:"switch_rules,omitempty" json:"switch_rules,omitempty"`
	Expandable        bool       `bson:"expandable,omitempty" json:"expandable,omitempty"`
	Cards             []*Widget  `bson:"cards,omitempty" json:"cards,omitempty"`
}