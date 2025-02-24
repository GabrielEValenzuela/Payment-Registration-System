# Source Code (`src`)

This directory contains the core implementation of the application, including business logic, API handlers, storage, configuration, and documentation.

## Directory Structure

```
src/
│── benchmark/                        # Performance benchmarking utilities
│── cmd/                               # Main application commands
│   ├── handlers/                      # API route handlers
│   ├── server/                        # Server initialization and routing
│── docs/                              # API documentation and specifications
│── internal/                          # Private application logic
│   ├── config/                        # Configuration management
│   ├── models/                        # Domain models and DTOs
│   ├── services/                      # Business logic services
│   ├── storage/                       # Data access layer
│   │   ├── entities/                   # ORM entities and schema definitions
│   │   ├── non_relational/             # NoSQL storage implementation
│   │   │   ├── repository/             # MongoDB repository layer
│   │   ├── relational/                 # SQL storage implementation
│   │   │   ├── repository/             # SQL repository layer
│   ├── testutils/                      # Test utilities and mocks
│── pkg/                                # Shared utilities and libraries
│   ├── logger/                         # Structured logging package
```

## Overview

The `src` folder follows a **clean architecture** approach by organizing the code into separate layers:

1. **Command (`cmd/`)**  
   - Entry points for running the application, including API handlers and server initialization.

2. **Documentation (`docs/`)**  
   - API specifications, OpenAPI documentation, and other reference materials.

3. **Internal (`internal/`)**  
   - Private business logic and application components:
     - **Configuration (`config/`)**: Manages environment variables and settings.
     - **Models (`models/`)**: Defines domain entities and DTOs.
     - **Services (`services/`)**: Implements core business logic.
     - **Storage (`storage/`)**: Handles data persistence with relational (SQL) and non-relational (NoSQL) databases.
     - **Test Utilities (`testutils/`)**: Helpers and mock data for testing.

4. **Package (`pkg/`)**  
   - Reusable utilities such as logging, error handling, and middleware.

## Storage Implementation

The application supports both **relational** and **non-relational** databases:

- **Relational (SQL)**
  - Uses ORM-based entities (`storage/relational/entities/`).
  - Implements data persistence through repositories (`storage/relational/repository/`).

- **Non-Relational (MongoDB)**
  - Implements MongoDB storage with repositories (`storage/non_relational/repository/`).

## Logging

The `pkg/logger/` package provides a structured logging mechanism using a centralized logging strategy.

## Benchmarking

The `benchmark/` directory includes performance analysis tools for profiling key application functionalities.

## Testing

- Unit and integration tests use mocks and test utilities (`internal/testutils/`).
- Supports structured testing strategies to ensure reliability.

## Contribution Guidelines

- Follow the directory structure and maintain modularity.
- Use dependency injection for services and storage.
- Write unit tests for all business logic components.

---
