package main

import (
	"context"
	"fmt"
	"log"

	"example.com/noman-ali/controllers"
	"example.com/noman-ali/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// init function is called before the main function and is used for initialization tasks.
func init() {
	// Create a context for MongoDB operations.
	ctx := context.TODO()

	// Configure MongoDB connection options.
	mongoconn := options.Client().ApplyURI("mongodb+srv://nextjs13:nextjs13@cluster0.v4qjg7w.mongodb.net/?retryWrites=true&w=majority")

	// Connect to MongoDB using the specified options.
	mongoclient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		log.Fatal("err while connecting with MongoDB Atlas", err)
	}

	// Ping the MongoDB server to check if it's available.
	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Error while trying to ping MongoDB Atlas", err)
	}

	fmt.Println("MongoDB Atlas connection established")

	// Create a MongoDB collection for user data.
	userc := mongoclient.Database("user_data").Collection("users")

	// Create a user service with the MongoDB collection.
	us := services.NewUserService(userc, ctx)

	// Create a user controller with the user service.
	uc := controllers.New(us)

	// Create a Gin server instance with default middleware.
	server := gin.Default()

	// Configure CORS (Cross-Origin Resource Sharing) middleware to allow requests from localhost:8080.
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:8080"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	}))

	defer mongoclient.Disconnect(ctx)

	// Create a basepath group for API versioning.
	basepath := server.Group("/v1")

	// Register user-related routes using the user controller.
	uc.RegisterUserRoutes(basepath)

	// Start the Gin server on port 9090 and handle incoming requests.
	log.Fatal(server.Run(":9090"))
}

// The main function is the entry point of the application, but it is empty in this case.
func main() {

}
