package domain

// -----------------------------------------------------------------------------
// DocumentID
//
// Represents the stable domain identity for document metadata.
// -----------------------------------------------------------------------------
type DocumentID string

// -----------------------------------------------------------------------------
// Validate
//
// Verifies that a document identity contains a non-empty value.
// -----------------------------------------------------------------------------
func (id DocumentID) Validate() error {
	return requireNonEmptyField("document id", string(id))
}

// -----------------------------------------------------------------------------
// String
//
// Returns the text representation of a document identity.
// -----------------------------------------------------------------------------
func (id DocumentID) String() string {
	return string(id)
}

// -----------------------------------------------------------------------------
// Document
//
// Represents metadata for a resume, cover letter, portfolio file, or attachment.
// -----------------------------------------------------------------------------
type Document struct {
	ID            DocumentID
	ApplicationID ApplicationID
	Name          string
	Kind          string
	Path          string
}

// -----------------------------------------------------------------------------
// NewDocument
//
// Creates document metadata after applying basic domain validation.
// -----------------------------------------------------------------------------
func NewDocument(id DocumentID, applicationID ApplicationID, name string, kind string, path string) (Document, error) {
	document := Document{
		ID:            id,
		ApplicationID: applicationID,
		Name:          name,
		Kind:          kind,
		Path:          path,
	}

	if err := document.Validate(); err != nil {
		return Document{}, err
	}

	return document, nil
}

// -----------------------------------------------------------------------------
// Validate
//
// Verifies that document metadata has identity, application ownership, and required identifying fields.
// -----------------------------------------------------------------------------
func (document Document) Validate() error {
	if err := document.ID.Validate(); err != nil {
		return err
	}

	if err := document.ApplicationID.Validate(); err != nil {
		return err
	}

	if err := requireNonEmptyField("document name", document.Name); err != nil {
		return err
	}

	return requireNonEmptyField("document kind", document.Kind)
}
