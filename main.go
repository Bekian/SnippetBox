package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// define a home handler func which writes a byte slice as the res body
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from SnippetBox!"))
}

// display a specfic snippet
func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		// this might use notfoundhandler later
		http.NotFound(w, r)
		return
	}

	msg := fmt.Sprintf("Displaying a snippet with ID: %d", id)
	w.Write([]byte(msg))
}

// create a snippet
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("creating a snippet"))
}

func main() {
	// create a servemux and assign the home (base/root) route
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("/snippet/view/{id}", snippetView) // add ID wildcard segment
	mux.HandleFunc("/snippet/create", snippetCreate)
	// log start
	log.Println("Starting server on :8080")
	// listen and serve
	err := http.ListenAndServe(":http", mux)
	log.Fatal(err)
}
