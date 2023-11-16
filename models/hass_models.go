package models

// Service repräsentiert einen Service mit einer Domäne, einem Namen und optional einer Beschreibung.
type Service struct {
	Domain      string `bson:"domain"`                  // Domäne des Dienstes
	Name        string `bson:"name"`                    // Name des Dienstes
	Description string `bson:"description,omitempty"`  // Optionale Beschreibung des Dienstes
}

// Services repräsentiert eine Sammlung von Services mit einer eindeutigen ID und einer Liste von Service-Objekten.
type Services struct {
	Id       string     `bson:"_id, omitempty"`  // Eindeutige ID der Services
	Services []*Service `bson:"services, omitempty"`  // Liste von Service-Objekten
}
