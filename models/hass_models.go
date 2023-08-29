package models

type Service struct {
	Domain      string `bson:"domain"`
	Name        string `bson:"name"`
	Description string `bson:"description,omitempty"`
}
