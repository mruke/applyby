# Quickstart

ApplyBy currently has a runnable local full-stack workflow:

- Go backend API server
- PostgreSQL development database
- React frontend UI
- TypeScript frontend implementation
- Layered backend and frontend tests
## PostgreSQL Development Database

ApplyBy uses Docker Compose to run PostgreSQL locally for development and integration tests.

Start PostgreSQL:

```powershell
docker compose up -d
```

Run all tests without PostgreSQL integration tests:

```powershell
go test ./...
```

Run PostgreSQL integration tests:

```powershell
$env:APPLYBY_INTEGRATION_TESTS = "1"
$env:APPLYBY_DATABASE_URL = "postgres://applyby:applyby@localhost:5432/applyby?sslmode=disable"
go test ./internal/storage/postgres
```

Stop PostgreSQL:

```powershell
docker compose down
```

## Frontend Setup

Install frontend dependencies:

```powershell
cd web
npm install
```

Run frontend tests:

```powershell
npm test
```

Run the frontend development server without the backend API:

```powershell
npm run dev
```

Build the frontend:

```powershell
npm run build
```

## Backend API Server

Start PostgreSQL:

```powershell
docker compose up -d
```

Run the Go API server:

```powershell
$env:APPLYBY_DATABASE_URL = "postgres://applyby:applyby@localhost:5432/applyby?sslmode=disable"
$env:APPLYBY_HTTP_ADDR = ":8080"
go run ./cmd/applyby-api
```

The API server listens on:

```text
http://localhost:8080
```

For local frontend development, start the frontend in a separate terminal:

```powershell
cd web
$env:VITE_API_BASE_URL = "http://localhost:8080"
npm run dev
```

## Validation

## Current Workflow Coverage Note

The validation flow covers the current supported workflows: create application, search/list applications, open detail, update status, schedule/complete reminders, add contacts, add document metadata, and verify activity history.

Full maintenance workflows are still planned for editing/removing contacts, editing/removing document metadata, editing/canceling reminders, and editing non-status application details.


Run backend formatting and tests from the repository root:

```
powershell
gofmt -w cmd internal
go test ./...
git diff --check
```

Run frontend tests and production build:

```
powershell
cd web
npm test
npm run build
cd ..
```

The current UI supports application creation, application search, application detail views, status updates, reminders, contacts, document metadata, activity history, and a basic dashboard overview.
