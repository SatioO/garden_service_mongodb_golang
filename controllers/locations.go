package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/satioO/gardens/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// LocationController ...
type LocationController struct {
	db *mongo.Database
}

// NewLocationController ...
func NewLocationController(db *mongo.Database) *LocationController {
	return &LocationController{db}
}

// ListLocation ...
func (l *LocationController) ListLocation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := l.db.Collection("location").Find(ctx, bson.M{})
	locations := []models.Location{}
	err = cursor.All(ctx, &locations)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went Wrong ..!"))
		return
	}

	response, _ := json.Marshal(locations)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// SaveLocation ...
func (l *LocationController) SaveLocation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var location models.Location
	err := json.NewDecoder(r.Body).Decode(&location)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cursor, err := l.db.Collection("location").InsertOne(ctx, &location)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id := cursor.InsertedID.(primitive.ObjectID).String()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Inserted: " + id))
}

// FindByZipcode ...
func (l *LocationController) FindByZipcode(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	zipcode, _ := strconv.Atoi(mux.Vars(r)["zipcode"])

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"zipcode": zipcode,
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "garden",
				"as":           "garden",
				"localField":   "_id",
				"foreignField": "locationID",
			},
		},
		bson.M{"$unwind": "$garden"},
	}

	cursor, err := l.db.Collection("location").Aggregate(ctx, pipeline)
	locations := []models.Location{}
	err = cursor.All(ctx, &locations)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went Wrong ..!"))
		return
	}

	response, _ := json.Marshal(locations)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
