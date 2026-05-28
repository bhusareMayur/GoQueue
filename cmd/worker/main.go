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
	"github.com/bhusareMayur/goqueue/internal/reaper"
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

	// Repository, Queue, Service setup
	repo := postgres.NewJobRepository(dbPool)
	queue := redisqueue.NewQueue(redisClient)
	service := job.NewService(repo, queue)

	// Context and graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		os.Interrupt,
		syscall.SIGTERM,
	)

	go func() {
		<-signalChan
		log.Println("shutdown signal received")
		cancel()
	}()

	var wg sync.WaitGroup

	// 1. Start Delayed Job Scheduler
	sched := scheduler.NewDelayedScheduler(redisClient)
	wg.Add(1)
	go func() {
		defer wg.Done()
		sched.Start(ctx)
	}()

	// ==========================================
	// 2. NEW: Start the Reaper Service
	// ==========================================
	reaperSvc := reaper.NewReaper(service, queue)
	wg.Add(1)
	go reaperSvc.Start(ctx, &wg)

	// 3. Start Workers
	log.Printf("Starting %d workers...\n", workerConcurrency)
	for i := 1; i <= workerConcurrency; i++ {
		w := worker.NewWorker(i, queue, service)
		wg.Add(1)
		go w.Start(ctx, &wg)
	}

	// Block until everything shuts down cleanly
	wg.Wait()
	log.Println("graceful shutdown complete")
}