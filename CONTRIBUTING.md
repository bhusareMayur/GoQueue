<div align="center">

# 🚀 Contributing to GoQueue

**Thank you for your interest in making GoQueue better!**
We welcome all contributions — from bug reports and docs to full feature implementations.

[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat-square&logo=go)](https://golang.org/dl/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen?style=flat-square)](https://github.com/bhusaremayur/goqueue/pulls)
[![Conventional Commits](https://img.shields.io/badge/Commits-Conventional-FE5196?style=flat-square)](https://www.conventionalcommits.org/)

</div>

---

## 📌 Table of Contents

- [🐛 Reporting Bugs](#-reporting-bugs)
- [💡 Suggesting Enhancements](#-suggesting-enhancements)
- [💻 Local Development Setup](#-local-development-setup)
- [📁 Project Structure](#-project-structure)
- [🛠 Development Workflow](#-development-workflow)
- [📝 Commit Message Guidelines](#-commit-message-guidelines)
- [🚀 Pull Request Process](#-pull-request-process)

---

## 🐛 Reporting Bugs

Found something broken? Please [open an issue](https://github.com/bhusaremayur/goqueue/issues) with the following details:

| Field | What to Include |
|---|---|
| **Title** | A clear and descriptive summary of the bug |
| **Steps to Reproduce** | Exact steps to trigger the issue |
| **Expected vs Actual** | What you expected to happen vs what did happen |
| **System Context** | OS, Go version, Docker version |
| **Logs** | Relevant output from the API, Worker, or Redis containers |

> 💡 The more detail you provide, the faster we can triage and fix it.

---

## 💡 Suggesting Enhancements

Have an idea to make GoQueue faster, more reliable, or easier to use?

1. Open an issue with the label `enhancement` or `feature`
2. Clearly describe the **use case** and the **problem it solves**
3. Optionally outline a potential approach or architecture

We love well-scoped proposals with concrete motivations!

---

## 💻 Local Development Setup

### Prerequisites

Before you begin, make sure you have the following installed:

- [**Go**](https://golang.org/dl/) — version **1.20 or higher**
- [**Docker**](https://docs.docker.com/get-docker/) — for running infrastructure
- [**Docker Compose**](https://docs.docker.com/compose/install/) — for orchestrating services
- [**Make**](https://www.gnu.org/software/make/) *(optional but recommended)*

### Step-by-Step Installation

**1. Fork the repository** on GitHub, then clone your fork:

```bash
git clone https://github.com/YOUR-USERNAME/goqueue.git
cd goqueue
```

**2. Add the upstream remote** so you can keep in sync with the main repo:

```bash
git remote add upstream https://github.com/bhusaremayur/goqueue.git
```

**3. Boot up the infrastructure** (PostgreSQL, Redis, Prometheus, Grafana):

```bash
docker-compose -f deployments/docker-compose.yml up -d
```

**4. Run the services locally** for development and debugging:

```bash
# Start the HTTP API server
go run cmd/api/main.go

# Start the distributed background worker
go run cmd/worker/main.go
```

---

## 📁 Project Structure

Here's a quick map of the codebase to help you find your way around:

```
goqueue/
├── cmd/
│   ├── api/            → Entry point for the HTTP API server
│   └── worker/         → Entry point for the background worker process
│
├── internal/           → Core domain logic
│   ├── domain/         → Jobs, queues, schedulers
│   ├── storage/
│   │   └── postgres/
│   │       └── migrations/  → SQL migration files
│   └── reaper/         → Dead job reaper service
│
├── pkg/                → Reusable, public-facing packages
│   ├── metrics/        → Prometheus metrics
│   ├── middleware/      → HTTP middleware
│   ├── retry/          → Retry policies
│   └── logging/        → Structured logging
│
└── deployments/        → Infrastructure configs
    ├── docker-compose.yml
    └── prometheus.yml
```

---

## 🛠 Development Workflow

Follow these steps for every contribution:

**1. Sync your fork** with the latest upstream changes:

```bash
git checkout main
git pull upstream main
```

**2. Create a feature branch** for your work:

```bash
git checkout -b feature/your-feature-name
```

> Branch naming convention: `feature/`, `fix/`, `docs/`, `refactor/`

**3. Write your code** — follow standard Go idioms and keep things idiomatic.

**4. Format your code** before committing:

```bash
go fmt ./...
```

**5. Write and run tests** — all new features must include tests:

```bash
go test -v ./...
```

> Tests live alongside their source files, e.g. `internal/domain/job/service_test.go`

---

## 📝 Commit Message Guidelines

We follow the [**Conventional Commits**](https://www.conventionalcommits.org/) specification for a clean, readable git history and automated release notes.

### Format

```
<type>(<optional scope>): <short description>
```

### Types

| Type | When to Use |
|---|---|
| `feat` | ✨ A new feature |
| `fix` | 🐛 A bug fix |
| `docs` | 📚 Documentation-only changes |
| `style` | 💅 Formatting, missing semicolons — no logic change |
| `refactor` | ♻️ Code change that is neither a fix nor a feature |
| `test` | 🧪 Adding or correcting tests |
| `chore` | 🔧 Tooling, config, or dependency updates |

### Examples

```bash
feat: implement multi-queue priority dispatching
fix(scheduler): add backoff delay on redis error
docs: update README with architecture diagram
test(worker): add integration test for retry logic
chore: initialize base project structure
```

---

## 🚀 Pull Request Process

**1. Push your branch** to your fork:

```bash
git push origin feature/your-feature-name
```

**2. Open a Pull Request** against the `main` branch of [`bhusaremayur/goqueue`](https://github.com/bhusaremayur/goqueue).

**3. Describe your changes** clearly in the PR description:
- What problem does this solve?
- How did you test it?
- Any relevant context or screenshots?

**4. Review process** — maintainers will review your PR and may request changes or ask questions. Stay engaged!

**5. Merge** — once approved and all CI checks pass, your PR will be merged. 🎉

---

> **Thank you** for taking the time to contribute to GoQueue.
> Every improvement, no matter how small, makes a difference. Happy coding! 💙