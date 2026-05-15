package api

import (
	"context"
	"net/http"

	"github.com/mruke/applyby/internal/application"
	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// addDocumentExecutor
//
// Defines the application behavior needed by the add document handler.
// -----------------------------------------------------------------------------
type addDocumentExecutor interface {
	Execute(ctx context.Context, input application.AddDocumentInput) (domain.Document, error)
}

// -----------------------------------------------------------------------------
// listDocumentsExecutor
//
// Defines the application behavior needed by the list documents handler.
// -----------------------------------------------------------------------------
type listDocumentsExecutor interface {
	Execute(ctx context.Context, input application.ListDocumentsInput) ([]domain.Document, error)
}

// -----------------------------------------------------------------------------
// updateDocumentExecutor
//
// Defines the application behavior needed by the update document handler.
// -----------------------------------------------------------------------------
type updateDocumentExecutor interface {
	Execute(ctx context.Context, input application.UpdateDocumentInput) (domain.Document, error)
}

// -----------------------------------------------------------------------------
// removeDocumentExecutor
//
// Defines the application behavior needed by the remove document handler.
// -----------------------------------------------------------------------------
type removeDocumentExecutor interface {
	Execute(ctx context.Context, input application.RemoveDocumentInput) error
}

// -----------------------------------------------------------------------------
// DocumentWorkflowHandlers
//
// Groups HTTP handlers for document metadata workflows owned by an application.
// -----------------------------------------------------------------------------
type DocumentWorkflowHandlers struct {
	addDocument    addDocumentExecutor
	listDocuments  listDocumentsExecutor
	updateDocument updateDocumentExecutor
	removeDocument removeDocumentExecutor
}

// -----------------------------------------------------------------------------
// DocumentWorkflowDependencies
//
// Collects document workflow dependencies for API route construction.
// -----------------------------------------------------------------------------
type DocumentWorkflowDependencies struct {
	AddDocument    addDocumentExecutor
	ListDocuments  listDocumentsExecutor
	UpdateDocument updateDocumentExecutor
	RemoveDocument removeDocumentExecutor
}

// -----------------------------------------------------------------------------
// NewDocumentWorkflowHandlers
//
// Creates document workflow handlers from application-layer dependencies.
// -----------------------------------------------------------------------------
func NewDocumentWorkflowHandlers(dependencies DocumentWorkflowDependencies) DocumentWorkflowHandlers {
	return DocumentWorkflowHandlers{
		addDocument:    dependencies.AddDocument,
		listDocuments:  dependencies.ListDocuments,
		updateDocument: dependencies.UpdateDocument,
		removeDocument: dependencies.RemoveDocument,
	}
}

// -----------------------------------------------------------------------------
// HandleCollection
//
// Routes document collection requests for one application by HTTP method.
// -----------------------------------------------------------------------------
func (handlers DocumentWorkflowHandlers) HandleCollection(
	response http.ResponseWriter,
	request *http.Request,
	applicationID domain.ApplicationID,
) {
	switch request.Method {
	case http.MethodGet:
		handlers.handleListDocuments(response, request, applicationID)
	case http.MethodPost:
		handlers.handleAddDocument(response, request, applicationID)
	default:
		writeJSON(response, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
	}
}

// -----------------------------------------------------------------------------
// HandleResource
//
// Routes item-level document maintenance requests for one application by method.
// -----------------------------------------------------------------------------
func (handlers DocumentWorkflowHandlers) HandleResource(
	response http.ResponseWriter,
	request *http.Request,
	applicationID domain.ApplicationID,
	documentID domain.DocumentID,
) {
	switch request.Method {
	case http.MethodPatch:
		handlers.handleUpdateDocument(response, request, applicationID, documentID)
	case http.MethodDelete:
		handlers.handleRemoveDocument(response, request, applicationID, documentID)
	default:
		writeJSON(response, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
	}
}

// -----------------------------------------------------------------------------
// handleAddDocument
//
// Decodes, validates, and executes the add document metadata workflow.
// -----------------------------------------------------------------------------
func (handlers DocumentWorkflowHandlers) handleAddDocument(response http.ResponseWriter, request *http.Request, applicationID domain.ApplicationID) {
	if handlers.addDocument == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "add document service is not configured"})
		return
	}

	var body documentRequest

	if err := decodeJSON(request, &body); err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	document, err := handlers.addDocument.Execute(request.Context(), body.toInput(applicationID))
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(response, http.StatusCreated, documentToResponse(document))
}

// -----------------------------------------------------------------------------
// handleListDocuments
//
// Executes the list documents workflow for one application.
// -----------------------------------------------------------------------------
func (handlers DocumentWorkflowHandlers) handleListDocuments(response http.ResponseWriter, request *http.Request, applicationID domain.ApplicationID) {
	if handlers.listDocuments == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "list documents service is not configured"})
		return
	}

	documents, err := handlers.listDocuments.Execute(request.Context(), application.ListDocumentsInput{
		ApplicationID: applicationID,
	})
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(response, http.StatusOK, documentsToResponse(documents))
}

// -----------------------------------------------------------------------------
// handleUpdateDocument
//
// Decodes, validates, and executes the update document metadata workflow.
// -----------------------------------------------------------------------------
func (handlers DocumentWorkflowHandlers) handleUpdateDocument(
	response http.ResponseWriter,
	request *http.Request,
	applicationID domain.ApplicationID,
	documentID domain.DocumentID,
) {
	if handlers.updateDocument == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "update document service is not configured"})
		return
	}

	var body documentRequest

	if err := decodeJSON(request, &body); err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	document, err := handlers.updateDocument.Execute(request.Context(), body.toUpdateInput(applicationID, documentID))
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(response, http.StatusOK, documentToResponse(document))
}

// -----------------------------------------------------------------------------
// handleRemoveDocument
//
// Executes the remove document metadata workflow for one application document.
// -----------------------------------------------------------------------------
func (handlers DocumentWorkflowHandlers) handleRemoveDocument(
	response http.ResponseWriter,
	request *http.Request,
	applicationID domain.ApplicationID,
	documentID domain.DocumentID,
) {
	if handlers.removeDocument == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "remove document service is not configured"})
		return
	}

	err := handlers.removeDocument.Execute(request.Context(), application.RemoveDocumentInput{
		ApplicationID: applicationID,
		DocumentID:    documentID,
	})
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	response.WriteHeader(http.StatusNoContent)
}
