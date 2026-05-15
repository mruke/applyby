package main

import (
	"net/http"

	"github.com/mruke/applyby/internal/api"
	"github.com/mruke/applyby/internal/application"
	"github.com/mruke/applyby/internal/storage/postgres"
)

// -----------------------------------------------------------------------------
// buildRouter
//
// Composes API handlers for the HTTP server.
// -----------------------------------------------------------------------------
func buildRouter(repository postgres.ApplicationRepository) http.Handler {
	applicationHandlers := buildApplicationHandlers(repository)
	workflowHandlers := buildWorkflowHandlers(repository)

	return api.NewExpandedRouter(applicationHandlers, workflowHandlers)
}

// -----------------------------------------------------------------------------
// buildApplicationHandlers
//
// Composes application-resource handlers from application workflow services.
// -----------------------------------------------------------------------------
func buildApplicationHandlers(repository postgres.ApplicationRepository) api.ApplicationHandlers {
	return api.NewApplicationHandlers(
		application.NewCreateApplicationService(repository, repository),
		application.NewListApplicationsService(repository),
		application.NewGetApplicationService(repository),
		application.NewUpdateApplicationDetailsService(repository, repository),
		application.NewUpdateApplicationStatusService(repository, repository),
		application.NewRemoveApplicationService(repository),
	)
}

// -----------------------------------------------------------------------------
// buildWorkflowHandlers
//
// Composes expanded workflow handlers from application workflow services.
// -----------------------------------------------------------------------------
func buildWorkflowHandlers(repository postgres.ApplicationRepository) api.WorkflowHandlers {
	return api.NewWorkflowHandlers(api.WorkflowHandlerDependencies{
		SearchApplications: application.NewSearchApplicationsService(repository),
		ListActivityEvents: application.NewListActivityEventsService(repository),
		Reminders: api.ReminderWorkflowDependencies{
			ScheduleReminder: application.NewScheduleReminderService(repository, repository),
			ListReminders:    application.NewListRemindersService(repository),
			CompleteReminder: application.NewCompleteReminderService(repository, repository),
			UpdateReminder:   application.NewUpdateReminderService(repository, repository),
			RemoveReminder:   application.NewRemoveReminderService(repository, repository),
		},
		Contacts: api.ContactWorkflowDependencies{
			AddContact:    application.NewAddContactService(repository, repository),
			ListContacts:  application.NewListContactsService(repository),
			UpdateContact: application.NewUpdateContactService(repository, repository),
			RemoveContact: application.NewRemoveContactService(repository, repository),
		},
		Documents: api.DocumentWorkflowDependencies{
			AddDocument:    application.NewAddDocumentService(repository, repository),
			ListDocuments:  application.NewListDocumentsService(repository),
			UpdateDocument: application.NewUpdateDocumentService(repository, repository),
			RemoveDocument: application.NewRemoveDocumentService(repository, repository),
		},
	})
}
