package models

type DetailedDashboard struct {
	ID      string          `bson:"_id"`
	Widgets []*DetailWidget `bson:"widgets"`
}

const (
	GroupWidgetType int = iota
	EntityWidgetType
	DoubleEntityWidgetType
)

type DetailWidget struct {
	WidgetType               int               `bson:"widgetType"`
	Title                    string            `bson:"title"`
	PrimaryAlternativeStates map[string]string `bson:"PrimaryAlternativeStates"`
	PrimaryEntityID          int32             `bson:"primaryEntityID"`
	SecondaryEntityID        int32             `bson:"secondaryEntityID"`
	PrimaryUnit              string            `bson:"primaryUnit"`
	SecondaryUnit            string            `bson:"secondaryUnit"`
	Childs                   []*DetailWidget   `bson:"childs"`
}
