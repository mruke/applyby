# Quickstart

This guide explains how to start ApplyBy locally and manually verify the main user workflows.

## Requirements

Install:

- Go
- Node.js and npm
- Docker Desktop or another Docker-compatible runtime
- Git
- PowerShell

## 1. Install Frontend Dependencies

From the repository root:

```powershell
cd web
npm install
cd ..
```

## 2. Start the Local Database

From the repository root:

```powershell
docker compose up -d
```

## 3. Start the Backend API

In a dedicated terminal from the repository root:

```powershell
$env:APPLYBY_DATABASE_URL = "postgres://applyby:applyby@localhost:5432/applyby?sslmode=disable"
$env:APPLYBY_HTTP_ADDR = ":8080"
go run ./cmd/applyby-api
```

The API server listens on:

```text
http://localhost:8080
```

## 4. Start the Frontend

In a second terminal from the repository root:

```powershell
cd web
$env:VITE_API_BASE_URL = "http://localhost:8080"
npm run dev
```

Open the Vite URL printed by the terminal, usually:

```text
http://localhost:5173
```

## 5. Try the Main Workflows

After the database, backend, and frontend are running:

1. Open the dashboard.
2. Use the dashboard summary cards to filter the application list.
3. Open the applications workbench.
4. Create an application.
5. Open the application detail page.
6. Edit application details.
7. Change the application status.
8. Add, edit, and remove a contact.
9. Add, edit, and remove document metadata.
10. Schedule, edit, complete, and remove a reminder.
11. Confirm activity history updates.
12. Remove the application.
13. Confirm the application disappears from the list.

## Stop the Application

Stop the frontend and backend terminals with `Ctrl+C`.

Stop the database from the repository root:

```powershell
docker compose down
```

## Notes

`npm run dev` runs the frontend from source. It does not require `web/dist`.

`web/dist` is generated only when a production frontend build is created.
