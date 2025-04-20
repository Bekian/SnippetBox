package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Bekian/SnippetBox/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecorder   *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	/// init configs
	// addr port string uses the 8080 port by default
	addr := flag.String("addr", ":8080", "HTTP network address")
	// dsn string for mysql connection
	dsn := flag.String("dsn", "web:1234@/snippetbox?parseTime=true", "MySQL data source name string")
	flag.Parse()

	// logger init
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// init db
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// init template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// init form decoder
	formDecoder := form.NewDecoder()

	// init session manager
	// use mysql db as session store and a lifetime of 12 hours
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// pass above initialized objects to app struct for outside use
	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecorder:   formDecoder,
		sessionManager: sessionManager,
	}

	// init the server struct
	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	// logger
	logger.Info("starting server", "addr", *addr)
	// start the server
	err = srv.ListenAndServe()
	// pass any errors that arise to our logger
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
