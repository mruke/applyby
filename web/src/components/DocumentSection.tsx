import { DocumentForm } from "./DocumentForm";
import { DocumentList } from "./DocumentList";
import type { CreateDocumentFormValues, DocumentResponse } from "../types/application";

// -----------------------------------------------------------------------------
// DocumentSectionProps
//
// Defines document section data and workflow callbacks for the detail page.
// -----------------------------------------------------------------------------
type DocumentSectionProps = {
  applicationId: string;
  documents: DocumentResponse[];
  errorMessage: string | null;
  isAdding: boolean;
  isRemoving: boolean;
  onAdd: (values: CreateDocumentFormValues) => Promise<void>;
  onRemove: (documentId: string) => Promise<void>;
};

// -----------------------------------------------------------------------------
// DocumentSection
//
// Renders document creation, document load errors, and document maintenance actions.
// -----------------------------------------------------------------------------
export function DocumentSection({
  applicationId,
  documents,
  errorMessage,
  isAdding,
  isRemoving,
  onAdd,
  onRemove
}: DocumentSectionProps) {
  return (
    <>
      <DocumentForm isSubmitting={isAdding} onSubmit={onAdd} />

      {errorMessage ? (
        <section className="state-card" aria-labelledby="documents-heading">
          <h2 id="documents-heading">Documents</h2>
          <p className="form-message form-message--error" role="alert">
            {errorMessage}
          </p>
        </section>
      ) : (
        <DocumentList applicationId={applicationId} documents={documents} isRemoving={isRemoving} onRemove={onRemove} />
      )}
    </>
  );
}