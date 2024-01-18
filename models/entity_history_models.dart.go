package models

type EntityHistory struct {
	ID       string            `bson:"_id" json:"id"`
	Elements []*HistoryElement `bson:"elements" json:"elements"`
}

type HistoryElement struct {
	X int64   `bson:"x" json:"x"`
	Y float64 `bson:"y" json:"y"`
}
