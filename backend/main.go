package main

import (
	"app/dislikes"
	"log"
	"net/http"
	"os"
)

func main() {
	writer := dislikes.NewDataStoreDislikeWriter()
	defer writer.DSClient.Close()

	http.HandleFunc("GET /tweet/dislike", writer.HandleGetDislikes)
	http.HandleFunc("POST /tweet/dislike", writer.HandlePostDislikes)

	// start server
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
