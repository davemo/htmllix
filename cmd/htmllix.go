package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/davemo/htmllix/pkg/view"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	serverPort := os.Getenv("SERVER_PORT")

	clerkEnv := view.ClerkEnv{
		PublishableKey: os.Getenv("CLERK_PUBLISHABLE_KEY"),
		FrontendApi: os.Getenv("CLERK_FRONTEND_API"),
	}

	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		fmt.Println("Error creating temporary directory", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, dbName)
	syncInterval := time.Minute

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl, libsql.WithAuthToken(authToken), libsql.WithSyncInterval(syncInterval))
	if err != nil {
		fmt.Println("Error creating connector", err)
		os.Exit(1)
	}
	defer connector.Close()

	db := sql.OpenDB(connector)
	defer db.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	// e.Static("/public", "public")

	e.GET("/", func(c echo.Context) error {
		index := view.Index()
		layout := view.Layout(index, clerkEnv)
		return layout.Render(context.Background(), c.Response().Writer)
	})

	e.Logger.Fatal(e.Start(serverPort))
}
