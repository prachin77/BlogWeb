package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// this is a pointer(reference) to collection in mongo db
var collection *mongo.Collection

func init() {
	if err := godotenv.Load("P:/BlogWeb/.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ConnectionString := os.Getenv("ConnectionString")
	DbName := os.Getenv("DbName")
	CollectionName := os.Getenv("CollectionName")

	clientOpt := options.Client().ApplyURI(ConnectionString)

	// connect to mongodb
	client, err := mongo.Connect(context.TODO(), clientOpt)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connection to mongo db successfull ✌️✌️")

	collection = client.Database(DbName).Collection(CollectionName)

	// collection instance
	fmt.Println("collection instance is ready")
}

