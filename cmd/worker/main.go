package main

import (
	"context"
	"log/slog"
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
	"github.com/bhusareMayur/goqueue/pkg/logger"
)

func main() {
	// Initialize JSON Logger
	logger.InitJSONLogger()

	err := godotenv.Load()
	if err != nil {
		slog.Warn("no .env file found, using system environment variables")
	}

	postgresURL := os.Getenv("POSTGRES_URL")
	redisAddr := os.Getenv("REDIS_ADDR")
	
	workerConcurrency := 1
	if wc, err := strconv.Atoi(os.Getenv("WORKER_CONCURRENCY")); err == nil {
		workerConcurrency = wc
	}

	dbPool, err := postgres.NewPool(postgresURL)
	if err != nil {
		slog.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer func() {
		slog.Info("closing postgres connection")
		dbPool.Close()
	}()

	redisClient := redisqueue.NewClient(redisAddr)
	defer func() {
		slog.Info("closing redis connection")
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
		slog.Info("shutdown signal received")
		cancel()
	}()

	var wg sync.WaitGroup

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		slog.Info("worker metrics server running", "port", "2112")
		if err := http.ListenAndServe(":2112", mux); err != nil {
			slog.Error("metrics server error", "error", err)
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

	slog.Info("starting workers", "concurrency", workerConcurrency)
	for i := 1; i <= workerConcurrency; i++ {
		w := worker.NewWorker(i, queue, service)
		wg.Add(1)
		go w.Start(ctx, &wg)
	}

	wg.Wait()
	slog.Info("graceful shutdown complete")
}