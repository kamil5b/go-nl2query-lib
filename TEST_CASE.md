# Test Cases for Natural Language Database Query Interface

---

## 1. Workspace & Connection Management Tests

### 1.1 Add Workspace Tests (US-001)

#### TC-WS-001: Valid Database URL Submission
- **Scenario:** User submits a valid PostgreSQL database URL
- **Given:** Valid DB URL: `postgres://user:pass@localhost:5432/testdb`
- **When:** `POST /workspace` is called with valid DB URL
- **Then:** 
  - Returns `200 OK` with `tenant_id` (hashed)
  - Status is `IN_PROGRESS` or `DONE`
  - `message` field contains appropriate status text
  - Async ingestion task is enqueued

#### TC-WS-002: Invalid Database URL Format
- **Scenario:** User submits malformed database URL
- **Given:** Invalid DB URL: `not-a-valid-url`
- **When:** `POST /workspace` is called with invalid format
- **Then:**
  - Returns `400 Bad Request`
  - Error message indicates "invalid database URL format"

#### TC-WS-003: Database URL with Invalid Credentials
- **Scenario:** User submits valid URL format but unreachable database
- **Given:** DB URL: `postgres://wronguser:wrongpass@localhost:5432/testdb`
- **When:** `POST /workspace` is called
- **And:** No existing metadata for this workspace
- **Then:**
  - Returns `200 OK` with warning message
  - Status is `ERROR`
  - Message contains "Connection failed"
  - No async ingestion task is enqueued

#### TC-WS-004: Existing Workspace with Cached Metadata
- **Scenario:** User provides DB URL that was previously added but now unreachable
- **Given:** Workspace exists with cached metadata
- **And:** Database is currently unreachable
- **When:** `POST /workspace` is called with same DB URL
- **Then:**
  - Returns `200 OK` with same `tenant_id`
  - Status is `DONE`
  - Warning message: "Using cached metadata, live DB unreachable"

#### TC-WS-005: Metadata Checksum Unchanged
- **Scenario:** User refreshes connection to database with no schema changes
- **Given:** Workspace exists with matching metadata checksum
- **When:** `POST /workspace` is called with same DB URL
- **Then:**
  - Returns `200 OK` with `tenant_id`
  - Status is `DONE`
  - No ingestion task is enqueued
  - Message indicates "No changes detected"

#### TC-WS-006: Metadata Checksum Changed
- **Scenario:** User refreshes connection and database schema has changed
- **Given:** Workspace exists but new metadata checksum differs
- **When:** `POST /workspace` is called with same DB URL
- **Then:**
  - Returns `200 OK` with `tenant_id`
  - Status is `IN_PROGRESS`
  - Async ingestion task is enqueued
  - Previous error status in Redis is cleared

#### TC-WS-007: Concurrent Workspace Addition
- **Scenario:** Multiple users simultaneously add same database
- **Given:** Two requests for same DB URL at same time
- **When:** Both `POST /workspace` calls are made concurrently
- **Then:**
  - Both receive same `tenant_id`
  - Only one ingestion task is enqueued
  - No duplicate processing occurs

---

### 1.2 List Workspaces Tests (US-002)

#### TC-WS-008: List Empty Workspaces
- **Scenario:** User has no database connections
- **When:** `GET /workspace` is called on fresh system
- **Then:**
  - Returns `200 OK`
  - `workspaces` array is empty

#### TC-WS-009: List Multiple Workspaces
- **Scenario:** User has multiple database connections with different statuses
- **Given:** 3 workspaces: 1 DONE, 1 IN_PROGRESS, 1 ERROR
- **When:** `GET /workspace` is called
- **Then:**
  - Returns `200 OK` with all 3 workspaces
  - Each workspace contains: `tenant_id`, `db_name`, `status`
  - Statuses are correctly reflected

#### TC-WS-010: Pagination of Workspaces
- **Scenario:** User has many database connections (>50)
- **Given:** 100 workspaces created
- **When:** `GET /workspace?limit=10&offset=0` is called
- **Then:**
  - Returns first 10 workspaces
  - Subsequent calls with `offset=10` return next batch

---

### 1.3 Get Status Tests (US-003)

