package models

type EntityInformationElement struct {
	EntityID     string            `bson:"entity_id" json:"entity_id"`
	FriendlyName string            `bson:"friendly_name" json:"friendly_name"`
	State        string            `bson:"state" json:"state"`
	Attributes   map[string]string `bson:"attributes" json:"attributes"`
}

type EntityInformation struct {
	ID       string                      `bson:"_id" json:"id"`
	Elements []*EntityInformationElement `bson:"elements" json:"elements"`
}
