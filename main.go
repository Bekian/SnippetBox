package main

import (
	"log"
	"net/http"
)

// define a home handler func which writes a byte slice as the res body
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from SnippetBox!")) // here the casing is slightly different cause i think it looks nicer
}

// display a specfic snippet
func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("a specific snippet"))
}

func main() {
	// create a servemux and assign the home (base/root) route
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	// log start
	log.Println("Starting server on :8080")
	// listen and serve
	err := http.ListenAndServe(":http", mux)
	log.Fatal(err)
}
