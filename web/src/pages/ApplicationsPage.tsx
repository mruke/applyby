import { useCallback, useEffect, useState } from "react";

import { createApplication, getApplications } from "../api/applications";
import { ApplicationForm } from "../components/ApplicationForm";
import { ApplicationList } from "../components/ApplicationList";
import { EmptyState } from "../components/EmptyState";
import { ErrorState } from "../components/ErrorState";
import { LoadingState } from "../components/LoadingState";
import type { ApplicationResponse, CreateApplicationFormValues } from "../types/application";

/**
 * ApplicationsPageState
 *
 * Represents the frontend state for loading, displaying, and creating applications.
 */
type ApplicationsPageState = {
  applications: ApplicationResponse[];
  errorMessage: string | null;
  isLoading: boolean;
  isSubmitting: boolean;
  successMessage: string | null;
};

/**
 * ApplicationsPage
 *
 * Loads tracked applications and provides the first frontend write workflow
 * for creating a new application through the backend API.
 */
export function ApplicationsPage() {
  const [state, setState] = useState<ApplicationsPageState>({
    applications: [],
    errorMessage: null,
    isLoading: true,
    isSubmitting: false,
    successMessage: null
  });

  /**
   * loadApplications
   *
   * Loads applications from the backend and updates the page list state.
   */
  const loadApplications = useCallback(async () => {
    const response = await getApplications();

    setState((currentState) => ({
      ...currentState,
      applications: response.applications,
      errorMessage: null,
      isLoading: false
    }));
  }, []);

  useEffect(() => {
    let isCurrentRequest = true;

    async function loadInitialApplications() {
      try {
        const response = await getApplications();

        if (!isCurrentRequest) {
          return;
        }

        setState((currentState) => ({
          ...currentState,
          applications: response.applications,
          errorMessage: null,
          isLoading: false
        }));
      } catch {
        if (!isCurrentRequest) {
          return;
        }

        setState((currentState) => ({
          ...currentState,
          applications: [],
          errorMessage: "Applications could not be loaded. Check that the backend is running and try again.",
          isLoading: false
        }));
      }
    }

    void loadInitialApplications();

    return () => {
      isCurrentRequest = false;
    };
  }, []);

  /**
   * handleCreateApplication
   *
   * Creates an application, refreshes the list, and exposes success or error feedback.
   */
  async function handleCreateApplication(values: CreateApplicationFormValues) {
    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isSubmitting: true,
      successMessage: null
    }));

    try {
      await createApplication(values);
      await loadApplications();

      setState((currentState) => ({
        ...currentState,
        isSubmitting: false,
        successMessage: "Application added."
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Application could not be added. Check the form and try again.",
        isSubmitting: false,
        successMessage: null
      }));
    }
  }

  return (
    <>
      <header className="page-header">
        <h1>Applications</h1>
        <p>Search, filter, and scan tracked applications.</p>
      </header>

      <ApplicationForm isSubmitting={state.isSubmitting} onSubmit={handleCreateApplication} />

      {state.successMessage ? (
        <p className="form-message form-message--success" role="status">
          {state.successMessage}
        </p>
      ) : null}

      {state.isLoading ? (
        <LoadingState message="Loading applications..." />
      ) : state.errorMessage ? (
        <ErrorState title="Applications need attention" message={state.errorMessage} />
      ) : state.applications.length === 0 ? (
        <EmptyState
          title="No applications yet"
          message="Add your first application to start tracking your job search."
        />
      ) : (
        <ApplicationList applications={state.applications} />
      )}
    </>
  );
}