#### TC-WS-011: Status In Progress
- **Scenario:** Check status of workspace during ingestion
- **Given:** Workspace with `IN_PROGRESS` status
- **When:** `GET /status/{workspace_id}` is called
- **Then:**
  - Returns `200 OK`
  - Body: `{"status": "IN_PROGRESS"}`

#### TC-WS-012: Status Done
- **Scenario:** Check status of completed workspace
- **Given:** Workspace ingestion is complete
- **When:** `GET /status/{workspace_id}` is called
- **Then:**
  - Returns `200 OK`
  - Body: `{"status": "DONE"}`

#### TC-WS-013: Status Error with Message
- **Scenario:** Check status of failed workspace
- **Given:** Workspace ingestion failed with error
- **When:** `GET /status/{workspace_id}` is called
- **Then:**
  - Returns `200 OK`
  - Body: `{"status": "ERROR", "message": "Connection timed out"}`

#### TC-WS-014: Non-existent Workspace Status
- **Scenario:** Check status of workspace that doesn't exist
- **Given:** Invalid `workspace_id`
- **When:** `GET /status/{invalid_id}` is called
- **Then:**
  - Returns `404 Not Found`
  - Error message: "Workspace not found"

---

## 2. Background Ingestion Service Tests

### 2.1 Metadata Extraction Tests

#### TC-ING-001: Extract Metadata from PostgreSQL
- **Scenario:** Ingest metadata from valid PostgreSQL database
- **Given:** Connected PostgreSQL with 3 tables (users, orders, products)
- **When:** Ingestion task processes metadata extraction
- **Then:**
  - All 3 tables are extracted with correct metadata
  - Columns, types, and constraints are captured
  - Indexes and relations are identified
  - Status in Redis set to `IN_PROGRESS`

#### TC-ING-002: Extract Metadata with Complex Relations
- **Scenario:** Database has foreign key relationships
- **Given:** Schema with users → orders (FK) → products (FK)
- **When:** Metadata extraction processes relations
- **Then:**
  - All relations are captured correctly
  - Relation types (one-to-many, many-to-many) are identified
  - Source and target tables/columns are recorded

#### TC-ING-003: Extract Metadata with Null Columns
- **Scenario:** Database contains nullable columns
- **Given:** Table with columns: `id (NOT NULL)`, `email (NULL)`, `name (NULL)`
- **When:** Metadata extraction processes columns
- **Then:**
  - Nullable property is correctly set
  - Defaults are recorded where applicable

#### TC-ING-004: Checksum Generation from Metadata
- **Scenario:** Generate checksum for metadata comparison
- **Given:** Extracted metadata with 5 tables
- **When:** ChecksumGenerator processes metadata
- **Then:**
  - Deterministic checksum is generated (BLAKE3)
  - Same metadata always produces same checksum
  - Different metadata produces different checksum

#### TC-ING-005: Checksum Consistency on Refresh
- **Scenario:** Re-extract unchanged schema and compare checksums
- **Given:** Database with no schema changes
- **When:** Metadata is extracted twice
- **Then:**
  - Both checksums match
  - No re-ingestion is triggered

#### TC-ING-006: Checksum Mismatch on Schema Change
- **Scenario:** Schema changes (new column added)
- **Given:** Original metadata with checksum
- **When:** New column is added to table and re-extracted
- **Then:**
  - New checksum differs from original
  - Ingestion task is re-triggered

---

### 2.2 Embedding Generation Tests

#### TC-ING-007: Single Text Embedding
- **Scenario:** Generate embedding for table metadata
- **Given:** Table description: "users table with id, name, email columns"
- **When:** Embedder.Embed() processes the text
- **Then:**
  - Returns float32 array (embedding vector)
  - Vector dimension matches configured size (e.g., 1536 for OpenAI)
  - Same text produces consistent embedding

#### TC-ING-008: Batch Embedding Generation
- **Scenario:** Generate embeddings for multiple metadata items
- **Given:** 10 table/column descriptions
- **When:** Embedder.EmbedBatch() processes all texts
- **Then:**
  - Returns array of 10 embedding vectors
  - All vectors have consistent dimensions
  - Processing completes faster than sequential embedding

#### TC-ING-009: Embedding Chunking Strategy
- **Scenario:** Handle metadata exceeding token limits
- **Given:** Very large table with 100+ columns
- **When:** Metadata chunking strategy is applied
- **Then:**
  - Metadata is split into manageable chunks
  - Each chunk produces separate embeddings
  - Chunks are tagged with source table/column info

