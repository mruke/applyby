package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mruke/applyby/internal/api"
	"github.com/mruke/applyby/internal/application"
	"github.com/mruke/applyby/internal/config"
	"github.com/mruke/applyby/internal/storage/postgres"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// -----------------------------------------------------------------------------
// main
//
// Starts the ApplyBy HTTP API server.
// -----------------------------------------------------------------------------
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// run
//
// Wires configuration, database persistence, application services, API handlers,
// and the HTTP server lifecycle.
// -----------------------------------------------------------------------------
func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	databaseConfig := config.LoadDatabaseConfig()
	serverConfig := config.LoadServerConfig()

	db, err := openDatabase(ctx, databaseConfig.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := postgres.RunMigrations(ctx, db); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	repository := postgres.NewApplicationRepository(db)

	applicationHandlers := api.NewApplicationHandlers(
		application.NewCreateApplicationService(repository, repository),
		application.NewListApplicationsService(repository),
		application.NewGetApplicationService(repository),
		application.NewUpdateApplicationDetailsService(repository, repository),
		application.NewUpdateApplicationStatusService(repository, repository),
	)

	workflowHandlers := api.NewWorkflowHandlers(api.WorkflowHandlerDependencies{
		SearchApplications: application.NewSearchApplicationsService(repository),
		ListActivityEvents: application.NewListActivityEventsService(repository),
		ScheduleReminder:   application.NewScheduleReminderService(repository, repository),
		ListReminders:      application.NewListRemindersService(repository),
		CompleteReminder:   application.NewCompleteReminderService(repository, repository),
		AddContact:         application.NewAddContactService(repository, repository),
		ListContacts:       application.NewListContactsService(repository),
		UpdateContact:      application.NewUpdateContactService(repository, repository),
		RemoveContact:      application.NewRemoveContactService(repository, repository),
		AddDocument:        application.NewAddDocumentService(repository, repository),
		ListDocuments:      application.NewListDocumentsService(repository),
	})

	router := api.NewExpandedRouter(applicationHandlers, workflowHandlers)

	server := &http.Server{
		Addr:              serverConfig.HTTPAddress,
		Handler:           allowLocalDevelopmentCORS(router),
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("ApplyBy API listening on %s", serverConfig.HTTPAddress)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown server: %w", err)
		}

		return nil
	case err := <-serverErrors:
		if err == http.ErrServerClosed {
			return nil
		}

		return fmt.Errorf("listen and serve: %w", err)
	}
}

// -----------------------------------------------------------------------------
// openDatabase
//
// Opens and verifies a PostgreSQL database connection.
// -----------------------------------------------------------------------------
func openDatabase(ctx context.Context, databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

// -----------------------------------------------------------------------------
// allowLocalDevelopmentCORS
//
// Adds local development CORS headers so the Vite frontend can call the Go API.
// This is intentionally narrow and should be revisited before production deployment.
// -----------------------------------------------------------------------------
func allowLocalDevelopmentCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		origin := request.Header.Get("Origin")

		if origin == "http://localhost:5173" || origin == "http://127.0.0.1:5173" {
			response.Header().Set("Access-Control-Allow-Origin", origin)
			response.Header().Set("Vary", "Origin")
			response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}

		if request.Method == http.MethodOptions {
			response.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(response, request)
	})
}
