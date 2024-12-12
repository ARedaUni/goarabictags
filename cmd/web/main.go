package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/amrojjeh/arabic-tags/internal/models"
	_ "github.com/lib/pq"
)

type application struct {
	logger     *slog.Logger
	u          url
	page       map[string]*template.Template
	user       models.UserModel
	excerpt    models.ExcerptModel
	word       models.WordModel
	manuscript models.ManuscriptModel
	session    *scs.SessionManager
	mutex      sync.Mutex
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP Address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"))
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("cannot open db", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Error("cannot open connection with db", slog.String("error", err.Error()))
		os.Exit(1)
	}

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Store = postgresstore.New(db)

	app := application{
		logger:     logger,
		user:       models.UserModel{Db: db},
		excerpt:    models.ExcerptModel{Db: db},
		manuscript: models.ManuscriptModel{Db: db},
		word:       models.WordModel{Db: db},
		session:    session,
	}

	server := &http.Server{
		Handler:      app.routes(),
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("starting server", slog.String("addr", *addr))
	err = server.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
