package dislikes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

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
	os.WriteFile(writer.Path, fileDataStr, os.ModeAppend|os.ModePerm)

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
