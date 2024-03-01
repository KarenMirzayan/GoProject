package main

import (
	"database/sql"
	"flag"
	"github.com/KarenMirzayan/Project/pkg/messenger/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models models.Models
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8080", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://beezy:2202264mir@localhost/messenger?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	// Connect to DB
	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	app := &application{
		config: cfg,
		models: models.NewModels(db),
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	app.run()
}

func (app *application) run() {
	r := mux.NewRouter()

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Menu Singleton
	// Create a new menu
	v1.HandleFunc("/users", app.createUsersHandler).Methods("POST")
	// Get a specific menu
	v1.HandleFunc("/users/{userId:[0-9]+}", app.getUsersHandler).Methods("GET")
	////Update a specific menu
	v1.HandleFunc("/users/{userId:[0-9]+}", app.updateUsersHandler).Methods("PUT")
	// Delete a specific menu
	v1.HandleFunc("/users/{userId:[0-9]+}", app.deleteUsersHandler).Methods("DELETE")

	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	log.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config // struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
