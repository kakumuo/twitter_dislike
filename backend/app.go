package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	// "github.com/rs/cors"

	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DislikeRecord struct {
	OwnerId   string `bson:"ownerId"`
	UpdatedAt int    `bson:"updatedAt,$date"`
}

type TweetRecord struct {
	OwnerId   string                   `bson:"ownerId"`
	TweetId   string                   `bson:"tweetId"`
	DislikeBy map[string]DislikeRecord `bson:"dislikedBy"`
	CreatedAt int                      `bson:"createdAt,$date"`
	UpdatedAt int                      `bson:"updatedAt,$date"`
}

type Config struct {
	ConnectionString string `bson:"connectionString"`
	UserName         string `bson:"username"`
	Password         string `bson:"password"`
	HostName         string `bson:"hostname"`
	Port             int    `bson:"port"`
}

const HOST = "127.0.0.1"
const PORT = 8888
const FILE_PATH = "./data.json"

var coll *mongo.Collection

// TODO: move functionality to mongodb
func main() {
	mux := http.NewServeMux()

	var dislikeWriter DBDislikeWriter
	config := LoadConfig("./config.json")

	dislikeWriter = NewDBDislikeWriter(config.ConnectionString)
	defer func() {
		if err := dislikeWriter.Client.Disconnect(context.TODO()); err != nil {
			log.Fatal("Error while disconnecting from mongo instance:", err.Error())
		}
	}()

	mux.HandleFunc("GET /tweet/dislike", dislikeWriter.HandleGetDislikes)
	mux.HandleFunc("POST /tweet/dislike", dislikeWriter.HandlePostDislikes)

	// start server
	corsHandler := cors.AllowAll().Handler(mux)
	serverEndpoint := fmt.Sprintf("%s:%d", config.HostName, config.Port)
	log.Println("Starting server at:", serverEndpoint)
	log.Fatal(http.ListenAndServe(serverEndpoint, corsHandler))
}

func LoadConfig(filepath string) *Config {
	res := &Config{}

	fileData, err := os.ReadFile(filepath)

	if err != nil {
		log.Fatal("Error: ", err.Error())
	}

	err = json.Unmarshal(fileData, res)

	if err != nil {
		log.Fatal("Unable to read config json file:", err.Error())
	}

	return res
}

/*
API SPEC
	/tweet/dislike
		GET(tweetId:string):int
		POST(ownerId:string, tweetId:string, profileId:string):int
*/
