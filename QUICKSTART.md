# Quickstart

ApplyBy currently has a runnable Go backend test workflow and a Docker Compose PostgreSQL development database.

The selected implementation direction is:

- Go for the backend service
- PostgreSQL for persistence
- React for the frontend UI
- TypeScript for frontend implementation
- Layered testing with unit, integration, end-to-end, and helper test areas

This document will be expanded as runnable implementation slices are added.

## Planned Sections

- Prerequisites
- Backend setup
- Database setup
- Frontend setup
- Development commands
- Test commands
- Benchmark commands
- Generated files

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

Run the frontend development server:

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
