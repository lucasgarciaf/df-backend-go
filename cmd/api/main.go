package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lucasgarciaf/df-backend-go/config"
	"github.com/lucasgarciaf/df-backend-go/internal/middleware"
	"github.com/lucasgarciaf/df-backend-go/internal/router"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Initialize Keycloak
	if err := middleware.InitKeycloak(); err != nil {
		log.Fatalf("Failed to initialize Keycloak: %v", err)
	}

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(config.MongoDBURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	db := client.Database(config.DatabaseName)

	fmt.Println("Starting the application...")
	r := gin.Default()

	// Setup the router
	router.SetupRouter(r, db)

	fmt.Println("Starting server on :8081")
	log.Fatal(r.Run(":8081")) // Ensure the server runs on port 8081
}
