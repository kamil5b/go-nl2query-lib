# Go-NL2Query for Database Queries

A Go library for converting natural language queries into database queries (SQL or NoSQL) using embeddings and LLM-powered query generation. Query any database type using plain English.

## Overview

`Go-NL2Query` is a modular, domain-driven library that enables applications to transform user-friendly natural language questions into executable database queries. It supports both SQL and NoSQL databases, leveraging vector embeddings for semantic search and large language models to intelligently generate queries based on your database schema and structure.

## Key Features

- **Natural Language to Database Query Translation**: Convert natural language into SQL or NoSQL queries
- **Multi-Database Support**: Works with SQL databases (PostgreSQL, MySQL, etc.) and NoSQL databases (MongoDB, DynamoDB, etc.)
- **Vector Embeddings**: Semantic search capabilities for context-aware query generation
- **LLM-Powered Generation**: Intelligent query synthesis using language models
- **Modular Architecture**: Clean separation of concerns using hexagonal architecture (ports & adapters)
- **Multi-Tenant Support**: Built-in tenant isolation across all components
- **Async Task Processing**: Asynchronous job processing using task queues
- **Type-Safe Domain Models**: Well-defined domain objects for vectors, queries, and workspaces
- **Extensible Ports**: Interface-driven design for easy integration with different providers
- **Flexible Query Validation**: Validate and execute queries before returning results

## Project Structure

```
go-nl-sql/
├── domains/           # Core business logic and domain models
├── ports/            # Interface definitions for external dependencies
├── services/         # Business logic and use case implementations
│   ├── ingestion/   # Data ingestion and vectorization
│   ├── query/       # Query generation and validation
│   └── workspace/   # Workspace management
├── adapters/        # Concrete implementations of ports
├── testsuites/      # Shared test utilities and mocks
└── Makefile         # Development automation
```

### Domains

Core domain models:
- **Vector**: Represents embeddings with metadata and content
- **Query**: Natural language query with result tracking
- **Metadata**: Contextual information storage
- **Error**: Domain-specific error handling
- **Workspace**: Multi-tenant workspace abstraction

### Ports (Interfaces)

External service abstractions:
- **EmbedderPort**: Text embedding generation
- **VectorStorePort**: Vector storage and retrieval
- **LLMPort**: Query generation using language models
- **ClientDatabasePort**: Source database access (SQL/NoSQL)
- **InternalDatabasePort**: Query result persistence
- **EncryptPort**: Data encryption/decryption
- **HashPort**: Cryptographic hashing
- **QueryValidatorPort**: Query validation before execution
- **TaskQueuePort**: Async job management
- **StatusPort**: Operation status tracking
- **WorkspacePort**: Workspace management

### Services

Business logic implementations:
- **IngestionService**: Handles data ingestion, vectorization, and storage
- **VectorizeAndStoreService**: Processes and stores vectors
- **QueryService**: Natural language to database query conversion

## Development

### Installation

```bash
# Initialize workspace with all modules
make work-init

# Sync dependencies across modules
make work-sync
```

### Testing

```bash
# Run all unit tests
make test

# Generate coverage report
make coverage
```

### Mock Generation

```bash
# Install mockgen tool
make mocks-install

# Generate all mocks from ports
make mocks
```

### Dependency Management

```bash
# Tidy all go.mod files
make work-tidy

# Vendor dependencies
make work-vendor
```

### Available Commands

Run `make help` to see all available commands.

## Architecture Principles

This project follows **Hexagonal Architecture** (Ports & Adapters) principles:

- **Domains**: Pure business logic independent of external concerns
- **Ports**: Interfaces defining contracts for external services
- **Adapters**: Concrete implementations of ports for specific technologies
- **Services**: Orchestrate domain logic and port usage
- **Testsuites**: Shared testing utilities and mock implementations

This design enables:
- Easy testing through mock implementations
- Simple swapping of external service providers
- Simple addition of new database support
- Clear separation of concerns
- Framework independence

## Usage Example

```go
// Create service with injected dependencies
ingestionService := ingestion.NewIngestionService(
    config,
    embedderAdapter,
    vectorStoreAdapter,
    statusAdapter,
)

// Process natural language queries through the system
// Implementation details depend on your adapter implementations
```

## Query Flow

1. **Ingestion**: Documents are vectorized and stored with metadata
2. **Natural Language Input**: User provides a natural language query
3. **Vector Search**: Query is embedded and semantically similar documents are retrieved
4. **Context Building**: Retrieved documents provide context for LLM
5. **Query Generation**: LLM generates appropriate SQL or NoSQL query
6. **Validation**: Generated query is validated before execution
7. **Execution**: Query is executed against the target database
8. **Result Formatting**: Results are formatted and returned to user

## Multi-Tenancy

All services are designed with multi-tenant support. Each operation accepts a `tenantID` parameter to ensure data isolation and security across different organizations or users.

## Go Versions

Requires Go 1.25.6 or later.

## Contributing

When adding new features:
1. Define domain models in `domains/`
2. Create port interfaces in `ports/`
3. Implement business logic in `services/`
4. Implement concrete adapters in `adapters/` for specific database types
5. Add tests alongside implementations
6. Run `make coverage` to ensure test coverage
7. Generate mocks with `make mocks` for testing

## License

MIT License - See [LICENSE](./LICENSE)

This project is open-source and free to use, modify, and distribute. For details, see the LICENSE file in this repository.
