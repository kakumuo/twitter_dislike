package function

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("HelloWorld", helloWorld)
}

// helloWorld writes "Hello, World!" to the HTTP response.
func helloWorld(w http.ResponseWriter, r *http.Request) {
	l := log.New(funcframework.LogWriter(r.Context()), "", 0)

	l.Println("Try logging with executionID!")
	fmt.Fprintln(w, "Hello, World!")
}