#### TC-ING-010: Failed Embedding Service
- **Scenario:** Embedder service is unavailable
- **Given:** Embedder API returns error
- **When:** Ingestion task attempts embedding
- **Then:**
  - Ingestion fails gracefully
  - Status is set to `ERROR: Embedding service unavailable`
  - Error is stored in Redis for user visibility

---

### 2.3 Vector Storage Tests

#### TC-ING-011: Upsert Vectors to VectorStore
- **Scenario:** Store embeddings in vector database
- **Given:** 50 vectors from metadata chunks
- **When:** VectorStore.Upsert() is called with `tenant_id`
- **Then:**
  - All 50 vectors are stored
  - Vectors are tagged with `tenant_id` for filtering
  - Metadata (table name, column info) is preserved

#### TC-ING-012: Vector Storage Completion
- **Scenario:** Complete ingestion workflow
- **Given:** All embeddings successfully created
- **When:** VectorStore.Upsert() completes
- **Then:**
  - `status:{tenant_id}` in Redis is deleted (cleanup)
  - Workspace status is set to `DONE`
  - User can immediately query

#### TC-ING-013: Vector Storage Failure
- **Scenario:** Vector database becomes unavailable
- **Given:** Vector store throws connection error
- **When:** Ingestion task attempts upsert
- **Then:**
  - Status is set to `ERROR: Vector store unavailable`
  - Workspace status reflects error
  - Task can be retried

#### TC-ING-014: Verify Vector Existence
- **Scenario:** Check if vectors exist for tenant before query
- **Given:** Completed ingestion for `tenant_id`
- **When:** VectorStore.Exists() is called
- **Then:**
  - Returns `true`
  - Query service can proceed safely

#### TC-ING-015: Vector Existence Check - Not Found
- **Scenario:** Check for vectors that don't exist
- **Given:** No ingestion performed for `tenant_id`
- **When:** VectorStore.Exists() is called
- **Then:**
  - Returns `false`
  - Query service returns 404

---

## 3. Query Service Tests

### 3.1 Pre-Query Validation Tests (US-004)

#### TC-QRY-001: Ingestion In Progress Check
- **Scenario:** User attempts query while ingestion is ongoing
- **Given:** Workspace status is `IN_PROGRESS`
- **When:** `POST /query` is called for workspace
- **Then:**
  - Returns `409 Conflict`
  - Error message: "Ingestion in progress, please try again later"

#### TC-QRY-002: Vector Data Not Found
- **Scenario:** User attempts query for non-existent workspace
- **Given:** No vectors exist for `tenant_id`
- **When:** `POST /query` is called with invalid `tenant_id`
- **Then:**
  - Returns `404 Not Found`
  - Error message: "Workspace not found or not yet ingested"

#### TC-QRY-003: Valid Workspace Status Check
- **Scenario:** Workspace is ready for querying
- **Given:** Workspace status is `DONE` and vectors exist
- **When:** Query validation runs
- **Then:**
  - Pre-validation passes
  - Proceeds to RAG retrieval

---

### 3.2 RAG Retrieval Tests

#### TC-QRY-004: Embed User Prompt
- **Scenario:** Convert natural language prompt to vector
- **Given:** User prompt: "Show me top 10 users by purchase amount"
- **When:** Query service embeds the prompt
- **Then:**
  - Prompt embedding is generated
  - Vector dimension matches metadata embeddings
  - Embedding captures semantic meaning

#### TC-QRY-005: Search Vector Database
- **Scenario:** Retrieve relevant metadata using similarity search
- **Given:** Prompt embedding and stored vectors
- **When:** VectorStore.Search() is called with `tenant_id` filter, limit=5
- **Then:**
  - Top 5 most similar vectors are returned
  - Results contain table/column metadata relevant to "users" and "purchase"
  - Similarity scores indicate relevance

#### TC-QRY-006: Empty Search Results
- **Scenario:** Query has no matching vectors
- **Given:** Prompt about non-existent tables
- **When:** VectorStore.Search() is called
- **Then:**
  - Returns empty array or very low-relevance results
  - LLM receives warning context
  - SQL may be generated but with lower confidence

