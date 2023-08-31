package models

type Service struct {
	Domain      string `bson:"domain"`
	Name        string `bson:"name"`
	Description string `bson:"description,omitempty"`
}

type Services struct {
	Id       string     `bson:"_id, omitempty"`
	Services []*Service `bson:"services, omitempty"`
}
