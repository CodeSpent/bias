package http

import (
	"bias/config"
	"bias/http/router"
	"bias/pkg/database/postgres"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"os"
	"os/signal"
	"syscall"
)

var e = echo.New()

func StartServer() {

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, error=${error}, time=${time_rfc3339}, remote_ip=${remote_ip}, path=${path}, query=${query}, request=${request}\n",
	}))

	// Initialize configuration
	err := config.Initialize()
	if err != nil {
		panic("Failed to initialize configuration: " + err.Error())
	}
	defer config.Close()

	godotenv.Load()

	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable TimeZone=America/New_York",
		host, port, username, dbname, password,
	)

	// Connect to PostgreSQL database
	db, err := postgres.NewPostgresDB(connStr)

	if err != nil {
		log.Fatal(err)
	}

	// Migrate models

	// Set up routes
	// TODO: This does more than just setup routes.
	router.SetupRoutes(e, config.RedisClient, db)

	// Start the server in a separate goroutine
	go func() {
		if err := e.Start(":5032"); err != nil {
			e.Logger.Fatal(err)
		}
	}()

	// Handle SIGINT and SIGTERM signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Shutdown the server
	if err := e.Shutdown(context.Background()); err != nil {
		e.Logger.Fatal(err)
	}
}

func StopServer() {
	if e != nil {
		fmt.Println("Stopping the server...")
		if err := e.Shutdown(context.Background()); err != nil {
			e.Logger.Fatal(err)
		}
	}
}
