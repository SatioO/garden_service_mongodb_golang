package response

import (
	"github.com/satioO/gardens/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Garden struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"_name,omitempty"`
	Location models.Location    `json:"location,omitempty" bson:"location,omitempty"`
}
