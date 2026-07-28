package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/CR45-NITT/cr45-reduced/backend/internal/api"
	"github.com/CR45-NITT/cr45-reduced/backend/internal/auth"
	"github.com/CR45-NITT/cr45-reduced/backend/internal/store"
	"github.com/CR45-NITT/cr45-reduced/backend/migrations"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.LUTC)
	addr := envOrDefault("HTTP_ADDR", ":8080")
	databaseURL := envOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/cr45_reduced?sslmode=disable")
	appSecret := envOrDefault("APP_SECRET", "cr45-reduced-dev-secret")
	logRequests := envOrDefault("LOG_REQUESTS", "") == "true"

	db := mustOpenDB(databaseURL)
	defer db.Close()
	if err := migrations.Up(db); err != nil {
		logger.Fatalf("migrations failed: %v", err)
	}

	repo := store.NewPostgresStore(db)
	authService, err := auth.NewService(repo, appSecret)
	if err != nil {
		logger.Fatalf("auth init failed: %v", err)
	}
	server := api.NewServer(repo, authService, logger, logRequests)

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           server.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Printf("starting cr45-reduced backend on %s", addr)
	if err := httpServer.ListenAndServe(); err != nil {
		logger.Fatalf("server stopped: %v", err)
	}
}

func mustOpenDB(databaseURL string) *sql.DB {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatalf("db open failed: %v", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	const (
		maxRetries = 20
		retryDelay = 2 * time.Second
	)

	for i := 1; i <= maxRetries; i++ {
		err = db.Ping()

		if err == nil {
			log.Println("Connected to PostgreSQL successfully!")
			return db
		}

		log.Printf("Database not ready (attempt %d/%d): %v", i, maxRetries, err)

		time.Sleep(retryDelay)
	}

	log.Fatalf("Could not connect to PostgreSQL after %d attempts", maxRetries)
	return nil
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
