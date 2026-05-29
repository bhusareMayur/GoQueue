package main

import (
	"log/slog"
	"net/http"
	"os"
	"strconv" // NEW: Required for parsing int

	"github.com/joho/godotenv"

	"github.com/bhusareMayur/goqueue/internal/api"
	"github.com/bhusareMayur/goqueue/internal/api/handlers"
	"github.com/bhusareMayur/goqueue/internal/domain/job"
	redisqueue "github.com/bhusareMayur/goqueue/internal/queue/redis"
	"github.com/bhusareMayur/goqueue/internal/storage/postgres"
	"github.com/bhusareMayur/goqueue/pkg/logger"
)

func main() {
	// Initialize JSON Logger
	logger.InitJSONLogger()

	err := godotenv.Load()
	if err != nil {
		slog.Warn("no .env file found, using system environment variables")
	}

	port := os.Getenv("APP_PORT")
	postgresURL := os.Getenv("POSTGRES_URL")
	redisAddr := os.Getenv("REDIS_ADDR")

	// NEW: Read Max Queue Capacity for Backpressure
	maxQueueCapStr := os.Getenv("MAX_QUEUE_CAPACITY")
	var maxQueueCapacity int64 = 50000 // Default watermark
	if maxQueueCapStr != "" {
		if parsed, err := strconv.ParseInt(maxQueueCapStr, 10, 64); err == nil {
			maxQueueCapacity = parsed
		} else {
			slog.Warn("invalid MAX_QUEUE_CAPACITY, using default 50000", "error", err)
		}
	}

	dbPool, err := postgres.NewPool(postgresURL)
	if err != nil {
		slog.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	redisClient := redisqueue.NewClient(redisAddr)
	repo := postgres.NewJobRepository(dbPool)
	queue := redisqueue.NewQueue(redisClient)
	service := job.NewService(repo, queue)

	// NEW: Pass maxQueueCapacity to Handler
	jobHandler := handlers.NewJobHandler(service, maxQueueCapacity)
	healthHandler := handlers.NewHealthHandler(dbPool, redisClient)

	router := api.NewRouter(jobHandler, healthHandler)

	slog.Info("API server starting", "port", port, "max_queue_capacity", maxQueueCapacity)

	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}