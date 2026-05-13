import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";

import { getApplicationById, updateApplicationDetails } from "../api/applications";
import { ApplicationEditForm } from "../components/ApplicationEditForm";
import { EmptyState } from "../components/EmptyState";
import { ErrorState } from "../components/ErrorState";
import { LoadingState } from "../components/LoadingState";
import type { ApplicationResponse, UpdateApplicationDetailsFormValues } from "../types/application";

/**
 * ApplicationEditPageState
 *
 * Represents loading and submit state for editing application details.
 */
type ApplicationEditPageState = {
  application: ApplicationResponse | null;
  errorMessage: string | null;
  isLoading: boolean;
  isSubmitting: boolean;
};

/**
 * ApplicationEditPage
 *
 * Loads one application and renders the non-status application detail edit workflow.
 */
export function ApplicationEditPage() {
  const navigate = useNavigate();
  const { applicationId } = useParams<{ applicationId: string }>();

  const [state, setState] = useState<ApplicationEditPageState>({
    application: null,
    errorMessage: null,
    isLoading: true,
    isSubmitting: false
  });

  useEffect(() => {
    let isCurrentRequest = true;

    async function loadApplication() {
      try {
        if (!applicationId) {
          throw new Error("missing application id");
        }

        const application = await getApplicationById(applicationId);

        if (!isCurrentRequest) {
          return;
        }

        setState((currentState) => ({
          ...currentState,
          application,
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
          errorMessage: "Application could not be loaded. Check that the backend is running and try again.",
          isLoading: false
        }));
      }
    }

    void loadApplication();

    return () => {
      isCurrentRequest = false;
    };
  }, [applicationId]);

  /**
   * handleSubmit
   *
   * Updates application details and returns to the application detail page.
   */
  async function handleSubmit(values: UpdateApplicationDetailsFormValues) {
    if (!applicationId) {
      return;
    }

    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isSubmitting: true
    }));

    try {
      await updateApplicationDetails(applicationId, values);
      void navigate(`/applications/${applicationId}`);
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Application details could not be updated. Check the form and try again.",
        isSubmitting: false
      }));
    }
  }

  if (state.isLoading) {
    return <LoadingState message="Loading application..." />;
  }

  if (state.errorMessage && !state.application) {
    return <ErrorState title="Application could not be loaded" message={state.errorMessage} />;
  }

  if (!state.application) {
    return (
      <EmptyState
        title="Application not found"
        message="No application matched this route. Return to the applications list and choose an existing application."
      />
    );
  }

  return (
    <>
      <header className="page-header">
        <Link to={`/applications/${state.application.id}`}>Back to application</Link>
        <h1>Edit application</h1>
        <p>
          {state.application.title} · {state.application.company_name}
        </p>
      </header>

      {state.errorMessage ? (
        <p className="form-message form-message--error" role="alert">
          {state.errorMessage}
        </p>
      ) : null}

      <ApplicationEditForm application={state.application} isSubmitting={state.isSubmitting} onSubmit={handleSubmit} />

      <p>
        <Link className="secondary-button" to={`/applications/${state.application.id}`}>
          Cancel
        </Link>
      </p>
    </>
  );
}