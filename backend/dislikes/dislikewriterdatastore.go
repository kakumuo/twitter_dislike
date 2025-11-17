package dislikes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
)

type Entity struct {
	Value string
}

/*
FILE DISLIKE WRITER
*/
type DataStoreDislikeWriter struct {
	Context  *context.Context
	DSClient *datastore.Client
}

func NewDataStoreDislikeWriter() DataStoreDislikeWriter {
	ctx := context.Background()

	// Create a datastore client. In a typical application, you would create
	// a single client which is reused for every datastore operation.
	dsClient, err := datastore.NewClient(ctx, "my-project")
	if err != nil {
		log.Fatal("Unable to create datastore:", err.Error())
	}
	defer dsClient.Close()

	return DataStoreDislikeWriter{
		Context:  &ctx,
		DSClient: dsClient,
	}
}

func (writer DataStoreDislikeWriter) HandleGetDislikes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ownerId, tweetId, profileId := r.URL.Query().Get("ownerId"), r.URL.Query().Get("tweetId"), r.URL.Query().Get("profileId")
	outputObj := make(map[string]any)
	log.Printf("Disliking tweet... {disliker: %s, tweetId: %s, tweetOwner: %s}\n", profileId, tweetId, ownerId)

	defer func() {
		outputData, _ := json.MarshalIndent(outputObj, "", " ")
		fmt.Fprintf(w, "%s", string(outputData))
	}()

	_, err := writer.DSClient.RunInTransaction(*writer.Context, func(tx *datastore.Transaction) error {
		filter := datastore.NameKey("TweetRecord", "tweetId", nil)
		timeStamp := int(time.Now().Unix())
		var tweetRecord TweetRecord

		err := tx.Get(filter, &tweetRecord)
		if err == datastore.ErrNoSuchEntity {
			tweetRecord = TweetRecord{OwnerId: ownerId, TweetId: tweetId, DislikeBy: make(map[string]DislikeRecord), CreatedAt: timeStamp, UpdatedAt: timeStamp}
		} else if err != nil {
			return err
		}

		tweetRecord.UpdatedAt = timeStamp
		dislikeByMap := tweetRecord.DislikeBy
		if _, ok := dislikeByMap[profileId]; !ok {
			dislikeByMap[profileId] = DislikeRecord{OwnerId: profileId, UpdatedAt: timeStamp}
		} else {
			delete(dislikeByMap, profileId)
		}

		if _, err := tx.Put(filter, &tweetRecord); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatalln("Error while running transaction", err.Error())
	}
}

func (writer DataStoreDislikeWriter) HandlePostDislikes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ownerId, tweetId, profileId := r.URL.Query().Get("ownerId"), r.URL.Query().Get("tweetId"), r.URL.Query().Get("profileId")
	outputObj := make(map[string]any)
	log.Printf("Disliking tweet... {disliker: %s, tweetId: %s, tweetOwner: %s}\n", profileId, tweetId, ownerId)

	defer func() {
		outputData, _ := json.MarshalIndent(outputObj, "", " ")
		fmt.Fprintf(w, "%s", string(outputData))
	}()

	filter := datastore.NameKey("TweetRecord", "tweetId", nil)
	timeStamp := int(time.Now().Unix())
	var tweetRecord TweetRecord

	err := writer.DSClient.Get((*writer.Context), filter, tweetRecord)
	if err == datastore.ErrNoSuchEntity {
		tweetRecord = TweetRecord{OwnerId: ownerId, TweetId: tweetId, DislikeBy: make(map[string]DislikeRecord), CreatedAt: timeStamp, UpdatedAt: timeStamp}
	} else if err != nil {
		log.Fatal("Uanble to query entity:", err.Error())
	}

	parseDataResponse(&outputObj, &tweetRecord, profileId)
}

/*

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

*/
