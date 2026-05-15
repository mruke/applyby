import { Link } from "react-router-dom";

import type { DocumentResponse } from "../types/application";

/**
 * DocumentListProps
 *
 * Defines the document metadata records rendered by the document list component.
 */
type DocumentListProps = {
  applicationId: string;
  documents: DocumentResponse[];
  isRemoving?: boolean;
  onRemove?: (documentId: string) => Promise<void>;
};

/**
 * DocumentList
 *
 * Renders document metadata records. This does not perform file upload or download.
 */
export function DocumentList({ applicationId, documents, isRemoving = false, onRemove }: DocumentListProps) {
  if (documents.length === 0) {
    return (
      <section className="state-card" aria-labelledby="documents-heading">
        <h2 id="documents-heading">Documents</h2>
        <p>No document metadata added yet.</p>
      </section>
    );
  }

  return (
    <section className="state-card" aria-labelledby="documents-heading">
      <h2 id="documents-heading">Documents</h2>

      <ul className="stack-list">
        {documents.map((document) => (
          <li key={document.id} className="stack-list__item">
            <div>
              <strong>{document.name}</strong>
              <p>
                {document.kind}
                {document.path ? ` · ${document.path}` : ""}
              </p>
            </div>

            <div className="stack-list__actions">
              <Link className="secondary-button" to={`/applications/${applicationId}/documents/${document.id}/edit`}>
                Edit
              </Link>

              {onRemove ? (
                <button type="button" disabled={isRemoving} onClick={() => void onRemove(document.id)}>
                  {isRemoving ? "Removing..." : "Remove"}
                </button>
              ) : null}
            </div>
          </li>
        ))}
      </ul>
    </section>
  );
}