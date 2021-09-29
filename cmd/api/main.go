package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"github.com/jwaldner/go-movies-backend/models"
	_ "github.com/lib/pq"
)

type AppStatus struct {
	Status  string `json:"status"`
	Port    int    `json:"listening on port"`
	Env     string `json:"environment"`
	Version string `json:"version"`
}

type application struct {
	config config
	logger *log.Logger
	models models.Models
}

var (
	// Get current file full path from runtime
	_, b, _, _ = runtime.Caller(0)

	// projectRootPath folder of this project
	projectRootPath = filepath.Join(filepath.Dir(b), "../")
)

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment (development/production)")
	flag.StringVar(&cfg.db.user, "POSTGRES_USER", "movie_user", "database user")
	flag.StringVar(&cfg.db.password, "POSTGRES_PASSWORD", "magic_ticket", "database password")
	flag.StringVar(&cfg.db.host, "POSTGRES_HOST", "localhost", "database host")
	flag.StringVar(&cfg.db.defaultDb, "POSTGRES_DB", "movie_database", "database")

	//flag.StringVar(&cfg.db.dsn, "dsn", "postgres://movie_user:magic_ticket@localhost/movie_database?sslmode=disable", "postgres connection string")
	//flag.StringVar(&cfg.db.dsn, "dsn", "postgres://movie_user:magic_ticket@db/movie_database?sslmode=disable", "postgres connection string")

	flag.Parse()

	if os.Getenv("POSTGRES_HOST") != "" {
		cfg.db.user = os.Getenv("POSTGRES_USER")
		cfg.db.password = os.Getenv("POSTGRES_PASSWORD")
		cfg.db.host = os.Getenv("POSTGRES_HOST")
		cfg.db.defaultDb = os.Getenv("POSTGRES_DB")
	}

	url := fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable",
        cfg.db.user,
        cfg.db.password,
        cfg.db.host,
        cfg.db.defaultDb)

	cfg.db.dsn = url

	// load .env file
	err := godotenv.Load(projectRootPath + "/.env")

	if err != nil {
		log.Printf("Error .env file: %s", err)
		//log.Fatalf("Error  .env file%s", err)
	}

	if os.Getenv("GO_MOVIES_JWT") == "" {
		log.Fatalf("environment variables not set")
	} else {
		cfg.jwt.secret = os.Getenv("GO_MOVIES_JWT")
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDb(cfg)

	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	app := application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	logger.Printf("version %s listening on port %v in '%s' mode\n", version, cfg.port, cfg.env)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err = srv.ListenAndServe()

	if err != nil {
		app.logger.Fatal(err.Error())
	}
}

func openDb(cfg config) (*sql.DB, error) {

	db, err := sql.Open("postgres", cfg.db.dsn)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
