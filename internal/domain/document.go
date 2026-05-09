package domain

// -----------------------------------------------------------------------------
// Document
//
// Represents metadata for a resume, cover letter, portfolio file, or attachment.
// -----------------------------------------------------------------------------
type Document struct {
	Name string
	Kind string
	Path string
}

// -----------------------------------------------------------------------------
// NewDocument
//
// Creates document metadata after applying basic domain validation.
// -----------------------------------------------------------------------------
func NewDocument(name string, kind string, path string) (Document, error) {
	document := Document{
		Name: name,
		Kind: kind,
		Path: path,
	}

	if err := document.Validate(); err != nil {
		return Document{}, err
	}

	return document, nil
}

// -----------------------------------------------------------------------------
// Validate
//
// Verifies that document metadata has the minimum required identifying fields.
// -----------------------------------------------------------------------------
func (document Document) Validate() error {
	if err := requireNonEmptyField("document name", document.Name); err != nil {
		return err
	}

	return requireNonEmptyField("document kind", document.Kind)
}
