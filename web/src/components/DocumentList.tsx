import type { DocumentResponse } from "../types/application";

/**
 * DocumentListProps
 *
 * Defines the document metadata records rendered by the document list component.
 */
type DocumentListProps = {
  documents: DocumentResponse[];
};

/**
 * DocumentList
 *
 * Renders document metadata records. This does not perform file upload or download.
 */
export function DocumentList({ documents }: DocumentListProps) {
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
          </li>
        ))}
      </ul>
    </section>
  );
}