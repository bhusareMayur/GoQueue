package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/bhusareMayur/goqueue/internal/api"
	"github.com/bhusareMayur/goqueue/internal/api/handlers"
	"github.com/bhusareMayur/goqueue/internal/domain/job"
	redisqueue "github.com/bhusareMayur/goqueue/internal/queue/redis"
	"github.com/bhusareMayur/goqueue/internal/storage/postgres"
)

func main() {

	// Load .env
	err := godotenv.Load()

	if err != nil {
		log.Println("no .env file found")
	}

	// Read environment variables
	port := os.Getenv("APP_PORT")

	postgresURL := os.Getenv("POSTGRES_URL")

	redisAddr := os.Getenv("REDIS_ADDR")

	// PostgreSQL connection
	dbPool, err := postgres.NewPool(
		postgresURL,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer dbPool.Close()

	// Redis connection
	redisClient := redisqueue.NewClient(
		redisAddr,
	)

	// Repository
	repo := postgres.NewJobRepository(
		dbPool,
	)

	// Queue
	queue := redisqueue.NewQueue(
		redisClient,
	)

	// Service
	service := job.NewService(
		repo,
		queue,
	)

	// Handlers
	jobHandler := handlers.NewJobHandler(
		service,
	)

	// NEW: Initialize Health Handler
	healthHandler := handlers.NewHealthHandler(
		dbPool,
		redisClient,
	)

	// Router
	router := api.NewRouter(
		jobHandler,
		healthHandler, // NEW: Inject into router
	)

	log.Println("API server running on port :", port)

	// Start server
	err = http.ListenAndServe(
		":"+port,
		router,
	)

	if err != nil {
		log.Fatal(err)
	}
}