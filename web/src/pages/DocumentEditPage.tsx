import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";

import { getApplicationById } from "../api/applications";
import { getDocuments, updateDocument } from "../api/documents";
import { DocumentEditForm } from "../components/DocumentEditForm";
import { EmptyState } from "../components/EmptyState";
import { ErrorState } from "../components/ErrorState";
import { LoadingState } from "../components/LoadingState";
import type { ApplicationResponse, DocumentResponse, UpdateDocumentFormValues } from "../types/application";

// -----------------------------------------------------------------------------
// DocumentEditPageState
//
// Represents loading and submit state for editing one document metadata record.
// -----------------------------------------------------------------------------
type DocumentEditPageState = {
  application: ApplicationResponse | null;
  document: DocumentResponse | null;
  errorMessage: string | null;
  isLoading: boolean;
  isSubmitting: boolean;
};

// -----------------------------------------------------------------------------
// findDocumentById
//
// Finds one document in an application document collection.
// -----------------------------------------------------------------------------
function findDocumentById(documents: DocumentResponse[], documentId: string): DocumentResponse | null {
  return documents.find((document) => document.id === documentId) ?? null;
}

// -----------------------------------------------------------------------------
// DocumentEditPage
//
// Loads an application document and renders the document metadata edit workflow.
// -----------------------------------------------------------------------------
export function DocumentEditPage() {
  const navigate = useNavigate();
  const { applicationId, documentId } = useParams<{ applicationId: string; documentId: string }>();

  const [state, setState] = useState<DocumentEditPageState>({
    application: null,
    document: null,
    errorMessage: null,
    isLoading: true,
    isSubmitting: false
  });

  useEffect(() => {
    let isCurrentRequest = true;

    async function loadDocument() {
      try {
        if (!applicationId || !documentId) {
          throw new Error("missing route ids");
        }

        const [application, documentsResponse] = await Promise.all([
          getApplicationById(applicationId),
          getDocuments(applicationId)
        ]);

        if (!isCurrentRequest) {
          return;
        }

        setState((currentState) => ({
          ...currentState,
          application,
          document: findDocumentById(documentsResponse.documents, documentId),
          errorMessage: null,
          isLoading: false
        }));
      } catch {
        if (!isCurrentRequest) {
          return;
        }

        setState((currentState) => ({
          ...currentState,
          application: null,
          document: null,
          errorMessage: "Document metadata could not be loaded. Check that the backend is running and try again.",
          isLoading: false
        }));
      }
    }

    void loadDocument();

    return () => {
      isCurrentRequest = false;
    };
  }, [applicationId, documentId]);

  async function handleSubmit(values: UpdateDocumentFormValues) {
    if (!applicationId || !documentId) {
      return;
    }

    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isSubmitting: true
    }));

    try {
      await updateDocument(applicationId, documentId, values);
      void navigate(`/applications/${applicationId}`);
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Document metadata could not be updated. Check the form and try again.",
        isSubmitting: false
      }));
    }
  }

  if (state.isLoading) {
    return <LoadingState message="Loading document metadata..." />;
  }

  if (state.errorMessage && !state.application && !state.document) {
    return <ErrorState title="Document metadata could not be loaded" message={state.errorMessage} />;
  }

  if (!state.application) {
    return (
      <EmptyState
        title="Application not found"
        message="No application matched this route. Return to the applications list and choose an existing application."
      />
    );
  }

  if (!state.document) {
    return (
      <EmptyState
        title="Document metadata not found"
        message="No document metadata matched this route. Return to the application detail page and choose an existing document."
      />
    );
  }

  return (
    <>
      <header className="page-header">
        <Link to={`/applications/${state.application.id}`}>Back to application</Link>
        <h1>Edit document metadata</h1>
        <p>
          {state.document.name} · {state.application.title}
        </p>
      </header>

      {state.errorMessage ? (
        <p className="form-message form-message--error" role="alert">
          {state.errorMessage}
        </p>
      ) : null}

      <DocumentEditForm document={state.document} isSubmitting={state.isSubmitting} onSubmit={handleSubmit} />

      <p>
        <Link className="secondary-button" to={`/applications/${state.application.id}`}>
          Cancel
        </Link>
      </p>
    </>
  );
}