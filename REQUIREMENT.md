# 📄 PRD: Natural Language Database Query Interface

| Property | Value |
| --- | --- |
| **Status** | 🟢 Draft / 🟡 In Review / 🔴 Blocked |
| **Owner** | Product Engineering Team |
| **Tech Stack** | Go, Asynq, Redis, Vector DB, LLM (OpenAI/Anthropic) |
| **Type** | Backend System & API |
| **Last Updated** | February 7, 2026 |

---

## 1. 🎯 Overview

This project aims to build a backend system that allows users to connect their existing SQL databases via a URL and query them using natural language. The system handles the complex task of ingesting database metadata, generating embeddings, and using Large Language Models (LLMs) to translate natural language into SQL.

### Key Value Proposition

* **Seamless Onboarding:** Simple Database URL connection.
* **Intelligent Metadata Management:** Automatic change detection (checksums) and asynchronous ingestion.
* **Safe Execution:** Blocks DDL/DML operations to prevent data loss.
* **Resilience:** Self-healing SQL generation loops and graceful degradation if the DB is unreachable.

---

## 2. 👤 User Stories & Functional Requirements

### 2.1 Workspace & Connection Management

| ID | User Story | Acceptance Criteria |
| --- | --- | --- |
| **US-001** | As a user, I want to add a database connection so that I can query it later. | • POST `/workspace` accepts a Database URL.<br>

<br>• System validates basic URL format.<br>

<br>• System initiates async ingestion/validation.<br>

<br>• Returns a unique `tenant_id` (hashed). |
| **US-002** | As a user, I want to see a list of my database connections. | • GET `/workspace` returns list of connections.<br>

<br>• Shows status (Success/In Progress). |
| **US-003** | As a user, I want to check the status of a specific connection. | • GET `/status/{workspace_id}` returns:<br>

<br>  - `IN_PROGRESS` (if ingesting)<br>

<br>  - `ERROR: {msg}` (if failed)<br>

<br>  - `DONE` (if ready). |

### 2.2 Querying & Data Retrieval

| ID | User Story | Acceptance Criteria |
| --- | --- | --- |
| **US-004** | As a user, I want to prompt the system with natural language to get data. | • POST `/query` accepts `workspace_id` and `prompt`.<br>

<br>• Returns SQL Query and Result Set (tabular). |
| **US-005** | As a user, I want to be blocked from modifying data via natural language. | • If prompt implies DML/DDL (INSERT, DROP, etc.), system returns the generated SQL but **does not execute it**.<br>

<br>• Returns a specific warning message. |
| **US-006** | As a user, I want to know if the system fails to generate valid SQL. | • System attempts to self-correct SQL errors.<br>

<br>• If it fails after `N` retries, return the broken SQL with a Warning. |
| **US-007** | As a user, I want to receive SQL even if the live DB is down. | • If metadata exists but live connection fails, return generated SQL with a warning (Data unretrievable). |

---

## 3. 🏗️ System Architecture & Tech Stack

### 3.1 Technical Constraints & Stack

* **Language:** Go (Golang).
* **Task Queue:** Asynq (Redis-backed).
* **Infrastructure:**
* **Redis:** For State/Status management and Task Queue.
* **Vector DB:** Configurable Interface (implementations: `sqlite-vec`, `pgvector`, `qdrant`).
* **Embedder:** Configurable Interface (implementations: `openai`, `huggingface`).
* **LLM Service:** Configurable Interface (implementations: `OpenAI-compatible`, `Anthropic`).


* **Architecture Pattern:** DDD (Domain-Driven Design), SFF (Single-Function File), Dependency Injection.
* **Encryption:** BLAKE3 (for ID hashing and Checksums).
* **Development Methodology:** TDD (Test-Driven Development).
* *Constraint:* Define Interface -> Generate Mock -> Write Unit Test -> Write Implementation.
* *Coverage:* >80% strictly enforced.



### 3.2 High-Level Diagram

---

## 4. 🔄 System Services & Logic Flow

### 4.1 Server Service (Workspace)

**Responsibility:** Handle HTTP requests, manage connection logic, and trigger ingestion.

* **Logic: Adding/Refreshing Connection**
1. **ID Generation:** Hash the DB URL using BLAKE3 -> `tenant_id`.
2. **Vector Check:** Check if `tenant_id` exists in Vector DB.
3. **Live Connection Probe:** Attempt to connect to the provided DB URL.
* *If Connection Fails & Data Exists:* Return `tenant_id` with `WARN: Using cached metadata, live DB unreachable`.
* *If Connection Fails & No Data:* Throw `ERROR`.


4. **Metadata Extraction:** Extract Tables, Columns, Relations, Constraints, Indexes, Comments.
5. **Checksum:** specific metadata -> BLAKE3 Hash.
6. **Comparison:**
* *If Not Found or Checksum Mismatch:* Trigger **Ingestion Service** (Async Job) and return `tenant_id`.
* *If Match:* Return `tenant_id` (No action needed).


