package main

import (
	"fmt"
	"net/http"
	"transfers-svc/internal/application"
	"transfers-svc/internal/infrastructure/sqlite"

	"github.com/gocraft/dbr/v2"
	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbConn, err := dbr.Open("sqlite3", "qonto_accounts.sqlite", nil)
	if err != nil {
		log.Fatal(fmt.Errorf("unexpected error connecting to db: %w", err))
	}
	defer dbConn.Close()

	dbSession := dbConn.NewSession(nil)
	if err != nil {
		log.Fatal(fmt.Errorf("unexpected error opening db session: %w", err))
	}

	var app = application.App{
		Repository: sqlite.Repository{DbSession: dbSession},
	}
	r := mux.NewRouter()
	r.HandleFunc("/transfers", app.TransfersHandler).Methods("POST")
	http.ListenAndServe(":8081", r)
}
