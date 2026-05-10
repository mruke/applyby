package api

import "net/http"

// -----------------------------------------------------------------------------
// NewRouter
//
// Creates the HTTP route table for the ApplyBy API.
// -----------------------------------------------------------------------------
func NewRouter(applicationHandlers ApplicationHandlers) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/applications", applicationHandlers.HandleApplications)
	router.HandleFunc("/applications/", applicationHandlers.HandleApplicationStatus)

	return router
}
