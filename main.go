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

func init() {
	ctx := context.TODO()

	mongoconn := options.Client().ApplyURI("mongodb://localhost:27017/")

	mongoclient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		log.Fatal("err while connecting with MongoDB Atlas", err)
	}

	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Error while trying to ping MongoDB Atlas", err)
	}

	fmt.Println("MongoDB Atlas connection established")

	// Update the database and collection names as needed.
	userc := mongoclient.Database("user_data").Collection("users")
	us := services.NewUserService(userc, ctx)
	uc := controllers.New(us)
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:8080"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	}))

	defer mongoclient.Disconnect(ctx)

	basepath := server.Group("/v1")
	uc.RegisterUserRoutes(basepath)

	log.Fatal(server.Run(":9090"))
}
func main() {

}
