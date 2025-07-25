package main

import (
	"final13/pkg/api"
	"final13/pkg/db"
	"final13/tests"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	// Initialization of port and DB
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = strconv.Itoa(tests.Port)
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = tests.DBFile
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// .env load
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load")
	}

	// Web server
	webDir := "web"
	absWebDir, err := filepath.Abs(webDir)
	if err != nil {
		log.Fatalf("Failed to get path to web/: %v", err)
	}
	http.Handle("/", http.FileServer(http.Dir(absWebDir)))

	// API Init
	api.Init()

	// Starting server
	log.Printf("Server started at http://localhost:%s/", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server startup error: %v", err)
	}
}
