package api

import (
	"github.com/mruke/applyby/internal/application"
	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// contactRequest
//
// Represents the JSON request body for adding or updating a contact.
// -----------------------------------------------------------------------------
type contactRequest struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// -----------------------------------------------------------------------------
// contactResponse
//
// Represents the JSON response shape for a contact.
// -----------------------------------------------------------------------------
type contactResponse struct {
	ID            string `json:"id"`
	ApplicationID string `json:"application_id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Role          string `json:"role"`
}

// -----------------------------------------------------------------------------
// contactsResponse
//
// Represents the JSON response shape for a contact collection.
// -----------------------------------------------------------------------------
type contactsResponse struct {
	Contacts []contactResponse `json:"contacts"`
}

// -----------------------------------------------------------------------------
// toInput
//
// Converts a contact request into an add contact workflow input model.
// -----------------------------------------------------------------------------
func (request contactRequest) toInput(applicationID domain.ApplicationID) application.AddContactInput {
	return application.AddContactInput{
		ID:            domain.ContactID(request.ID),
		ApplicationID: applicationID,
		Name:          request.Name,
		Email:         request.Email,
		Role:          request.Role,
	}
}

// -----------------------------------------------------------------------------
// toUpdateInput
//
// Converts a contact request into an update contact workflow input model.
// -----------------------------------------------------------------------------
func (request contactRequest) toUpdateInput(applicationID domain.ApplicationID, contactID domain.ContactID) application.UpdateContactInput {
	return application.UpdateContactInput{
		ApplicationID: applicationID,
		ContactID:     contactID,
		Name:          request.Name,
		Email:         request.Email,
		Role:          request.Role,
	}
}

// -----------------------------------------------------------------------------
// contactToResponse
//
// Converts a domain contact into an API response model.
// -----------------------------------------------------------------------------
func contactToResponse(contact domain.Contact) contactResponse {
	return contactResponse{
		ID:            contact.ID.String(),
		ApplicationID: contact.ApplicationID.String(),
		Name:          contact.Name,
		Email:         contact.Email,
		Role:          contact.Role,
	}
}

// -----------------------------------------------------------------------------
// contactsToResponse
//
// Converts domain contacts into an API collection response model.
// -----------------------------------------------------------------------------
func contactsToResponse(contacts []domain.Contact) contactsResponse {
	response := contactsResponse{
		Contacts: make([]contactResponse, 0, len(contacts)),
	}

	for _, contact := range contacts {
		response.Contacts = append(response.Contacts, contactToResponse(contact))
	}

	return response
}
