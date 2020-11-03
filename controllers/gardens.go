package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/satioO/gardens/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GardenController ...
type GardenController struct {
	db *mongo.Database
}

// NewGardenController ...
func NewGardenController(db *mongo.Database) *GardenController {
	return &GardenController{db}
}

// ListGardens ...
func (g *GardenController) ListGardens(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := []bson.M{
		bson.M{
			"$lookup": bson.M{
				"from":         "location",
				"as":           "location",
				"localField":   "locationID",
				"foreignField": "_id",
			},
		},
		bson.M{"$unwind": "$location"},
	}

	cursor, err := g.db.Collection("garden").Aggregate(ctx, pipeline)
	gardens := []*models.Garden{}
	err = cursor.All(ctx, &gardens)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went Wrong ..!"))
		return
	}

	response, _ := json.Marshal(gardens)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// SaveGarden ...
func (g *GardenController) SaveGarden(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var garden models.Garden
	err := json.NewDecoder(r.Body).Decode(&garden)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cursor, err := g.db.Collection("garden").InsertOne(ctx, &garden)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id := cursor.InsertedID.(primitive.ObjectID).String()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Inserted: " + id))
}

// FindByLocation ...
func (g *GardenController) FindByLocation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	locationID, _ := primitive.ObjectIDFromHex(mux.Vars(r)["locationID"])

	var gardens *models.Garden
	cursor, err := g.db.Collection("garden").Find(ctx, bson.M{"locationID": locationID})

	err = cursor.All(ctx, &gardens)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(&gardens)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(response)
}
