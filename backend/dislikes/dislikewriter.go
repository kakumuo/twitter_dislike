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
func parseDataResponse(outputObj *map[string]any, tweetRecord *TweetRecord, profileIndex int) {
	(*outputObj)["dislikeCount"] = len(tweetRecord.DislikeBy)
	(*outputObj)["tweetId"] = tweetRecord.TweetId
	(*outputObj)["userDislike"] = false

	if profileIndex == -1 {
		(*outputObj)["userDislike"] = true
	}
}
