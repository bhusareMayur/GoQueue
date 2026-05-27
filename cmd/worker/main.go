package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/bhusareMayur/goqueue/internal/domain/job"
	redisqueue "github.com/bhusareMayur/goqueue/internal/queue/redis"
	"github.com/bhusareMayur/goqueue/internal/storage/postgres"
	"github.com/bhusareMayur/goqueue/internal/worker"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found")
	}

	// Read environment variables
	postgresURL := os.Getenv("POSTGRES_URL")
	redisAddr := os.Getenv("REDIS_ADDR")
	concurrencyStr := os.Getenv("WORKER_CONCURRENCY")

	// Parse concurrency, default to 5 if not set or invalid
	concurrency, err := strconv.Atoi(concurrencyStr)
	if err != nil || concurrency <= 0 {
		log.Printf("invalid or missing WORKER_CONCURRENCY, defaulting to 5")
		concurrency = 5
	}

	// PostgreSQL connection
	dbPool, err := postgres.NewPool(postgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	// Redis connection
	redisClient := redisqueue.NewClient(redisAddr)

	// Repository
	repo := postgres.NewJobRepository(dbPool)

	// Queue
	queue := redisqueue.NewQueue(redisClient)

	// Service
	service := job.NewService(repo, queue)

	log.Printf("starting worker pool with concurrency: %d", concurrency)

	// Spawn worker goroutines
	for i := 1; i <= concurrency; i++ {
		w := worker.NewWorker(i, queue, service)
		go w.Start() // Start each worker in a separate goroutine
	}

	// Block main thread from exiting forever
	select {}
}