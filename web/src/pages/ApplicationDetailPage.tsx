import { useCallback, useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";

import { getApplicationById, updateApplicationStatus } from "../api/applications";
import { EmptyState } from "../components/EmptyState";
import { ErrorState } from "../components/ErrorState";
import { LoadingState } from "../components/LoadingState";
import { StatusBadge } from "../components/StatusBadge";
import { StatusUpdateForm } from "../components/StatusUpdateForm";
import type { ApplicationResponse, ApplicationStatus } from "../types/application";

/**
 * ApplicationDetailPageState
 *
 * Represents the loading and status update state for the application detail page.
 */
type ApplicationDetailPageState = {
  application: ApplicationResponse | null;
  errorMessage: string | null;
  isLoading: boolean;
  isSubmittingStatus: boolean;
  successMessage: string | null;
};

/**
 * formatDate
 *
 * Converts an API timestamp into a readable date label for detail display.
 */
function formatDate(timestamp: string): string {
  return new Intl.DateTimeFormat(undefined, {
    year: "numeric",
    month: "long",
    day: "numeric"
  }).format(new Date(timestamp));
}

/**
 * ApplicationDetailPage
 *
 * Loads one application from the existing applications API and exposes
 * summary details plus a status update workflow.
 */
export function ApplicationDetailPage() {
  const { applicationId } = useParams<{ applicationId: string }>();

  const [state, setState] = useState<ApplicationDetailPageState>({
    application: null,
    errorMessage: null,
    isLoading: true,
    isSubmittingStatus: false,
    successMessage: null
  });

  /**
   * loadApplication
   *
   * Loads the current application by route identity.
   */
  const loadApplication = useCallback(async () => {
    if (!applicationId) {
      setState({
        application: null,
        errorMessage: "Application id is missing from the route.",
        isLoading: false,
        isSubmittingStatus: false,
        successMessage: null
      });
      return;
    }

    const application = await getApplicationById(applicationId);

    setState((currentState) => ({
      ...currentState,
      application,
      errorMessage: null,
      isLoading: false
    }));
  }, [applicationId]);

  useEffect(() => {
    let isCurrentRequest = true;

    async function loadInitialApplication() {
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

    void loadInitialApplication();

    return () => {
      isCurrentRequest = false;
    };
  }, [applicationId]);

  /**
   * handleStatusUpdate
   *
   * Updates the application status, refreshes the detail state, and displays feedback.
   */
  async function handleStatusUpdate(status: ApplicationStatus) {
    if (!applicationId) {
      return;
    }

    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isSubmittingStatus: true,
      successMessage: null
    }));

    try {
      const updatedApplication = await updateApplicationStatus(applicationId, status);

      setState((currentState) => ({
        ...currentState,
        application: updatedApplication,
        errorMessage: null,
        isSubmittingStatus: false,
        successMessage: "Status updated."
      }));

      await loadApplication();
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Status could not be updated. Check the selected status and try again.",
        isSubmittingStatus: false,
        successMessage: null
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
        <Link to="/applications">Back to applications</Link>
        <h1>{state.application.title}</h1>
        <p>{state.application.company_name}</p>
      </header>

      {state.successMessage ? (
        <p className="form-message form-message--success" role="status">
          {state.successMessage}
        </p>
      ) : null}

      {state.errorMessage ? (
        <p className="form-message form-message--error" role="alert">
          {state.errorMessage}
        </p>
      ) : null}

      <section className="detail-grid" aria-label="Application details">
        <article className="state-card">
          <h2>Summary</h2>
          <dl className="detail-list">
            <div>
              <dt>Status</dt>
              <dd>
                <StatusBadge status={state.application.status} />
              </dd>
            </div>
            <div>
              <dt>Company website</dt>
              <dd>
                {state.application.company_website ? (
                  <a href={state.application.company_website}>{state.application.company_website}</a>
                ) : (
                  "Not specified"
                )}
              </dd>
            </div>
            <div>
              <dt>Source</dt>
              <dd>{state.application.source || "Not specified"}</dd>
            </div>
            <div>
              <dt>Created</dt>
              <dd>{formatDate(state.application.created_at)}</dd>
            </div>
            <div>
              <dt>Notes</dt>
              <dd>{state.application.notes || "No notes added yet."}</dd>
            </div>
          </dl>
        </article>

        <StatusUpdateForm
          currentStatus={state.application.status}
          isSubmitting={state.isSubmittingStatus}
          onSubmit={handleStatusUpdate}
        />
      </section>
    </>
  );
}