#### TC-QRY-007: Context Assembly for LLM
- **Scenario:** Combine user prompt with retrieved context
- **Given:** 5 retrieved vectors with metadata
- **When:** LLM context is assembled
- **Then:**
  - Context includes: table names, column names, types, relations
  - Original prompt is preserved
  - Total context fits within LLM token limits

---

### 3.3 SQL Generation Tests (US-004)

#### TC-QRY-008: Simple SELECT Query Generation
- **Scenario:** Generate SELECT query from natural language
- **Given:** Prompt: "Show me all users"
- **And:** Context: users table with id, name, email columns
- **When:** LLMService.GenerateSQL() is called
- **Then:**
  - Returns valid SELECT query: `SELECT * FROM users;`
  - Query is syntactically correct

#### TC-QRY-009: Complex SELECT with Joins
- **Scenario:** Generate query involving multiple tables
- **Given:** Prompt: "Show me users and their orders"
- **And:** Context: users table, orders table, FK relation
- **When:** LLMService.GenerateSQL() is called
- **Then:**
  - Returns query with JOIN: `SELECT u.*, o.* FROM users u JOIN orders o ON u.id = o.user_id;`
  - JOIN condition matches foreign key relation

#### TC-QRY-010: Aggregation Query Generation
- **Scenario:** Generate aggregation query
- **Given:** Prompt: "Total spending by user"
- **And:** Context: orders table with user_id and amount columns
- **When:** LLMService.GenerateSQL() is called
- **Then:**
  - Returns: `SELECT user_id, SUM(amount) FROM orders GROUP BY user_id;`

#### TC-QRY-011: Query with ORDER BY and LIMIT
- **Scenario:** Generate query with sorting and pagination
- **Given:** Prompt: "Top 5 users by spend"
- **When:** LLMService.GenerateSQL() is called
- **Then:**
  - Returns: `SELECT * FROM users ORDER BY spend DESC LIMIT 5;`

---

### 3.4 SQL Validation & Self-Healing Tests (US-006)

#### TC-QRY-012: Validate Correct SQL
- **Scenario:** Validate syntactically correct SQL
- **Given:** SQL: `SELECT * FROM users WHERE id = 1;`
- **When:** SQLValidator.IsSafe() is called
- **Then:**
  - Returns `(true, nil)`
  - No DML/DDL detected
  - Safe to proceed to execution

#### TC-QRY-013: Dry-Run SQL Validation
- **Scenario:** Test SQL against schema without execution
- **Given:** Valid SQL: `SELECT * FROM users WHERE id = 1;`
- **When:** DatabaseConnection.ExecuteDryRun() is called
- **Then:**
  - Returns `nil` error
  - Confirms SQL is valid for current schema

#### TC-QRY-014: Invalid SQL with Syntax Error
- **Scenario:** LLM generates syntactically invalid SQL
- **Given:** SQL: `SELECT * FORM users;` (typo: FORM instead of FROM)
- **When:** Dry-run is attempted
- **Then:**
  - Returns error: "Syntax error near FORM"
  - Self-healing loop is triggered

#### TC-QRY-015: Self-Healing Loop - Retry 1
- **Scenario:** First retry to correct invalid SQL
- **Given:** Initial SQL with syntax error and error message
- **When:** LLMService.CorrectSQL() is called with error feedback
- **Then:**
  - Returns corrected SQL: `SELECT * FROM users;`
  - Retry count increments to 1

#### TC-QRY-016: Self-Healing Loop - Successful Correction
- **Scenario:** SQL is corrected after one retry
- **Given:** Retry count = 1, corrected SQL validates
- **When:** Dry-run validation succeeds
- **Then:**
  - Proceeds to execution
  - Self-healing loop exits successfully

#### TC-QRY-017: Self-Healing Loop - Max Retries Exceeded
- **Scenario:** SQL cannot be corrected after 3 attempts
- **Given:** MAX_RETRIES = 3, still failing after 3 attempts
- **When:** Self-healing loop exits
- **Then:**
  - Returns response with broken SQL
  - Warning: "WARN: Failed to generate valid SQL after 3 retries"
  - Data is null
  - User sees the problematic SQL for debugging

#### TC-QRY-018: Invalid Query with Column Reference Error
- **Scenario:** SQL references non-existent column
- **Given:** SQL: `SELECT non_existent_col FROM users;`
- **When:** Dry-run is attempted
- **Then:**
  - Returns error: "Column 'non_existent_col' does not exist"
  - Self-healing loop provides this specific error to LLM

