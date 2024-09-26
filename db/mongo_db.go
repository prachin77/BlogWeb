package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/prachin77/server/models"
	"go.mongodb.org/mongo-driver/bson"
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

func AddBlog(blog *models.Blog) (*models.Blog, error) {

	_, err := collection.InsertOne(context.TODO(), bson.M{
		"blog_title":         blog.BlogTitle,
		"blog_tag":           blog.Tag,
		"blog_content":       blog.BlogContent,
		"author_name":        blog.AuthorName,
		"blog_creation_date": blog.BlogCreationDate,
	})

	if err != nil {
		return nil, err
	}

	fmt.Println("Blog added successfully!")
	return blog, nil 
}
