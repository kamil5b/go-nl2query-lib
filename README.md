# go-nl-sql

**go-nl-sql** is a high-performance, production-ready Go agent that transforms natural language prompts into executable SQL queries. It uses a **Server-Worker architecture** and **RAG (Retrieval-Augmented Generation)** to intelligently understand your database schema and provide data insights safely.

## 🚀 Key Features

* **Asynchronous Ingestion:** Automatic schema metadata extraction and vectorization using `Asynq` and `Redis`.
* **Self-Healing SQL:** A configurable retry loop that feeds database execution errors back to the LLM for automatic query correction.
* **Built-in Safety:** Strict regex-based blocking for DDL/DML (INSERT, UPDATE, DELETE, DROP) to ensure read-only operations.
* **Schema Checksums:** Uses **BLAKE3** hashing to detect schema changes; only re-ingests metadata when the database structure actually changes.
* **Plug-and-Play Infra:** Swappable interfaces for Vector DBs (Qdrant/Pgvector), Embedders (OpenAI/HuggingFace), and LLMs (OpenAI/Anthropic).

---

## 🏗️ Architecture

The project follows **Domain-Driven Design (DDD)** and **Single-Function File (SFF)** principles to ensure maximum maintainability and testability.

### Tech Stack

* **Language:** Go 1.21+
* **Task Queue:** Asynq (Redis-backed)
* **Security:** BLAKE3 Hashing & Encryption
* **Patterns:** Dependency Injection, Strategy Pattern (for swappable providers)

---

## 📂 Project Structure

```text
├── cmd/api/                # Entry point
├── internal/
│   ├── domain/             # Core logic & interfaces (Business Rules)
│   ├── app/                # Use cases (Single-function files)
│   ├── infrastructure/     # SQL, VectorDB, and LLM implementations
│   └── worker/             # Background ingestion handlers
├── mocks/                  # Mockery generated files for TDD
└── docker-compose.yaml     # Local development environment

```

---

## 🛠️ Development & TDD

We maintain a strict **TDD (Test-Driven Development)** workflow. No implementation code is written before its corresponding interface and unit test.

1. **Generate Mocks:** `make mocks`
2. **Run Tests:** `go test -v -cover ./...`
3. **Coverage:** Strictly enforced **>80%** coverage.

---

## 🚦 API Reference

### 1. Workspace Management

* `POST /workspace`: Register a new DB connection and trigger async ingestion.
* `GET /workspace`: List all registered database connections.
* `GET /status/{workspace_id}`: Check if ingestion is `IN_PROGRESS`, `DONE`, or `ERROR`.

### 2. Querying

* `POST /query`: Submit a natural language prompt.
* **Returns:** Generated SQL and the resulting tabular data.
* **Safety:** Blocks destructive queries while still returning the generated SQL for review.



---

## 📦 Setup

1. **Clone and Install:**
```bash
git clone https://github.com/your-repo/go-nl-sql.git
cd go-nl-sql
go mod download

```


2. **Start Infrastructure:**
```bash
docker-compose up -d

```


3. **Run API:**
```bash
go run cmd/api/main.go

```