#### TC-QRY-019: Invalid Query with Table Reference Error
- **Scenario:** SQL references non-existent table
- **Given:** SQL: `SELECT * FROM non_existent_table;`
- **When:** Dry-run is attempted
- **Then:**
  - Returns error: "Relation 'non_existent_table' does not exist"
  - Self-healing loop triggered

---

### 3.5 DDL/DML Detection Tests (US-005)

#### TC-QRY-020: INSERT Detection
- **Scenario:** Detect INSERT statement in generated SQL
- **Given:** SQL: `INSERT INTO users (name) VALUES ('Alice');`
- **When:** SQLValidator.ContainsDDLDML() is called
- **Then:**
  - Returns `true`
  - Execution is blocked

#### TC-QRY-021: UPDATE Detection
- **Scenario:** Detect UPDATE statement
- **Given:** SQL: `UPDATE users SET email = 'new@example.com' WHERE id = 1;`
- **When:** SQLValidator.ContainsDDLDML() is called
- **Then:**
  - Returns `true`
  - Execution is blocked

#### TC-QRY-022: DELETE Detection
- **Scenario:** Detect DELETE statement
- **Given:** SQL: `DELETE FROM users WHERE id = 1;`
- **When:** SQLValidator.ContainsDDLDML() is called
- **Then:**
  - Returns `true`
  - Execution is blocked

#### TC-QRY-023: DROP Detection
- **Scenario:** Detect DROP statement
- **Given:** SQL: `DROP TABLE users;`
- **When:** SQLValidator.ContainsDDLDML() is called
- **Then:**
  - Returns `true`
  - Execution is blocked

#### TC-QRY-024: ALTER Detection
- **Scenario:** Detect ALTER statement
- **Given:** SQL: `ALTER TABLE users ADD COLUMN age INT;`
- **When:** SQLValidator.ContainsDDLDML() is called
- **Then:**
  - Returns `true`
  - Execution is blocked

#### TC-QRY-025: TRUNCATE Detection
- **Scenario:** Detect TRUNCATE statement
- **Given:** SQL: `TRUNCATE TABLE users;`
- **When:** SQLValidator.ContainsDDLDML() is called
- **Then:**
  - Returns `true`
  - Execution is blocked

#### TC-QRY-026: DDL/DML in Response with Warning
- **Scenario:** Return SQL but block execution with warning
- **Given:** SQL: `DELETE FROM users WHERE spend < 0;`
- **When:** `POST /query` completes
- **Then:**
  - Returns `200 OK`
  - Response includes: `{"sql": "DELETE FROM users...", "data": null, "warning": "WARN: Query containing DDL/DML, system not allowed to process further"}`
  - Execution is NOT performed

#### TC-QRY-027: Case-Insensitive DDL/DML Detection
- **Scenario:** Detect DDL/DML regardless of case
- **Given:** SQL: `delete from users;` (lowercase)
- **When:** SQLValidator.ContainsDDLDML() is called
- **Then:**
  - Returns `true` (case-insensitive match)

#### TC-QRY-028: DDL/DML in WHERE Clause String
- **Scenario:** Handle INSERT/DELETE appearing in string literals
- **Given:** SQL: `SELECT * FROM comments WHERE text LIKE '%DELETE%';`
- **When:** SQLValidator.ContainsDDLDML() is called
- **Then:**
  - Returns `false` (DELETE inside string literal)
  - Query proceeds to execution

---

### 3.6 SQL Execution Tests

#### TC-QRY-029: Successful SQL Execution
- **Scenario:** Execute valid SELECT query and return results
- **Given:** SQL: `SELECT * FROM users LIMIT 2;`
- **And:** Database contains users
- **When:** DatabaseConnection.Execute() is called
- **Then:**
  - Returns `[]map[string]interface{}` with 2 user records
  - Each record contains all columns as key-value pairs

#### TC-QRY-030: Empty Result Set
- **Scenario:** Query returns no results
- **Given:** SQL: `SELECT * FROM users WHERE id = 999999;`
- **And:** No user with id=999999
- **When:** DatabaseConnection.Execute() is called
- **Then:**
  - Returns empty array `[]`
  - No error occurs

