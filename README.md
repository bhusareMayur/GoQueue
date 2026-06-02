<div align="center">

<img src="https://readme-typing-svg.demolab.com?font=JetBrains+Mono&weight=700&size=40&pause=1000&color=00ADD8&center=true&vCenter=true&width=600&height=80&lines=GoQueue%2B%2B+🚀;Fault-Tolerant+Job+Engine;Zero+Job+Loss+Guarantee" alt="GoQueue++ Typing SVG" />

<br/>

**A fault-tolerant, distributed job processing engine built in Go.**
*Built for teams who can't afford to lose a single job.*

<br/>

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=for-the-badge&logo=postgresql&logoColor=white)](https://postgresql.org)
[![Redis](https://img.shields.io/badge/Redis-7-FF4438?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io)
[![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=for-the-badge&logo=prometheus&logoColor=white)](https://prometheus.io)
[![Grafana](https://img.shields.io/badge/Grafana-F46800?style=for-the-badge&logo=grafana&logoColor=white)](https://grafana.com)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://docker.com)

<br/>

[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](./LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-Welcome-brightgreen?style=flat-square)](./CONTRIBUTING.md)
[![Built in Public](https://img.shields.io/badge/Built_in_Public-12_posts-00ADD8?style=flat-square)](https://github.com/bhusaremayur/goqueue)
[![Version](https://img.shields.io/badge/version-1.0.0-blue?style=flat-square)](https://github.com/bhusaremayur/goqueue)

</div>

---

## 🤔 Why GoQueue?

Processing everything inside a synchronous HTTP request is **fragile**:

```
❌ Timeouts kill long-running work
❌ Retries duplicate side effects
❌ A single crash silently swallows jobs
```

GoQueue solves this by **decoupling job creation from job execution**:

```
✅ Persists every job to PostgreSQL the moment it's created
✅ Pushes to Redis for fast async dispatch
✅ Executes through concurrent Go workers
✅ Full crash recovery, priority scheduling, and tracing built in
```

---

## ✨ Core Features

<div align="center">

| Feature | Description |
|:---:|:---|
| 🗄️ **Dual-layer persistence** | PostgreSQL for durability · Redis for throughput |
| 👷 **Concurrent worker pools** | Scalable goroutine-based workers process jobs in parallel |
| 🔄 **Smart retry + backoff** | Non-blocking delayed queues with exponential backoff |
| ☠️ **Dead Letter Queue** | Permanently failed jobs routed for inspection and replay |
| 💀 **Crash recovery** | Visibility timeouts + Reaper service recover stuck jobs |
| 🎯 **Priority dispatching** | High · Medium · Low priority queues |
| 🔑 **Idempotency keys** | Safe job creation — no accidental duplicates |
| 🔍 **Correlation IDs** | Every job traceable end-to-end via structured JSON logs |
| 🛑 **Backpressure shedding** | Returns `HTTP 429` when queue capacity is exceeded |
| 📊 **Full observability** | Prometheus metrics + Grafana dashboards out of the box |
| 🧹 **Graceful shutdown** | Workers drain in-flight jobs safely on termination signals |

</div>

---

## 🏗️ Architecture

```
╔══════════════════════════════════════════════════════════════╗
║                      CLIENT REQUEST                         ║
╚══════════════════════════╤═══════════════════════════════════╝
                           │
                           ▼
              ┌────────────────────────┐
              │      API Server        │  :8080
              └────────┬───────────────┘
                       │
           ┌───────────┼───────────┐
           ▼                       ▼
  ┌────────────────┐    ┌────────────────────┐
  │  PostgreSQL    │    │       Redis        │
  │  (durable      │    │  (fast dispatch)   │
  │   job store)   │    └──────────┬─────────┘
  └────────────────┘               │
                        ┌──────────┼──────────┐
                        ▼          ▼           ▼
                   ┌─────────┐ ┌───────┐ ┌────────┐
                   │  HIGH   │ │  MED  │ │  LOW   │
                   │ priority│ │  pri  │ │  pri   │
                   └────┬────┘ └───┬───┘ └───┬────┘
                        └──────────┼──────────┘
                                   ▼
                    ┌──────────────────────────┐
                    │        Worker Pool        │
                    │  ┌────────┐  ┌────────┐  │
                    │  │Worker 1│  │Worker 2│  │
                    │  └────────┘  └────────┘  │
                    │        ┌────────┐         │
                    │        │Worker N│         │
                    │        └────────┘         │
                    └──────────────────────────┘
                                   │
               ┌───────────────────┼────────────────────┐
               ▼                   ▼                     ▼
      ┌────────────────┐  ┌─────────────────┐  ┌─────────────────┐
      │  Retry Queue   │  │       DLQ       │  │ Visibility TMO  │
      │ (exp. backoff) │  │  (failed jobs)  │  │ (crash detect)  │
      └────────────────┘  └─────────────────┘  └────────┬────────┘
                                                          │
                                                          ▼
                                               ┌─────────────────┐
                                               │ Reaper Service  │
                                               │ (crash recovery)│
                                               └────────┬────────┘
                                                        │
                                               ┌────────▼────────┐
                                               │ Job re-enqueued │
                                               └─────────────────┘

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Observability Layer  →  Prometheus ──▶ Grafana Dashboard
                          Structured JSON Logs + Correlation IDs
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

---

## 🚀 Quick Start

> The fastest way to run GoQueue is via Docker Compose — spins up the API, workers, PostgreSQL, Redis, Prometheus, and Grafana in **one command**.

### Prerequisites

- 🐳 Docker & Docker Compose
- `make` (optional but recommended)

### 1. Clone the repository

```bash
git clone https://github.com/bhusaremayur/goqueue.git
cd goqueue
```

### 2. Start all services

```bash
# With Docker Compose
docker-compose -f deployments/docker-compose.yml up -d

# Or with Make
make up
```

### 3. Verify services are running

```bash
docker-compose ps
```

<div align="center">

| Service | URL | Description |
|:---|:---|:---|
| 🟢 **Go API Server** | [`http://localhost:8080`](http://localhost:8080) | Job submission & status |
| 📊 **Grafana Dashboard** | [`http://localhost:3000`](http://localhost:3000) | Pre-built dashboards |
| 📈 **Prometheus** | [`http://localhost:9090`](http://localhost:9090) | Raw metrics |

</div>

---

## 📁 Repository Structure

```
goqueue/
│
├── cmd/
│   ├── api/                →  API server entry point
│   ├── worker/             →  Worker process entry point
│   └── example-app/        →  Example usage
│
├── internal/
│   ├── job/                →  Core job domain logic
│   ├── storage/
│   │   ├── postgres/       →  PostgreSQL repositories
│   │   │   └── migrations/ →  SQL schema migrations
│   │   └── redis/          →  Redis queue implementation
│   ├── scheduler/          →  Job scheduling logic
│   └── reaper/             →  Crash recovery service
│
├── pkg/
│   ├── logging/            →  Structured JSON logger
│   ├── metrics/            →  Prometheus metric definitions
│   ├── middleware/          →  HTTP middleware (correlation IDs, etc.)
│   └── retry/              →  Exponential backoff policies
│
└── deployments/
    ├── docker-compose.yml
    └── prometheus.yml
```

---

## 🔍 Observability

GoQueue treats observability as a **first-class concern** — not an afterthought.

### Structured JSON Logging

Every log line carries a `correlation_id` for full request traceability:

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

### Prometheus Metrics

Exposed at `/metrics` — all critical signals tracked:

```
goqueue_jobs_enqueued_total        →  Total jobs submitted
goqueue_jobs_processed_total       →  Successfully processed
goqueue_jobs_failed_total          →  Failed jobs (by queue)
goqueue_queue_depth                →  Current queue depth (by priority)
goqueue_worker_utilization         →  Worker pool saturation %
goqueue_retry_count                →  Retry attempts in flight
goqueue_dlq_size                   →  Dead Letter Queue depth
goqueue_processing_latency_seconds →  Job processing duration histogram
```

### Grafana Dashboards

Pre-built dashboards ship with the repo.

```bash
open http://localhost:3000   # after docker-compose up
```

---

## ⚡ Load Testing Results

Tested with [**k6**](https://k6.io) at high concurrency:

```
┌─────────────────────────────────────────────┐
│           Load Test Summary (k6)            │
├──────────────────────┬──────────────────────┤
│  Virtual Users       │  10,000              │
│  Total Requests      │  40,000+             │
├──────────────────────┼──────────────────────┤
│  Backpressure (429)  │  ✅ Validated        │
│  Crash Recovery      │  ✅ All jobs back    │
│  DLQ Routing         │  ✅ Clean isolation  │
└──────────────────────┴──────────────────────┘
```

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

## 🤝 Contributing

Contributions are welcome — bug reports, feature suggestions, and pull requests alike.

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/your-feature`
3. Commit your changes: `git commit -m 'feat: add your feature'`
4. Push to the branch: `git push origin feat/your-feature`
5. Open a Pull Request

Please read [**CONTRIBUTING.md**](./CONTRIBUTING.md) before submitting — it covers running the test suite locally and the PR process.

---

## 📄 License

Licensed under the terms in [**LICENSE**](./LICENSE).

---

<div align="center">

**GoQueue v1.0 — Built in Public, 12 posts later.**

*Built with ❤️ using Go · Redis · PostgreSQL · Prometheus · Grafana*

<br/>

⭐ **Star this repo if GoQueue helped you!** ⭐

</div>