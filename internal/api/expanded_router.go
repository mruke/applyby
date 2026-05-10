package api

import "net/http"

// -----------------------------------------------------------------------------
// NewExpandedRouter
//
// Creates the HTTP route table for the full backend workflow API.
// -----------------------------------------------------------------------------
func NewExpandedRouter(applicationHandlers ApplicationHandlers, workflowHandlers WorkflowHandlers) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/applications/search", workflowHandlers.HandleApplicationSearch)
	router.HandleFunc("/applications", applicationHandlers.HandleApplications)
	router.HandleFunc("/applications/", func(response http.ResponseWriter, request *http.Request) {
		if workflowHandlers.HandleApplicationWorkflow(response, request) {
			return
		}

		applicationHandlers.HandleApplicationStatus(response, request)
	})
	router.HandleFunc("/reminders/", workflowHandlers.HandleReminderComplete)

	return router
}
