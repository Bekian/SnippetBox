package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {
	/// config
	// addr port string uses the 4000 port by default
	addr := flag.String("addr", ":8080", "HTTP network address")
	flag.Parse()
	// logger init
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		logger: logger,
	}

	// server
	logger.Info("starting server", "addr", *addr)
	err := http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