#### TC-QRY-031: Query Execution Timeout
- **Scenario:** Query takes too long to execute
- **Given:** SQL with expensive join on large tables
- **When:** DatabaseConnection.Execute() is called with timeout
- **Then:**
  - Returns error: "Query execution timeout"
  - No partial results are returned

#### TC-QRY-032: Database Connection Lost During Execution
- **Scenario:** Database becomes unavailable mid-query
- **Given:** Connection established but DB crashes during execution
- **When:** DatabaseConnection.Execute() is called
- **Then:**
  - Returns error: "Database connection lost"
  - Triggers self-healing (if enabled)

---

### 3.7 Query Response Tests

#### TC-QRY-033: Successful Query Response Format
- **Scenario:** Return properly formatted query response
- **Given:** SQL and results from execution
- **When:** Query completes successfully
- **Then:**
  - Returns JSON: `{"sql": "...", "data": [...], "warning": null}`
  - All fields present

#### TC-QRY-034: Response with Null Warning
- **Scenario:** Query succeeds without warnings
- **When:** Query execution completes
- **Then:**
  - Response includes `"warning": null`

#### TC-QRY-035: Response with SQL-Only Flag
- **Scenario:** User requests SQL generation without execution
- **Given:** `return_sql_only: true` in request
- **When:** Query processing completes
- **Then:**
  - Returns SQL and empty data: `{"sql": "...", "data": [], "warning": null}`
  - Database is NOT queried

#### TC-QRY-036: Large Result Set Handling
- **Scenario:** Query returns very large dataset
- **Given:** SQL returns 100,000 rows
- **When:** DatabaseConnection.Execute() returns results
- **Then:**
  - All rows are serialized to JSON
  - Response completes without memory issues
  - Pagination may be needed for UX

#### TC-QRY-037: NULL Values in Result Set
- **Scenario:** Result set contains NULL values
- **Given:** Query returns rows with NULL columns
- **When:** Results are serialized
- **Then:**
  - NULL values appear as `null` in JSON
  - No data loss or conversion errors

#### TC-QRY-038: Special Characters in Result Data
- **Scenario:** Result data contains special characters
- **Given:** Query returns text with quotes, newlines, unicode
- **When:** Results are serialized to JSON
- **Then:**
  - All characters are properly escaped
  - Data remains intact and readable

---

## 4. Database Unreachable Scenarios (US-007)

#### TC-UNR-001: Cached Metadata Available
- **Scenario:** Live database is down but vectors exist
- **Given:** Vectors exist in Vector DB, database is unreachable
- **When:** `POST /query` is called
- **Then:**
  - Query generation proceeds using cached metadata
  - SQL is generated successfully
  - Response includes warning: "WARN: Database unreachable, using cached schema"
  - SQL is returned but data is null

#### TC-UNR-002: Database Timeout During Query
- **Scenario:** Query hits timeout while executing
- **Given:** Database is slow or hanging
- **When:** Execute() call exceeds timeout
- **Then:**
  - Error is caught and reported
  - SQL is returned with warning: "Database unreachable or query timeout"
  - Data is null

#### TC-UNR-003: Connection Pool Exhausted
- **Scenario:** All database connections are in use
- **Given:** Max connection pool size reached
- **When:** New query attempts to get connection
- **Then:**
  - Returns error: "Connection pool exhausted"
  - Query may fail or be queued

#### TC-UNR-004: Network Failure During Query
- **Scenario:** Network connectivity lost mid-query
- **Given:** Database is reachable initially but network fails
- **When:** DatabaseConnection.Execute() is in progress
- **Then:**
  - Returns error: "Network error"
  - Response includes warning about database unavailability

---

## 5. Concurrency & Race Condition Tests

#### TC-CON-001: Concurrent Queries to Same Workspace
- **Scenario:** Multiple users query same database concurrently
- **Given:** 5 concurrent query requests for same `tenant_id`
- **When:** All requests are processed simultaneously
- **Then:**
  - All 5 queries complete successfully
  - Results are independent and correct
  - No data corruption occurs

#### TC-CON-002: Concurrent Ingestion and Query
- **Scenario:** User tries to query while ingestion is happening
- **Given:** Ingestion in progress, new query request comes in
- **When:** Query request is received
- **Then:**
  - Returns `409 Conflict`
  - Message: "Ingestion in progress"
  - User is instructed to retry

