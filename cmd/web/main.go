package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	/// config
	// addr port string uses the 4000 port by default
	addr := flag.String("addr", ":8080", "HTTP network address")
	flag.Parse()
	// logger init
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// router
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// serve static files
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	// server
	logger.Info("starting server", "addr", *addr)
	err := http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
