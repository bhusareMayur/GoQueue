package main

import (
	"log"
	"os"

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

	// Worker
	w := worker.NewWorker(queue, service)

	// Start worker (this blocks forever in the for-loop)
	w.Start()
}