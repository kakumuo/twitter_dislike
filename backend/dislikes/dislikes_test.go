package dislikes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const HOST = "127.0.0.1"
const PORT = 8888
const FILE_PATH = "./data.json"

var coll *mongo.Collection

// TODO: move functionality to mongodb
func TestMongo(t *testing.T) {
	mux := http.NewServeMux()

	var mongoWriter MongoDislikeWriter
	config := LoadConfig("./config.json", t)

	mongoWriter = NewMongoDislikeWriter(config.ConnectionString)
	defer func() {
		if err := mongoWriter.Client.Disconnect(context.TODO()); err != nil {
			t.Fatal("Error while disconnecting from mongo instance:", err.Error())
		}
	}()

	mux.HandleFunc("GET /tweet/dislike", mongoWriter.HandleGetDislikes)
	mux.HandleFunc("POST /tweet/dislike", mongoWriter.HandlePostDislikes)

	// start server
	corsHandler := cors.AllowAll().Handler(mux)
	serverEndpoint := fmt.Sprintf("%s:%d", config.HostName, config.Port)
	t.Log("Starting server at:", serverEndpoint)
	t.Fatal(http.ListenAndServe(serverEndpoint, corsHandler))
}

func LoadConfig(filepath string, t *testing.T) *Config {
	res := &Config{}

	fileData, err := os.ReadFile(filepath)

	if err != nil {
		t.Fatal("Error: ", err.Error())
	}

	err = json.Unmarshal(fileData, res)

	if err != nil {
		t.Fatal("Unable to read config json file:", err.Error())
	}

	return res
}
