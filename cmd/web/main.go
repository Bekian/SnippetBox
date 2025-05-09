package main

import (
	"crypto/tls"
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
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
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
	// also use set to use tls
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// pass above initialized objects to app struct for outside use
	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecorder:   formDecoder,
		sessionManager: sessionManager,
	}

	// only use these 2 curves for assembly
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// init the server struct
	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
		// this writes any http server errors to the custom logger
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// logger
	logger.Info("starting server", "addr", *addr)
	// start the server with tls cert and key
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
