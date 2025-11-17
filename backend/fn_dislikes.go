package dislikes

import (
	"net/http"
	"sync"

	"company.com/app/dislikes"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

var writerOnce sync.Once
var writer dislikes.DataStoreDislikeWriter

func init() {
	functions.HTTP("TwitterPlus", TwitterPlusRouter)
}

// helloHTTP is an HTTP Cloud Function with a request parameter.
func TwitterPlusRouter(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tweet/dislike" {
		http.Error(w, "Invalid URL Path", http.StatusNotFound)
	} else {
		writerOnce.Do(func() {
			writer = dislikes.NewDataStoreDislikeWriter()
		})

		// Set CORS headers for the preflight request
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// Set CORS headers for the main request.
		w.Header().Set("Access-Control-Allow-Origin", "*")

		switch r.Method {
		case http.MethodGet:
			writer.HandleGetDislikes(w, r)
		case http.MethodPost:
			writer.HandlePostDislikes(w, r)
		default:
			http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}
