package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DislikeWriter interface {
	HandleGetDislikes(w http.ResponseWriter, r *http.Request)
	HandlePostDislikes(w http.ResponseWriter, r *http.Request)
}

/*
FILE DISLIKE WRITER
*/
type FileDislikeWriter struct {
	Path string
}

func NewFileDislikeWriter(path string) FileDislikeWriter {
	return FileDislikeWriter{
		Path: path,
	}
}

func (writer FileDislikeWriter) HandlePostDislikes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ownerId, tweetId, profileId := r.URL.Query().Get("ownerId"), r.URL.Query().Get("tweetId"), r.URL.Query().Get("profileId")
	outputObj := make(map[string]any)
	log.Printf("Disliking tweet... {disliker: %s, tweetId: %s, tweetOwner: %s}\n", profileId, tweetId, ownerId)

	defer func() {
		outputData, _ := json.MarshalIndent(outputObj, "", " ")
		fmt.Fprintf(w, "%s", string(outputData))
	}()

	fileData, err := os.ReadFile(writer.Path)
	if err != nil {
		outputObj["error"] = fmt.Sprintf("%s %s <<%s>>", "Could not open file", writer.Path, err.Error())
		log.Fatal(outputObj)
	}

	fileDataObj := make(map[string]TweetRecord)
	err = json.Unmarshal(fileData, &fileDataObj)
	if err != nil {
		outputObj["error"] = fmt.Sprintf("%s %s <<%s>>", "Could not parse data file", writer.Path, err.Error())
		log.Fatal(outputObj)
	}

	timeStamp := int(time.Now().Unix())
	if _, ok := fileDataObj[tweetId]; !ok {
		fileDataObj[tweetId] = TweetRecord{OwnerId: ownerId, TweetId: tweetId, DislikeBy: make(map[string]DislikeRecord), CreatedAt: timeStamp, UpdatedAt: timeStamp}
	}

	tweetRecord := fileDataObj[tweetId]
	tweetRecord.UpdatedAt = timeStamp
	fileDataObj[tweetId] = tweetRecord
	dislikeByMap := tweetRecord.DislikeBy
	if _, ok := dislikeByMap[profileId]; !ok {
		dislikeByMap[profileId] = DislikeRecord{OwnerId: profileId, UpdatedAt: timeStamp}
	} else {
		delete(dislikeByMap, profileId)
	}

	fileDataStr, _ := json.MarshalIndent(fileDataObj, "", " ")
	os.WriteFile(FILE_PATH, fileDataStr, os.ModeAppend|os.ModePerm)

	outputObj["dislikeCount"] = len(fileDataObj[tweetId].DislikeBy)
	outputObj["tweetId"] = tweetId
	outputObj["userDislike"] = false

	if _, ok := fileDataObj[tweetId].DislikeBy[profileId]; ok {
		outputObj["userDislike"] = true
	}

	parseDataResponse(&outputObj, &tweetRecord, profileId)
}

func (writer FileDislikeWriter) HandleGetDislikes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tweetId, profileId := r.URL.Query().Get("tweetId"), r.URL.Query().Get("profileId")
	outputObj := make(map[string]any)

	defer func() {
		outputData, _ := json.MarshalIndent(outputObj, "", " ")
		fmt.Fprintf(w, "%s", string(outputData))
	}()

	fileData, err := os.ReadFile(writer.Path)
	if err != nil {
		outputObj["error"] = fmt.Sprintf("%s %s <<%s>>", "Could not open file", writer.Path, err.Error())
		log.Fatal(outputObj)
	}

	fileDataObj := make(map[string]TweetRecord)
	err = json.Unmarshal(fileData, &fileDataObj)
	if err != nil {
		outputObj["error"] = fmt.Sprintf("%s %s <<%s>>", "Could not parse data file", writer.Path, err.Error())
		log.Fatal(outputObj)
	}

	outputObj["dislikeCount"] = len(fileDataObj[tweetId].DislikeBy)
	outputObj["tweetId"] = tweetId
	outputObj["userDislike"] = false

	if _, ok := fileDataObj[tweetId].DislikeBy[profileId]; ok {
		outputObj["userDislike"] = true
	}
}

/*
DB DISLIKE WRITER
*/
type DBDislikeWriter struct {
	Client *mongo.Client
	Coll   *mongo.Collection
}

func NewDBDislikeWriter(connstring string) DBDislikeWriter {
	client, err := mongo.Connect(options.Client().ApplyURI(connstring))
	if err != nil {
		log.Fatal("Unable to connect to mongo instance:", err.Error())
	}

	coll = client.Database("TwitterPlusDB").Collection("dislikes")

	return DBDislikeWriter{
		Client: client,
		Coll:   coll,
	}
}

func (writer DBDislikeWriter) HandlePostDislikes(w http.ResponseWriter, r *http.Request) {
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

func (writer DBDislikeWriter) HandleGetDislikes(w http.ResponseWriter, r *http.Request) {
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

/*
UTILS
*/
func parseDataResponse(outputObj *map[string]any, tweetRecord *TweetRecord, profileId string) {
	(*outputObj)["dislikeCount"] = len(tweetRecord.DislikeBy)
	(*outputObj)["tweetId"] = tweetRecord.TweetId
	(*outputObj)["userDislike"] = false

	if _, ok := tweetRecord.DislikeBy[profileId]; ok {
		(*outputObj)["userDislike"] = true
	}
}