7. **Cleanup:** If successful, delete any lingering error status in Redis.



### 4.2 Background Service (Ingestion)

**Responsibility:** Process heavy metadata operations asynchronously.

* **Logic:**
1. **Set Status:** Redis Key `status:{tenant_id}` -> `"IN_PROGRESS"`.
2. **Processing:**
* Chunk metadata (Column-based strategy).
* Encrypt/Hash for checksums.
* Generate Embeddings via Embedder Service.


3. **Storage:** Upsert vectors to Vector Database.
4. **Completion:**
* *Success:* Delete `status:{tenant_id}` from Redis.
* *Failure:* Set `status:{tenant_id}` -> `"ERROR: {msg}"`.





### 4.3 Query Service

**Responsibility:** Text-to-SQL generation and Execution.

* **Logic Flow:**
1. **State Check:**
* Check Redis `status:{tenant_id}`.
* If `"IN_PROGRESS"` -> Return **409 Conflict**.


2. **Vector Check:**
* Check Vector DB for `tenant_id`.
* If Missing -> Return **404 Not Found**.


3. **RAG (Retrieval):**
* Embed User Prompt.
* Search Vector DB (Filter: `tenant_id`) -> Get Context.


4. **LLM Generation:**
* Prompt = `Original Prompt + Context`.
* Call Processing LLM -> Get SQL.


5. **SQL Evaluation Loop (Self-Healing):**
* *Configurable Limit:* `MAX_RETRIES` (e.g., 3).
* Validate/Dry-run SQL.
* If Error -> Feed error back to LLM -> Regenerate.
* If `MAX_RETRIES` reached -> Return SQL with Warn.


6. **Execution & Safety:**
* **Safety Check:** Parse SQL for DDL/DML (INSERT, UPDATE, DELETE, ALTER).
* *If Found:* Return SQL + `WARN: DDL/DML detected, execution blocked`.


* **Execution:**
* If `return_sql_only` flag is true -> Return SQL.
* Else -> Execute on Client DB via **Data Processing Service**.
* *Runtime Error Handling:* If execution fails, feed error back to LLM (loop step 5).







### 4.4 Data Processing Service

**Responsibility:** Execute raw SQL and format results.

* Input: SQL Query, DB Connection.
* Output: `[]map[string]any` (Tabular JSON).

---

## 5. 🔌 API Specification

### 5.1 Workspace

#### **Add Workspace**

* **Endpoint:** `POST /workspace`
* **Payload:**
```json
{ "db_url": "postgres://user:pass@host:5432/db" }

```


* **Response (200 OK):**
```json
{
  "tenant_id": "tenant_hash_123",
  "status": "IN_PROGRESS", // or DONE
  "message": "Ingestion started"
}

```


* **Response (Error):**
* `WARN`: "Will use existing stored data because... {error}"



#### **List Workspaces**

* **Endpoint:** `GET /workspace`
* **Response:**
```json
{
  "workspaces": [
    { "tenant_id": "...", "db_name": "...", "status": "DONE" }
  ]
}

```



#### **Get Status**

* **Endpoint:** `GET /status/{workspace_id}`
* **Response:**
* `"IN_PROGRESS"`
* `"ERROR: Connection timed out"`
* `"DONE"`



### 5.2 Query

#### **Prompt Database**

* **Endpoint:** `POST /query`
* **Payload:**
```json
{
  "workspace_id": "tenant_hash_123",
  "prompt": "Show me top 5 users by spend",
  "return_sql_only": false
}

```


* **Response (Success):**
```json
{
  "sql": "SELECT * FROM users ORDER BY spend DESC LIMIT 5;",
  "data": [
    { "id": 1, "name": "Alice", "spend": 5000 },
    { "id": 2, "name": "Bob", "spend": 4500 }
  ],
  "warning": null
}

```


* **Response (Safety Warning):**
```json
{
  "sql": "DELETE FROM users WHERE spend < 0;",
  "data": null,
  "warning": "WARN: Query containing DDL/DML, system not allowed to process further"
}

```


* **Error Codes:**
* `404`: Workspace not found / not ingested.
* `409`: Ingestion in progress.



---

## 6. 🛡️ Non-Functional Requirements & Security

### 6.1 Security

* **Encryption:** Use **BLAKE3** for all internal hashing (Tenant IDs, Checksums).
* **Credential Handling:** DB URLs must be encrypted at rest if stored persistently (outside of the vector metadata context).
* **Injection Prevention:** The system relies on the LLM to generate valid SQL, but the **Data Processing Service** must use a read-only transaction mode where possible or strictly regex/parser check for destructive keywords (DROP, TRUNCATE, etc.) before execution.

### 6.2 Testing (TDD)

> **Strict Policy:** No code is written without a failing test.

1. **Interface Definition:** Define `Service` and `Repository` interfaces in domain layer.
2. **Mock Generation:** Use `gomock` to generate mocks.
3. **Unit Tests:** Write tests ensuring >80% coverage.
4. **Implementation:** Write the SFF code to pass the tests.
