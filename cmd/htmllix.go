package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/davemo/htmllix/pkg/view"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tursodatabase/go-libsql"
)

func ClerkAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionToken, err := c.Cookie("__session")
		if err != nil || sessionToken.Value == "" {
			fmt.Println("No Clerk session cookie found. Redirecting to /")
			// logging all request headers
			for name, headers := range c.Request().Header {
				for _, h := range headers {
					fmt.Printf("%v: %v\n", name, h)
				}
			}
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		pemPublicKey := os.Getenv("CLERK_PEM_PUBLIC_KEY")
		jwk, err := clerk.JSONWebKeyFromPEM(pemPublicKey)
		if err != nil {
			fmt.Printf("error parsing public key: %v. Redirecting to /\n", err)
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		claims, err := jwt.Verify(c.Request().Context(), &jwt.VerifyParams{
			Token: sessionToken.Value,
			JWK: jwk,
		})
		if err != nil {
			fmt.Printf("session token validation failed: %v. Redirecting to /\n", err)
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		fmt.Printf("User with ID: %s is authenticated\n", claims.Subject)
		c.Set("claims", claims)
		return next(c)
	}
}

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
		SecretKey: os.Getenv("CLERK_SECRET_KEY"),
	}
	clerk.SetKey(clerkEnv.SecretKey)

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
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowCredentials: true,
	// 	AllowOrigins: []string{"*"},
	// }))
	e.Use(middleware.Logger())
	// e.Static("/public", "public")

	e.GET("/", func(c echo.Context) error {
		index := view.Index()
		layout := view.Layout(index, clerkEnv)
		return layout.Render(context.Background(), c.Response().Writer)
	})

	e.GET("/home", ClerkAuthMiddleware(func(e echo.Context) error {
		home := view.Home()
		layout := view.Layout(home, clerkEnv)
		return layout.Render(context.Background(), e.Response().Writer)
	}))

	e.Logger.Fatal(e.Start(serverPort))
}
