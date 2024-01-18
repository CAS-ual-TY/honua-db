package models

// Dashboard repräsentiert ein Dashboard-Objekt mit einer eindeutigen ID und einer Liste von Widgets.
type Dashboard struct {
	ID      string    `bson:"_id, omitempty"`     // Eindeutige ID des Dashboards
	Widgets []*Widget `bson:"widgets, omitempty"` // Liste von Widgets auf dem Dashboard
}

// WidgetType definiert die verschiedenen Arten von Widgets, die auf einem Dashboard angezeigt werden können.
type WidgetType int32

const (
	ENTITY   WidgetType = iota // Entitätstyp
	DEVICE                     // Gerätetyp
	WEATHER                    // Wettertyp
	GROUP                      // Gruppentyp
	DEFAULT                    // Standardtyp
	HEATMODE                   // Heizmodustyp
)

// Widget repräsentiert ein einzelnes Widget auf einem Dashboard mit verschiedenen Eigenschaften.
type Widget struct {
	WidgetType        WidgetType `bson:"type" json:"type"`                       // Typ des Widgets
	Title             string     `bson:"title,omitempty" json:"title,omitempty"` // Titel des Widgets
	Icon              string     `bson:"icon,omitempty" json:"icon,omitempty"`
	Color             string     `bson:"color,omitempty" json:"color,omitempty"`
	Unit              string     `bson:"unit,omitempty" json:"unit,omitempty"`                               // Einheit des Widgets
	EntityID          int32      `bson:"entity_id,omitempty" json:"entity_id,omitempty"`                     // ID der Entität
	SecondaryEntityID int32      `bson:"secondary_entity_id,omitempty" json:"secondary_entity_id,omitempty"` // Sekundäre ID der Entität
	SecondTitle       string     `bson:"title_2,omitempty" json:"title_2,omitempty"`                         // Titel für die zweite Entität
	ThirdEntityID     int32      `bson:"third_entity_id,omitempty" json:"third_entity_id,omitempty"`         // Dritte ID der Entität
	ThirdTitle        string     `bson:"title_3,omitempty" json:"title_3,omitempty"`                         // Titel für die dritte Entität
	FourthEntityID    int32      `bson:"fourth_entity_id,omitempty" json:"fourth_entity_id,omitempty"`       // Vierte ID der Entität
	FourthTitle       string     `bson:"title_4,omitempty" json:"title_4,omitempty"`                         // Titel für die vierte Entität
	FifthEntityID     int32      `bson:"fifth_entity_id,omitempty" json:"fifth_entity_id,omitempty"`         // Fünfte ID der Entität
	FifthTile         string     `bson:"title_5,omitempty" json:"title_5,omitempty"`                         // Titel für die fünfte Entität
	Subtitle          string     `bson:"subtitle,omitempty" json:"subtitle,omitempty"`                       // Untertitel des Widgets
	SwitchRules       bool       `bson:"switch_rules,omitempty" json:"switch_rules,omitempty"`               // Schalter für Regeln
	Expandable        bool       `bson:"expandable,omitempty" json:"expandable,omitempty"`                   // Anzeige erweiterbar
	Cards             []*Widget  `bson:"cards,omitempty" json:"cards,omitempty"`                             // Liste von untergeordneten Widgets (Karten)
}
