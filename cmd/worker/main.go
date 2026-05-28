package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/bhusareMayur/goqueue/internal/domain/job"
	redisqueue "github.com/bhusareMayur/goqueue/internal/queue/redis"
	"github.com/bhusareMayur/goqueue/internal/reaper"
	"github.com/bhusareMayur/goqueue/internal/scheduler"
	"github.com/bhusareMayur/goqueue/internal/storage/postgres"
	"github.com/bhusareMayur/goqueue/internal/worker"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found")
	}

	postgresURL := os.Getenv("POSTGRES_URL")
	redisAddr := os.Getenv("REDIS_ADDR")
	
	workerConcurrency := 1
	if wc, err := strconv.Atoi(os.Getenv("WORKER_CONCURRENCY")); err == nil {
		workerConcurrency = wc
	}

	dbPool, err := postgres.NewPool(postgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Println("closing postgres connection")
		dbPool.Close()
	}()

	redisClient := redisqueue.NewClient(redisAddr)
	defer func() {
		log.Println("closing redis connection")
		redisClient.Close()
	}()

	repo := postgres.NewJobRepository(dbPool)
	queue := redisqueue.NewQueue(redisClient)
	service := job.NewService(repo, queue)

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Println("shutdown signal received")
		cancel()
	}()

	var wg sync.WaitGroup

	// ==========================================
	// NEW: Start Worker Metrics Server
	// ==========================================
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		log.Println("worker metrics server running on port :2112")
		if err := http.ListenAndServe(":2112", mux); err != nil {
			log.Printf("metrics server error: %v\n", err)
		}
	}()

	sched := scheduler.NewDelayedScheduler(redisClient)
	wg.Add(1)
	go func() {
		defer wg.Done()
		sched.Start(ctx)
	}()

	reaperSvc := reaper.NewReaper(service, queue)
	wg.Add(1)
	go reaperSvc.Start(ctx, &wg)

	log.Printf("Starting %d workers...\n", workerConcurrency)
	for i := 1; i <= workerConcurrency; i++ {
		w := worker.NewWorker(i, queue, service)
		wg.Add(1)
		go w.Start(ctx, &wg)
	}

	wg.Wait()
	log.Println("graceful shutdown complete")
}