package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Location ....
type Location struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name"`
	Zipcode int                `json:"zipcode" bson:"zipcode"`
	Garden  Garden             `json:"garden" bson:"garden"`
}