#### TC-CON-003: Concurrent Schema Refresh and Query
- **Scenario:** Schema is being refreshed while user queries
- **Given:** New ingestion starts, but old vectors still exist
- **When:** Query request arrives during refresh
- **Then:**
  - Can use old vectors OR wait for new ones (configurable)
  - Response indicates schema state

---

## 6. Edge Cases & Error Handling

#### TC-EDG-001: Empty Database Metadata
- **Scenario:** Database has no tables
- **Given:** Empty PostgreSQL database
- **When:** Ingestion processes empty metadata
- **Then:**
  - Metadata is extracted (empty arrays)
  - Checksum is generated
  - Status is set to DONE
  - Queries fail with "No tables found in database"

#### TC-EDG-002: Single Character Table/Column Names
- **Scenario:** Database schema uses very short names
- **Given:** Table `u`, column `a`
- **When:** Metadata is extracted
- **Then:**
  - Short names are handled correctly
  - Embeddings are generated
  - Queries may have lower accuracy due to ambiguity

#### TC-EDG-003: Very Long Table/Column Names
- **Scenario:** Names exceed typical limits
- **Given:** Column name: 500 characters
- **When:** Metadata is extracted
- **Then:**
  - Names are truncated or handled per database rules
  - Embeddings are still generated

#### TC-EDG-004: Special Characters in Table Names
- **Scenario:** Table names contain quotes, special chars
- **Given:** Table: `"user-data"`, `order_2024-01`
- **When:** Metadata is extracted and used in queries
- **Then:**
  - Names are properly quoted in generated SQL
  - Queries execute without syntax errors

#### TC-EDG-005: Reserved SQL Keywords as Column Names
- **Scenario:** Column is named with reserved keyword
- **Given:** Column named `select`, `from`, `where`
- **When:** Query is generated
- **Then:**
  - Column names are quoted: `"select"`, `"from"`
  - Queries execute successfully

#### TC-EDG-006: Duplicate Column Names Across Tables
- **Scenario:** Multiple tables have same column name
- **Given:** `users.id`, `orders.id`, `products.id`
- **When:** Query is generated with ambiguity
- **Then:**
  - LLM should qualify columns: `users.id`, `orders.id`
  - OR query fails with error about ambiguity

#### TC-EDG-007: Binary/Bytea Data Types
- **Scenario:** Database contains binary data
- **Given:** PostgreSQL BYTEA column
- **When:** Query returns binary data
- **Then:**
  - Binary data is base64 encoded in JSON response
  - OR handled per database driver rules

#### TC-EDG-008: UUID Data Types
- **Scenario:** Database uses UUID columns
- **Given:** PostgreSQL UUID columns
- **When:** Query returns UUID values
- **Then:**
  - UUIDs are returned as strings
  - Embeddings include UUID semantics

#### TC-EDG-009: JSONB/JSON Data Types
- **Scenario:** PostgreSQL JSONB columns
- **Given:** Column with nested JSON structure
- **When:** Query returns JSON data
- **Then:**
  - JSON is preserved in results
  - Metadata indicates JSON type

#### TC-EDG-010: Array/List Data Types
- **Scenario:** PostgreSQL array columns
- **Given:** Column: `int[]` or `text[]`
- **When:** Query returns array data
- **Then:**
  - Arrays are serialized to JSON arrays
  - Type information is preserved

---

## 7. Performance & Scalability Tests

#### TC-PER-001: Large Metadata Ingestion
- **Scenario:** Ingest database with 1000+ tables
- **Given:** Large enterprise database
- **When:** Ingestion task processes all metadata
- **Then:**
  - Completes within acceptable timeframe (e.g., < 5 minutes)
  - No memory leaks
  - All embeddings are generated

#### TC-PER-002: Vector Search Performance
- **Scenario:** Search across large vector dataset
- **Given:** 10,000+ vectors ingested
- **When:** VectorStore.Search() is called
- **Then:**
  - Returns top-K results within < 1 second
  - Search uses indexes for performance

#### TC-PER-003: Embedding Generation Throughput
- **Scenario:** Generate many embeddings efficiently
- **Given:** 1000 metadata chunks
- **When:** EmbedBatch() is called
- **Then:**
  - Batch embedding is faster than sequential
  - Rate limits are respected (if external API)

