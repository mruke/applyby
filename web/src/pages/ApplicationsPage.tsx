import { useCallback, useEffect, useState } from "react";

import { createApplication, getApplications, searchApplications } from "../api/applications";
import { ApplicationForm } from "../components/ApplicationForm";
import { ApplicationList } from "../components/ApplicationList";
import { ApplicationSearchForm } from "../components/ApplicationSearchForm";
import { EmptyState } from "../components/EmptyState";
import { ErrorState } from "../components/ErrorState";
import { LoadingState } from "../components/LoadingState";
import type { ApplicationResponse, ApplicationSearchCriteria, CreateApplicationFormValues } from "../types/application";

/**
 * emptySearchCriteria
 *
 * Provides default criteria for an unfiltered applications list.
 */
const emptySearchCriteria: ApplicationSearchCriteria = {
  companyName: "",
  source: "",
  statuses: [],
  text: ""
};

/**
 * ApplicationsPageState
 *
 * Represents the frontend state for loading, displaying, searching, and creating applications.
 */
type ApplicationsPageState = {
  applications: ApplicationResponse[];
  errorMessage: string | null;
  isLoading: boolean;
  isSearching: boolean;
  isSubmitting: boolean;
  searchCriteria: ApplicationSearchCriteria;
  successMessage: string | null;
};

/**
 * hasActiveSearchCriteria
 *
 * Reports whether any application search criteria are active.
 */
function hasActiveSearchCriteria(criteria: ApplicationSearchCriteria): boolean {
  return (
    criteria.companyName.trim() !== "" ||
    criteria.source.trim() !== "" ||
    criteria.text.trim() !== "" ||
    criteria.statuses.length > 0
  );
}

/**
 * ApplicationsPage
 *
 * Loads tracked applications and provides frontend workflows for creating,
 * searching, filtering, and clearing the applications list.
 */
export function ApplicationsPage() {
  const [state, setState] = useState<ApplicationsPageState>({
    applications: [],
    errorMessage: null,
    isLoading: true,
    isSearching: false,
    isSubmitting: false,
    searchCriteria: emptySearchCriteria,
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
      isLoading: false,
      isSearching: false,
      searchCriteria: emptySearchCriteria
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

  /**
   * handleSearchApplications
   *
   * Runs application search and exposes the active criteria in page state.
   */
  async function handleSearchApplications(criteria: ApplicationSearchCriteria) {
    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isSearching: true,
      searchCriteria: criteria,
      successMessage: null
    }));

    try {
      const response = await searchApplications(criteria);

      setState((currentState) => ({
        ...currentState,
        applications: response.applications,
        errorMessage: null,
        isLoading: false,
        isSearching: false,
        searchCriteria: criteria
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        applications: [],
        errorMessage: "Applications could not be searched. Check the filters and try again.",
        isSearching: false
      }));
    }
  }

  /**
   * handleClearSearch
   *
   * Clears active search criteria and reloads all applications.
   */
  async function handleClearSearch() {
    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isSearching: true,
      searchCriteria: emptySearchCriteria,
      successMessage: null
    }));

    try {
      await loadApplications();
    } catch {
      setState((currentState) => ({
        ...currentState,
        applications: [],
        errorMessage: "Applications could not be loaded. Check that the backend is running and try again.",
        isSearching: false
      }));
    }
  }

  return (
    <>
      <header className="page-header">
        <h1>Applications</h1>
        <p>Add new opportunities while keeping the current application list and filters visible.</p>
      </header>

      {state.successMessage ? (
        <p className="form-message form-message--success" role="status">
          {state.successMessage}
        </p>
      ) : null}

      <div className="applications-workspace">
        <div className="applications-workspace__pane applications-workspace__pane--sticky">
          <ApplicationForm isSubmitting={state.isSubmitting} onSubmit={handleCreateApplication} />
        </div>

        <div className="applications-workspace__pane">
          <section className="state-card" aria-labelledby="applications-list-heading">
            <div className="section-heading-row">
              <div>
                <h2 id="applications-list-heading">Tracked applications</h2>
                <p>{state.applications.length} currently shown</p>
              </div>
            </div>

            <ApplicationSearchForm
              criteria={state.searchCriteria}
              isSearching={state.isSearching}
              onClear={handleClearSearch}
              onSearch={handleSearchApplications}
            />

            {hasActiveSearchCriteria(state.searchCriteria) ? (
              <p className="active-filter-summary" role="status">
                Search filters are active.
              </p>
            ) : null}

            {state.isLoading ? (
              <LoadingState message="Loading applications..." />
            ) : state.errorMessage ? (
              <ErrorState title="Applications need attention" message={state.errorMessage} />
            ) : state.applications.length === 0 ? (
              <EmptyState
                title={hasActiveSearchCriteria(state.searchCriteria) ? "No matching applications" : "No applications yet"}
                message={
                  hasActiveSearchCriteria(state.searchCriteria)
                    ? "Try changing or clearing the search filters."
                    : "Add your first application to start tracking your job search."
                }
              />
            ) : (
              <ApplicationList applications={state.applications} />
            )}
          </section>
        </div>
      </div>
    </>
  );
}