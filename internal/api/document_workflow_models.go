package api

import (
	"github.com/mruke/applyby/internal/application"
	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// documentRequest
//
// Represents the JSON request body for adding or updating document metadata.
// -----------------------------------------------------------------------------
type documentRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"`
	Path string `json:"path"`
}

// -----------------------------------------------------------------------------
// documentResponse
//
// Represents the JSON response shape for document metadata.
// -----------------------------------------------------------------------------
type documentResponse struct {
	ID            string `json:"id"`
	ApplicationID string `json:"application_id"`
	Name          string `json:"name"`
	Kind          string `json:"kind"`
	Path          string `json:"path"`
}

// -----------------------------------------------------------------------------
// documentsResponse
//
// Represents the JSON response shape for a document metadata collection.
// -----------------------------------------------------------------------------
type documentsResponse struct {
	Documents []documentResponse `json:"documents"`
}

// -----------------------------------------------------------------------------
// toInput
//
// Converts a document request into an add document workflow input model.
// -----------------------------------------------------------------------------
func (request documentRequest) toInput(applicationID domain.ApplicationID) application.AddDocumentInput {
	return application.AddDocumentInput{
		ID:            domain.DocumentID(request.ID),
		ApplicationID: applicationID,
		Name:          request.Name,
		Kind:          request.Kind,
		Path:          request.Path,
	}
}

// -----------------------------------------------------------------------------
// toUpdateInput
//
// Converts a document request into an update document workflow input model.
// -----------------------------------------------------------------------------
func (request documentRequest) toUpdateInput(applicationID domain.ApplicationID, documentID domain.DocumentID) application.UpdateDocumentInput {
	return application.UpdateDocumentInput{
		ApplicationID: applicationID,
		DocumentID:    documentID,
		Name:          request.Name,
		Kind:          request.Kind,
		Path:          request.Path,
	}
}

// -----------------------------------------------------------------------------
// documentToResponse
//
// Converts domain document metadata into an API response model.
// -----------------------------------------------------------------------------
func documentToResponse(document domain.Document) documentResponse {
	return documentResponse{
		ID:            document.ID.String(),
		ApplicationID: document.ApplicationID.String(),
		Name:          document.Name,
		Kind:          document.Kind,
		Path:          document.Path,
	}
}

// -----------------------------------------------------------------------------
// documentsToResponse
//
// Converts domain documents into an API collection response model.
// -----------------------------------------------------------------------------
func documentsToResponse(documents []domain.Document) documentsResponse {
	response := documentsResponse{
		Documents: make([]documentResponse, 0, len(documents)),
	}

	for _, document := range documents {
		response.Documents = append(response.Documents, documentToResponse(document))
	}

	return response
}
