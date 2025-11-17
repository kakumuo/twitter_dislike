package dislikes

import (
	"net/http"
)

type DislikeWriter interface {
	HandleGetDislikes(w http.ResponseWriter, r *http.Request)
	HandlePostDislikes(w http.ResponseWriter, r *http.Request)
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
