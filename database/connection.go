package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Extract MongoDB connection details from environment variables
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")
	clusterURL := os.Getenv("MONGODB_CLUSTER_URL")
	dbName := os.Getenv("MONGODB_DATABASE")

	// Check if all required environment variables are set
	if username == "" || password == "" || clusterURL == "" || dbName == "" {
		return nil, fmt.Errorf("missing required MongoDB environment variables")
	}

	// Construct the MongoDB connection URI
	uri := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority",
		username, password, clusterURL, dbName)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the database to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	return client.Database(dbName), nil
}
