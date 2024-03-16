package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/tursodatabase/go-libsql"
)


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbName := os.Getenv("DB_NAME")
	primaryUrl := os.Getenv("DB_URL")
	authToken := os.Getenv("DB_AUTH_TOKEN")

	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		fmt.Println("Error creating temporary directory", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, dbName)

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl, libsql.WithAuthToken(authToken))
	if err != nil {
		fmt.Println("Error creating connector", err)
		os.Exit(1)
	}
	defer connector.Close()

	db := sql.OpenDB(connector)
	defer db.Close()
}
