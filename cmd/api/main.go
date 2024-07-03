package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lucasgarciaf/df-backend-go/config"
	"github.com/lucasgarciaf/df-backend-go/internal/middleware"
	"github.com/lucasgarciaf/df-backend-go/internal/router"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables
	log.Printf("Load environment variables")
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found: %v", err)
	}

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

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Successfully connected and pinged MongoDB!")

	// // Ensure MongoDB disconnection on shutdown
	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		log.Fatalf("Failed to disconnect MongoDB: %v", err)
	// 	}
	// }()

	fmt.Println("Starting the application...")

	// Create a new gin router
	r := gin.Default()

	// Setup the router
	router.SetupRouter(r, db)

	// Start the server in a goroutine
	server := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	fmt.Println("Server running on :8081")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	fmt.Println("Server exiting")
}
