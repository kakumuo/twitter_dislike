package dislikes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

/*
Mongo DISLIKE WRITER
*/
type MongoDislikeWriter struct {
	Client *mongo.Client
	Coll   *mongo.Collection
}

func NewMongoDislikeWriter(connstring string) MongoDislikeWriter {
	client, err := mongo.Connect(options.Client().ApplyURI(connstring))
	if err != nil {
		log.Fatal("Unable to connect to mongo instance:", err.Error())
	}

	coll := client.Database("TwitterPlusDB").Collection("dislikes")

	return MongoDislikeWriter{
		Client: client,
		Coll:   coll,
	}
}

func (writer MongoDislikeWriter) HandlePostDislikes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ownerId, tweetId, profileId := r.URL.Query().Get("ownerId"), r.URL.Query().Get("tweetId"), r.URL.Query().Get("profileId")
	outputObj := make(map[string]any)
	log.Printf("Disliking tweet... {disliker: %s, tweetId: %s, tweetOwner: %s}\n", profileId, tweetId, ownerId)

	defer func() {
		outputData, _ := json.MarshalIndent(outputObj, "", " ")
		fmt.Fprintf(w, "%s", string(outputData))
	}()

	filter := bson.D{{"tweetId", tweetId}}
	cursor := writer.Coll.FindOne(context.TODO(), filter)
	err := cursor.Err()
	var tweetRecord TweetRecord

	timeStamp := int(time.Now().Unix())

	if err == mongo.ErrNoDocuments {
		tweetRecord = TweetRecord{OwnerId: ownerId, TweetId: tweetId, DislikeBy: make(map[string]DislikeRecord), CreatedAt: timeStamp, UpdatedAt: timeStamp}
	} else if err != nil {
		log.Fatal("Uanble to query document:", err.Error())
	} else {
		cursor.Decode(&tweetRecord)
	}

	tweetRecord.UpdatedAt = timeStamp
	dislikeByMap := tweetRecord.DislikeBy
	if _, ok := dislikeByMap[profileId]; !ok {
		dislikeByMap[profileId] = DislikeRecord{OwnerId: profileId, UpdatedAt: timeStamp}
	} else {
		delete(dislikeByMap, profileId)
	}

	filter = bson.D{{"tweetId", tweetId}}
	update := bson.D{{"$set", tweetRecord}}
	opts := options.UpdateOne().SetUpsert(true)

	result, err := writer.Coll.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Fatal("Unable to update document:", err.Error())
	}

	log.Printf("Number of documents updated: %v\n", result.ModifiedCount)
	log.Printf("Number of documents upserted: %v\n", result.UpsertedCount)

	parseDataResponse(&outputObj, &tweetRecord, profileId)
}

func (writer MongoDislikeWriter) HandleGetDislikes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ownerId, tweetId, profileId := r.URL.Query().Get("ownerId"), r.URL.Query().Get("tweetId"), r.URL.Query().Get("profileId")
	outputObj := make(map[string]any)
	log.Printf("Disliking tweet... {disliker: %s, tweetId: %s, tweetOwner: %s}\n", profileId, tweetId, ownerId)

	defer func() {
		outputData, _ := json.MarshalIndent(outputObj, "", " ")
		fmt.Fprintf(w, "%s", string(outputData))
	}()

	filter := bson.D{{"tweetId", tweetId}}
	cursor := writer.Coll.FindOne(context.TODO(), filter)
	err := cursor.Err()
	var tweetRecord TweetRecord

	timeStamp := int(time.Now().Unix())

	if err == mongo.ErrNoDocuments {
		tweetRecord = TweetRecord{OwnerId: ownerId, TweetId: tweetId, DislikeBy: make(map[string]DislikeRecord), CreatedAt: timeStamp, UpdatedAt: timeStamp}
	} else if err != nil {
		log.Fatal("Uanble to query document:", err.Error())
	} else {
		cursor.Decode(&tweetRecord)
	}

	parseDataResponse(&outputObj, &tweetRecord, profileId)
}
