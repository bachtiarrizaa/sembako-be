# AGENTS.md - Sembako Backend

This document outlines the architectural guidelines, coding standards, and development rules for AI agents and developers working on the `sembako-be` project.

## 1. Project Architecture & Folder Responsibilities

The project follows a **Clean / Layered Architecture** pattern in Go:

- **`cmd/`**: Application entry points.
  - `cmd/api/main.go`: Initializes configuration and bootstraps the HTTP server.
  - `cmd/seeder/main.go`: Database seeder utility.
- **`internal/bootstrap/`**: Dependency injection wiring. Connects database, repositories, usecases, controllers, and routers.
- **`internal/config/`**: Environment configuration (`config.go`) and database initialization (`database.go`).
- **`internal/delivery/http/`**:
  - `controller/`: HTTP handlers. Binds incoming requests, performs input validation, invokes usecases, and formats responses.
  - `router/`: API routing setup using Gin (`/api/*`).
- **`internal/entity/`**: Domain database models mapped to database tables via GORM tags.
- **`internal/model/`**: Data Transfer Objects (DTOs), request/response structs, and pagination structs.
- **`internal/pkg/`**: Shared packages and utilities.
  - `errs/`: Custom application errors (`NotFound`, `Conflict`, `Internal`, etc.).
  - `utils/`: Standardized JSON response helpers and pagination formatters.
- **`internal/repository/`**: Data access layer. Executes database queries using GORM.
- **`internal/usecase/`**: Business logic layer. Enforces business rules, validation checks, and orchestrates repositories.
- **`migrations/`**: SQL migration files (`.up.sql` and `.down.sql`) for database schema management.

### Project Folder Tree
```text
sembako-be/
├── cmd/
│   ├── api/
│   │   └── main.go           # Application entry point & HTTP server bootstrapper
│   └── seeder/
│       └── main.go           # Database seeder utility
├── internal/
│   ├── bootstrap/
│   │   └── bootstrap.go      # Dependency injection wiring
│   ├── config/
│   │   ├── config.go         # Environment variables configuration
│   │   └── database.go       # GORM database connection setup
│   ├── delivery/
│   │   └── http/
│   │       ├── controller/   # HTTP handlers, payload binding, and validation
│   │       └── router/       # Gin API routing setup (/api/*)
│   ├── entity/               # GORM database models / tables
│   ├── model/                # DTOs, request/response structs, pagination structs
│   ├── pkg/
│   │   ├── errs/             # Custom application error definitions
│   │   └── utils/            # Standard response helpers & pagination formatters
│   ├── repository/           # Data access layer (GORM queries)
│   └── usecase/              # Business logic & domain rules
├── migrations/               # SQL migration files (.up.sql & .down.sql)
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── AGENTS.md
```

---

## 2. Strict Coding Standards & Rules

### No `any` Type
- **The use of `any` (`interface{}`) is strictly prohibited in Go code.**
- Always use explicit types, domain structs, DTOs, or strongly-typed interfaces to ensure type safety and maintainability.

### Error Handling
- Never return raw database errors directly to clients.
- Use custom error helpers from `internal/pkg/errs` (e.g., `errs.NewNotFound`, `errs.NewConflict`, `errs.NewInternal`).
- Controllers must catch errors and format them uniformly using helper error responses.

### Request Validation
- All incoming HTTP request payloads and query parameters must be bound and validated using `github.com/go-playground/validator/v10`.
- Validation errors should return HTTP status `422 Unprocessable Entity`.

### API Response Format
- Success responses must use `utils.SuccessResponse` or `utils.SuccessResponseWithPagination`.
- Error responses must use `utils.ErrorResponse`.

---

## 3. Dependency Injection & Routing

- Dependency injection must be performed bottom-up inside `internal/bootstrap/bootstrap.go`:
  1. Initialize Database connection (`config.NewDatabase`).
  2. Initialize Repositories (`repository.New...`).
  3. Initialize Usecases (`usecase.New...`).
  4. Initialize Controllers (`controller.New...`).
  5. Setup Routes (`router.Setup`).

---

## 4. Development Workflow & Verification

- **Code Quality:** Ensure all code adheres to Go formatting standards (`go fmt`, `go vet`).
- **Git Operations:** AI agents are strictly prohibited from performing git operations (`git add`, `git commit`, `git push`). All version control actions are reserved exclusively for the human developer.
- **Migrations:** Create new database changes using migration files in `migrations/`.
