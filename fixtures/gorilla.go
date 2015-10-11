package main

import (
	"log"
	"net/http"

	"github.com/gernest/helen"
	"github.com/gorilla/mux"
)

func main() {

	// Create the instance of your router
	server := mux.NewRouter()

	// Create a new helen.Static instance. We are passing "static" as the directory
	// we want to serve static content from.
	static := helen.NewStatic("fixtures")

	// We bind anything matching /static/ route to our static handler
	static.Bind("/static/", server)

	// You can register other handlers to your server or whatever you want to do with it.

	log.Fatal(http.ListenAndServe(":8000", server))
}
