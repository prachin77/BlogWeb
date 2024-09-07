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

func SearchBlogWithId(blogId string) (error, bool) {
	// Create a filter to search for the blog by its ID
	filter := bson.M{"blogid": blogId}

	// Perform the query
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err == mongo.ErrNoDocuments {
		// No document found with the given ID
		return nil, false
	} else if err != nil {
		// Some other error occurred
		return err, false
	}

	// Blog found
	return nil, true
}


func AddBlog(blogData map[string]interface{}, blog *models.Blog) {
	// Check if the blog already exists
	err, isBlogThere := SearchBlogWithId(blog.BlogId)
	if err != nil {
		fmt.Printf("Error checking if blog exists: %v\n", err)
		return
	}
	if isBlogThere {
		fmt.Println("Blog with the same ID already present")
		// now create new blog id & add blog 
		return
	}

	// Prepare the blog document for insertion
	blogDocument := bson.M{
		"blogid":                 blog.BlogId,
		"blogtitle":              blog.BlogTitle,
		"blogcontent":            blog.BlogContent,
		"tag":                    blog.Tags,
		"authorname":             blog.AuthorName,
		"image_data_in_bytes":    blogData["ImageDataInBytes"], // Use the base64-encoded string if needed
		"blog_image_data_length": blogData["BlogImageDataLength"],
		"userid":                 blog.UserId,
		"blog_creation_date":     blog.BlogCreationDate,
	}

	// Insert the blog document into the collection
	_, err = collection.InsertOne(context.TODO(), blogDocument)
	if err != nil {
		fmt.Printf("Error inserting blog into MongoDB: %v\n", err)
		return
	}

	fmt.Println("Blog successfully inserted into MongoDB.")
}


func GetAllBlogs() {

}
