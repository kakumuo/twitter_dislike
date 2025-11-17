package dislikes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

	var projId string
	envVars := []string{"DEVSHELL_PROJECT_ID", "GOOGLE_CLOUD_PROJECT_ID"}
	for _, envVar := range envVars {
		projId = os.Getenv(envVar)
		if projId != "" {
			break
		}
	}

	if projId == "" {
		projId = "twitter-plus-38c93"
	}

	// Create a datastore client. In a typical application, you would create
	// a single client which is reused for every datastore operation.
	dsClient, err := datastore.NewClient(ctx, projId)
	if err != nil {
		log.Fatal("Unable to create datastore:", err.Error())
	}

	return DataStoreDislikeWriter{
		Context:  &ctx,
		DSClient: dsClient,
	}
}

func (writer DataStoreDislikeWriter) HandlePostDislikes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ownerId, tweetId, profileId := r.URL.Query().Get("ownerId"), r.URL.Query().Get("tweetId"), r.URL.Query().Get("profileId")
	apiResponse := APIResponse{}
	log.Printf("Disliking tweet... {disliker: %s, tweetId: %s, tweetOwner: %s}\n", profileId, tweetId, ownerId)

	defer func() {
		outputData, _ := json.MarshalIndent(apiResponse, "", " ")
		fmt.Fprintf(w, "%s", string(outputData))
	}()

	_, err := writer.DSClient.RunInTransaction(*writer.Context, func(tx *datastore.Transaction) error {
		tweetKey := datastore.NameKey("TweetRecord", tweetId, nil)
		timeStamp := int(time.Now().Unix())
		var tweetRecord TweetRecord

		err := tx.Get(tweetKey, &tweetRecord)
		if err == datastore.ErrNoSuchEntity {
			tweetRecord = TweetRecord{OwnerId: ownerId, TweetId: tweetId, CreatedAt: timeStamp, UpdatedAt: timeStamp}
			tx.Put(tweetKey, &tweetRecord)
		} else if err != nil {
			return err
		}

		dislikeKey := datastore.NameKey("DislikeRecord", profileId, tweetKey)
		var dislikeRec DislikeRecord
		err = tx.Get(dislikeKey, &dislikeRec)
		switch err {
		case datastore.ErrNoSuchEntity:
			dislikeRec = DislikeRecord{OwnerId: profileId, UpdatedAt: timeStamp}
			_, err = tx.Put(dislikeKey, &dislikeRec)
			if err != nil {
				return err
			}
			apiResponse.UserDislike = true
		case nil:
			tx.Delete(dislikeKey)
		default:
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatalln("Error while running transaction =>", err.Error())
	}

	tweetKey := datastore.NameKey("TweetRecord", tweetId, nil)
	countQuery := datastore.NewQuery("DislikeRecord").Ancestor(tweetKey)
	n, err := writer.DSClient.Count(*writer.Context, countQuery)

	if err != nil {
		log.Fatal("Unable to count records:", err.Error())
	}

	apiResponse.DislikeCount = n
	apiResponse.TweetId = tweetId
}

func (writer DataStoreDislikeWriter) HandleGetDislikes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tweetId, profileId := r.URL.Query().Get("tweetId"), r.URL.Query().Get("profileId")
	apiResponse := APIResponse{}
	log.Printf("Disliking tweet... {disliker: %s, tweetId: %s}\n", profileId, tweetId)

	defer func() {
		outputData, _ := json.MarshalIndent(apiResponse, "", " ")
		fmt.Fprintf(w, "%s", string(outputData))
	}()

	tweetKey := datastore.NameKey("TweetRecord", tweetId, nil)
	dislikeKey := datastore.NameKey("DislikeRecord", profileId, tweetKey)

	hasDislikeQuery := datastore.NewQuery("DislikeRecord").FilterField("__key__", "=", dislikeKey)
	n, err := writer.DSClient.Count(*writer.Context, hasDislikeQuery)

	if err != nil {
		log.Fatal("Unable to count user dislike:", err.Error())
	}
	apiResponse.UserDislike = n > 0

	countQuery := datastore.NewQuery("DislikeRecord").Ancestor(tweetKey)
	n, err = writer.DSClient.Count(*writer.Context, countQuery)

	if err != nil {
		log.Fatal("Unable to count records:", err.Error())
	}
	apiResponse.DislikeCount = n
	apiResponse.TweetId = tweetId
}
