package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Garden ...
type Garden struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"_name,omitempty"`
	LocationID primitive.ObjectID `json:"locationID" bson:"locationID"`
	// Location   Location           `json:"location" bson:"location"`
}
