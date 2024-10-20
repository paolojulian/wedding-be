package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client

	DatabaseName = "wedding_db"
	// Collection names
	InvitationsCollection = "invitations"
	UsersCollection       = "users"
)

func ConnectMongoDB() *mongo.Client {
	// If there is a client initialized, return it
	if client != nil {
		return client
	}

	fmt.Println("Connecting to MongoDB Atlas...")

	// If not, create one

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://ocaysanity:2imsx3QdCp5uCcDm@weddingcluster.btbm6.mongodb.net/?retryWrites=true&w=majority&appName=WeddingCluster").SetServerAPIOptions(serverAPI)

	var err error
	client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}

	// Ping the database to confirm the connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Atlas!")
	return client
}

// DisconnectMongoDB safely disconnects from the MongoDB Atlas cluster
func DisconnectMongoDB() {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	}
}
