# GoQueue++ 🚀

> **A fault-tolerant, distributed job processing engine built in Go.**
> Built for teams who can't afford to lose a single job.

---

## What is GoQueue?

Processing everything inside a synchronous HTTP request is fragile — timeouts kill long-running work, retries duplicate side effects, and a single crash can silently swallow jobs.

GoQueue solves this by decoupling job creation from job execution. It persists every job to **PostgreSQL** the moment it's created, pushes it to **Redis** for fast async dispatch, and executes it through a pool of concurrent **Go workers** — with full crash recovery, priority scheduling, and end-to-end tracing built in.

---

## ✨ Core Features

| Feature | Description |
|---|---|
| 🗄️ **Dual-layer persistence** | PostgreSQL for durability, Redis for throughput |
| 👷 **Concurrent worker pools** | Scalable goroutine-based workers process jobs in parallel |
| 🔄 **Smart retry + backoff** | Non-blocking delayed queues with exponential backoff |
| ☠️ **Dead Letter Queue (DLQ)** | Permanently failed jobs are routed for inspection and replay |
| 💀 **Crash recovery** | Visibility timeouts + a Reaper service recover jobs from crashed workers |
| 🎯 **Priority dispatching** | High, Medium, and Low priority queues |
| 🔑 **Idempotency keys** | Safe job creation — no accidental duplicates |
| 🔍 **Correlation IDs** | Every job traceable end-to-end via structured JSON logs |
| 🛑 **Backpressure load shedding** | Returns `HTTP 429` when queue capacity is exceeded |
| 📊 **Full observability** | Prometheus metrics + Grafana dashboards out of the box |
| 🧹 **Graceful shutdown** | Workers drain in-flight jobs safely on termination signals |

---

## 🛠️ Tech Stack

```
Language        →  Go
Persistence     →  PostgreSQL
Queue / Cache   →  Redis
Observability   →  Prometheus + Grafana
Infrastructure  →  Docker + Docker Compose
```

---

## 🏗️ Architecture Overview

```
Client Request
     │
     ▼
┌─────────────┐        ┌──────────────┐
│  API Server │───────▶│  PostgreSQL  │  ← Persistent job store
└─────────────┘        └──────────────┘
     │
     ▼
┌─────────────┐
│    Redis    │  ← Fast async dispatch
└─────────────┘
     │
     ├──▶ [ High Priority Queue ]
     ├──▶ [ Medium Priority Queue ]
     └──▶ [ Low Priority Queue ]
                    │
                    ▼
          ┌──────────────────┐
          │   Worker Pool    │
          │  ┌────────────┐  │
          │  │  Worker 1  │  │
          │  │  Worker 2  │  │
          │  │  Worker 3  │  │
          └──┴────────────┴──┘
                    │
          ┌─────────┼─────────┐
          ▼         ▼         ▼
    Retry Queue   DLQ    Visibility
    (backoff)  (failed)   Timeout
                              │
                              ▼
                      Reaper Service
                      (crash recovery)
                              │
                              └──▶ Job returned to queue
```

**Observability layer** (monitors all components):

```
Prometheus ──▶ Grafana Dashboard
Structured JSON Logs + Correlation IDs
```

---

## 🚀 Quick Start

The fastest way to run GoQueue is via Docker Compose — it spins up the API, workers, PostgreSQL, Redis, Prometheus, and Grafana in one command.

### Prerequisites

- Docker & Docker Compose
- `make` (optional but recommended)

### 1. Clone the repository

```bash
git clone https://github.com/bhusaremayur/goqueue.git
cd goqueue
```

### 2. Start all services

```bash
docker-compose -f deployments/docker-compose.yml up -d
```

Or with Make:

```bash
make up
```

### 3. Verify everything is running

| Service | URL |
|---|---|
| Go API Server | `http://localhost:8080` |
| Grafana Dashboard | `http://localhost:3000` |

---

## 📁 Repository Structure

```
goqueue/
├── cmd/
│   ├── api/            → API server entry point
│   ├── worker/         → Worker process entry point
│   └── example-app/    → Example usage
│
├── internal/
│   ├── job/            → Core job domain logic
│   ├── storage/
│   │   ├── postgres/   → PostgreSQL repositories
│   │   │   └── migrations/  → SQL schema migrations
│   │   └── redis/      → Redis queue implementation
│   ├── scheduler/      → Job scheduling logic
│   └── reaper/         → Crash recovery service
│
├── pkg/
│   ├── logging/        → Structured JSON logger
│   ├── metrics/        → Prometheus metric definitions
│   ├── middleware/      → HTTP middleware (correlation IDs, etc.)
│   └── retry/          → Exponential backoff policies
│
└── deployments/
    ├── docker-compose.yml
    └── prometheus.yml
```

---

## 🔍 Observability

GoQueue is built with observability as a first-class concern.

**Structured logging** — every log line is JSON with a `correlation_id` field:

```json
{
  "level": "info",
  "correlation_id": "REQ-8472",
  "event": "job_processed",
  "job_id": "a1b2c3",
  "queue": "high",
  "worker": "worker-1",
  "duration_ms": 142
}
```

**Prometheus metrics** — queue depth, worker utilization, retry counts, DLQ size, and processing latency are all exposed at `/metrics`.

**Grafana dashboards** — pre-built dashboards ship with the repo. Access them at `http://localhost:3000` after running `docker-compose up`.

---

## ⚡ Load Testing Results

Tested with [k6](https://k6.io):

```
Virtual Users       →  10,000
Requests Tested     →  40,000+
Backpressure        →  ✓ Validated (HTTP 429 at capacity)
Crash Recovery      →  ✓ Jobs recovered correctly after worker kill
DLQ Routing         →  ✓ Permanently failed jobs isolated cleanly
```

---

## 🤝 Contributing

Contributions are welcome — bug reports, feature suggestions, and pull requests alike.

Please read [CONTRIBUTING.md](./CONTRIBUTING.md) before submitting anything. It covers how to run the test suite locally and the PR process.

---

## 📄 License

Licensed under the terms in [LICENSE](./LICENSE).

---

<div align="center">

**GoQueue v1.0**

*Built with Go · Redis · PostgreSQL · Prometheus · Grafana*

</div>