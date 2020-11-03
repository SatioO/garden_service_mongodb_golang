package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/satioO/gardens/controllers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	db := ConnectDB("gardens")
	gc := controllers.NewGardenController(db)
	lc := controllers.NewLocationController(db)

	r := mux.NewRouter()
	r.HandleFunc("/gardens", gc.ListGardens).Methods(http.MethodGet)
	r.HandleFunc("/gardens", gc.SaveGarden).Methods(http.MethodPost)
	r.HandleFunc("/gardens/locations/{locationID}", gc.FindByLocation).Methods(http.MethodGet)

	r.HandleFunc("/locations", lc.ListLocation).Methods(http.MethodGet)
	r.HandleFunc("/locations", lc.SaveLocation).Methods(http.MethodPost)
	r.HandleFunc("/locations/{zipcode}", lc.FindByZipcode).Methods(http.MethodGet)

	log.Println("Server started")
	log.Fatal(http.ListenAndServe(":3000", r))
}

// ConnectDB ...
func ConnectDB(dbName string) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)

	err = client.Connect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected!")
	}

	return client.Database(dbName)
}
