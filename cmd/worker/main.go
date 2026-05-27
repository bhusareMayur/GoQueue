package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/bhusareMayur/goqueue/internal/domain/job"
	redisqueue "github.com/bhusareMayur/goqueue/internal/queue/redis"
	"github.com/bhusareMayur/goqueue/internal/scheduler"
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
	
	workerConcurrency := 1
	if wc, err := strconv.Atoi(os.Getenv("WORKER_CONCURRENCY")); err == nil {
		workerConcurrency = wc
	}

	// PostgreSQL connection
	dbPool, err := postgres.NewPool(postgresURL)
	if err != nil {
		log.Fatal(err)
	}
	
	// STEP 8: Graceful Resource Cleanup
	defer func() {
		log.Println("closing postgres connection")
		dbPool.Close()
	}()

	// Redis connection
	redisClient := redisqueue.NewClient(redisAddr)
	defer func() {
		log.Println("closing redis connection")
		redisClient.Close()
	}()

	// Repository
	repo := postgres.NewJobRepository(dbPool)

	// Queue
	queue := redisqueue.NewQueue(redisClient)

	// Service
	service := job.NewService(repo, queue)

	// STEP 1: Root Context
	ctx, cancel := context.WithCancel(context.Background())

	// STEP 2: OS Signal Channel
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		os.Interrupt,
		syscall.SIGTERM,
	)

	// STEP 3: Shutdown Goroutine
	go func() {
		<-signalChan
		log.Println("shutdown signal received")
		cancel()
	}()

	// STEP 7: WaitGroup
	var wg sync.WaitGroup

	// ==========================================
	// NEW: Start the Delayed Job Scheduler
	// ==========================================
	sched := scheduler.NewDelayedScheduler(redisClient)
	wg.Add(1)
	go func() {
		// Pass the waitgroup down so it shuts down gracefully
		defer wg.Done()
		sched.Start(ctx)
	}()

	log.Printf("Starting %d workers...\n", workerConcurrency)

	// Start workers dynamically based on ENV variable
	for i := 1; i <= workerConcurrency; i++ {
		w := worker.NewWorker(i, queue, service)
		wg.Add(1)
		
		// Run worker in a goroutine
		go w.Start(ctx, &wg)
	}

	// Main thread blocks here until all workers execute wg.Done()
	wg.Wait()

	log.Println("graceful shutdown complete")
}