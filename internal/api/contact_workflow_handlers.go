package api

import (
	"context"
	"net/http"

	"github.com/mruke/applyby/internal/application"
	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// addContactExecutor
//
// Defines the application behavior needed by the add contact handler.
// -----------------------------------------------------------------------------
type addContactExecutor interface {
	Execute(ctx context.Context, input application.AddContactInput) (domain.Contact, error)
}

// -----------------------------------------------------------------------------
// listContactsExecutor
//
// Defines the application behavior needed by the list contacts handler.
// -----------------------------------------------------------------------------
type listContactsExecutor interface {
	Execute(ctx context.Context, input application.ListContactsInput) ([]domain.Contact, error)
}

// -----------------------------------------------------------------------------
// updateContactExecutor
//
// Defines the application behavior needed by the update contact handler.
// -----------------------------------------------------------------------------
type updateContactExecutor interface {
	Execute(ctx context.Context, input application.UpdateContactInput) (domain.Contact, error)
}

// -----------------------------------------------------------------------------
// removeContactExecutor
//
// Defines the application behavior needed by the remove contact handler.
// -----------------------------------------------------------------------------
type removeContactExecutor interface {
	Execute(ctx context.Context, input application.RemoveContactInput) error
}

// -----------------------------------------------------------------------------
// ContactWorkflowHandlers
//
// Groups HTTP handlers for contact workflows owned by an application.
// -----------------------------------------------------------------------------
type ContactWorkflowHandlers struct {
	addContact    addContactExecutor
	listContacts  listContactsExecutor
	updateContact updateContactExecutor
	removeContact removeContactExecutor
}

// -----------------------------------------------------------------------------
// ContactWorkflowDependencies
//
// Collects contact workflow dependencies for API route construction.
// -----------------------------------------------------------------------------
type ContactWorkflowDependencies struct {
	AddContact    addContactExecutor
	ListContacts  listContactsExecutor
	UpdateContact updateContactExecutor
	RemoveContact removeContactExecutor
}

// -----------------------------------------------------------------------------
// NewContactWorkflowHandlers
//
// Creates contact workflow handlers from application-layer dependencies.
// -----------------------------------------------------------------------------
func NewContactWorkflowHandlers(dependencies ContactWorkflowDependencies) ContactWorkflowHandlers {
	return ContactWorkflowHandlers{
		addContact:    dependencies.AddContact,
		listContacts:  dependencies.ListContacts,
		updateContact: dependencies.UpdateContact,
		removeContact: dependencies.RemoveContact,
	}
}

// -----------------------------------------------------------------------------
// HandleCollection
//
// Routes contact collection requests for one application by HTTP method.
// -----------------------------------------------------------------------------
func (handlers ContactWorkflowHandlers) HandleCollection(
	response http.ResponseWriter,
	request *http.Request,
	applicationID domain.ApplicationID,
) {
	switch request.Method {
	case http.MethodGet:
		handlers.handleListContacts(response, request, applicationID)
	case http.MethodPost:
		handlers.handleAddContact(response, request, applicationID)
	default:
		writeJSON(response, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
	}
}

// -----------------------------------------------------------------------------
// HandleResource
//
// Routes item-level contact maintenance requests for one application by method.
// -----------------------------------------------------------------------------
func (handlers ContactWorkflowHandlers) HandleResource(
	response http.ResponseWriter,
	request *http.Request,
	applicationID domain.ApplicationID,
	contactID domain.ContactID,
) {
	switch request.Method {
	case http.MethodPatch:
		handlers.handleUpdateContact(response, request, applicationID, contactID)
	case http.MethodDelete:
		handlers.handleRemoveContact(response, request, applicationID, contactID)
	default:
		writeJSON(response, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
	}
}

// -----------------------------------------------------------------------------
// handleAddContact
//
// Decodes, validates, and executes the add contact workflow.
// -----------------------------------------------------------------------------
func (handlers ContactWorkflowHandlers) handleAddContact(response http.ResponseWriter, request *http.Request, applicationID domain.ApplicationID) {
	if handlers.addContact == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "add contact service is not configured"})
		return
	}

	var body contactRequest

	if err := decodeJSON(request, &body); err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	contact, err := handlers.addContact.Execute(request.Context(), body.toInput(applicationID))
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(response, http.StatusCreated, contactToResponse(contact))
}

// -----------------------------------------------------------------------------
// handleListContacts
//
// Executes the list contacts workflow for one application.
// -----------------------------------------------------------------------------
func (handlers ContactWorkflowHandlers) handleListContacts(response http.ResponseWriter, request *http.Request, applicationID domain.ApplicationID) {
	if handlers.listContacts == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "list contacts service is not configured"})
		return
	}

	contacts, err := handlers.listContacts.Execute(request.Context(), application.ListContactsInput{
		ApplicationID: applicationID,
	})
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(response, http.StatusOK, contactsToResponse(contacts))
}

// -----------------------------------------------------------------------------
// handleUpdateContact
//
// Decodes, validates, and executes the update contact workflow.
// -----------------------------------------------------------------------------
func (handlers ContactWorkflowHandlers) handleUpdateContact(
	response http.ResponseWriter,
	request *http.Request,
	applicationID domain.ApplicationID,
	contactID domain.ContactID,
) {
	if handlers.updateContact == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "update contact service is not configured"})
		return
	}

	var body contactRequest

	if err := decodeJSON(request, &body); err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	contact, err := handlers.updateContact.Execute(request.Context(), body.toUpdateInput(applicationID, contactID))
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(response, http.StatusOK, contactToResponse(contact))
}

// -----------------------------------------------------------------------------
// handleRemoveContact
//
// Executes the remove contact workflow for one application contact.
// -----------------------------------------------------------------------------
func (handlers ContactWorkflowHandlers) handleRemoveContact(
	response http.ResponseWriter,
	request *http.Request,
	applicationID domain.ApplicationID,
	contactID domain.ContactID,
) {
	if handlers.removeContact == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "remove contact service is not configured"})
		return
	}

	err := handlers.removeContact.Execute(request.Context(), application.RemoveContactInput{
		ApplicationID: applicationID,
		ContactID:     contactID,
	})
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	response.WriteHeader(http.StatusNoContent)
}