#### TC-PER-004: LLM Call Latency
- **Scenario:** Measure LLM response time
- **When:** GenerateSQL() is called
- **Then:**
  - Completes within acceptable latency (e.g., < 10 seconds)
  - Network retries are handled

#### TC-PER-005: Database Execution Performance
- **Scenario:** Large query execution
- **Given:** Query returns 100,000 rows
- **When:** Execute() completes
- **Then:**
  - Results are fetched and serialized efficiently
  - Memory usage is reasonable

---

## 8. Integration Tests

#### TC-INT-001: End-to-End Workspace Creation
- **Scenario:** Complete workflow from DB connection to ready state
- **Given:** Valid database URL
- **When:** POST /workspace, wait for ingestion, verify DONE status
- **Then:**
  - Workspace is created
  - Metadata is ingested
  - Vectors are stored
  - Status becomes DONE

#### TC-INT-002: End-to-End Query Processing
- **Scenario:** Complete workflow from prompt to result
- **Given:** Ready workspace
- **When:** POST /query with natural language prompt
- **Then:**
  - SQL is generated
  - Query executes
  - Results are returned

#### TC-INT-003: Multi-Database Workspace Test
- **Scenario:** Add multiple databases to system
- **Given:** 3 different database connections
- **When:** All are ingested and queried
- **Then:**
  - Each workspace operates independently
  - No cross-contamination of metadata

#### TC-INT-004: Workspace Update Workflow
- **Scenario:** Update workspace after schema change
- **Given:** Workspace exists, schema is modified
- **When:** POST /workspace is called again
- **Then:**
  - New ingestion is triggered
  - Old vectors are replaced
  - Status transitions correctly

---

## 9. Security Tests

#### TC-SEC-001: SQL Injection Prevention
- **Scenario:** Attempt SQL injection through prompt
- **Given:** Prompt: "users'; DROP TABLE users; --"
- **When:** LLMService processes prompt
- **Then:**
  - LLM treats it as literal string, not code
  - Generated SQL is safe
  - DROP is detected and blocked

#### TC-SEC-002: Escape Sequence Handling
- **Scenario:** Prompt contains escape sequences
- **Given:** Prompt with backslashes, quotes
- **When:** Prompt is embedded and passed to LLM
- **Then:**
  - Escaping is handled correctly
  - No injection vectors exist

#### TC-SEC-003: Tenant ID Isolation
- **Scenario:** User A queries workspace of User B
- **Given:** User A has User B's tenant_id
- **When:** User A attempts to query User B's workspace
- **Then:**
  - System allows it (metadata is public) OR requires auth
  - Behavior is documented

#### TC-SEC-004: Vector Metadata Encryption
- **Scenario:** Sensitive schema info in vector metadata
- **Given:** Sensitive table/column names
- **When:** Vectors are stored
- **Then:**
  - Metadata is encrypted OR access is controlled
  - Security policy is followed

---

## 10. Configuration & Customization Tests

#### TC-CFG-001: Configurable Embedder Implementation
- **Scenario:** Switch between different embedder implementations
- **Given:** Configuration selects OpenAI embedder
- **When:** System initializes
- **Then:**
  - OpenAI embedder is used for all embedding calls

#### TC-CFG-002: Configurable VectorStore Implementation
- **Scenario:** Switch between sqlite-vec, pgvector, qdrant
- **Given:** Configuration selects pgvector
- **When:** Vectors are upserted
- **Then:**
  - pgvector implementation is used

#### TC-CFG-003: Configurable LLM Implementation
- **Scenario:** Switch between OpenAI and Anthropic
- **Given:** Configuration selects Anthropic
- **When:** SQL generation is requested
- **Then:**
  - Anthropic LLM is called

#### TC-CFG-004: MAX_RETRIES Configuration
- **Scenario:** Configure self-healing retry limit
- **Given:** MAX_RETRIES = 5
- **When:** SQL fails validation
- **Then:**
  - System retries up to 5 times before giving up

#### TC-CFG-005: Timeout Configuration
- **Scenario:** Configure query execution timeout
- **Given:** QUERY_TIMEOUT = 30 seconds
- **When:** Query execution exceeds timeout
- **Then:**
  - Query is cancelled after 30 seconds