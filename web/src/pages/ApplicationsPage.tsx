import { useCallback, useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";

import { createApplication, getApplications, searchApplications } from "../api/applications";
import { ApplicationForm } from "../components/ApplicationForm";
import { ApplicationList } from "../components/ApplicationList";
import { ApplicationSearchForm } from "../components/ApplicationSearchForm";
import { EmptyState } from "../components/EmptyState";
import { ErrorState } from "../components/ErrorState";
import { LoadingState } from "../components/LoadingState";
import { applicationStatusOptions } from "../constants/applicationStatuses";
import type { ApplicationResponse, ApplicationSearchCriteria, ApplicationStatus, CreateApplicationFormValues } from "../types/application";

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
 * isApplicationStatus
 *
 * Reports whether a query string value is a supported application status.
 */
function isApplicationStatus(value: string): value is ApplicationStatus {
  return applicationStatusOptions.some((status) => status.value === value);
}

/**
 * searchCriteriaFromParams
 *
 * Converts URL query params into application search criteria.
 */
function searchCriteriaFromParams(searchParams: URLSearchParams): ApplicationSearchCriteria {
  return {
    companyName: searchParams.get("company_name") ?? "",
    source: searchParams.get("source") ?? "",
    statuses: searchParams.getAll("status").filter(isApplicationStatus),
    text: searchParams.get("text") ?? ""
  };
}

/**
 * paramsFromSearchCriteria
 *
 * Converts application search criteria into URL query params.
 */
function paramsFromSearchCriteria(criteria: ApplicationSearchCriteria): URLSearchParams {
  const params = new URLSearchParams();

  if (criteria.companyName.trim() !== "") {
    params.set("company_name", criteria.companyName.trim());
  }

  if (criteria.source.trim() !== "") {
    params.set("source", criteria.source.trim());
  }

  if (criteria.text.trim() !== "") {
    params.set("text", criteria.text.trim());
  }

  for (const status of criteria.statuses) {
    params.append("status", status);
  }

  return params;
}

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
  const [searchParams, setSearchParams] = useSearchParams();
  const initialSearchCriteria = searchCriteriaFromParams(searchParams);

  const [state, setState] = useState<ApplicationsPageState>({
    applications: [],
    errorMessage: null,
    isLoading: true,
    isSearching: false,
    isSubmitting: false,
    searchCriteria: initialSearchCriteria,
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
        const criteria = searchCriteriaFromParams(searchParams);
        const response = hasActiveSearchCriteria(criteria) ? await searchApplications(criteria) : await getApplications();

        if (!isCurrentRequest) {
          return;
        }

        setState((currentState) => ({
          ...currentState,
          applications: response.applications,
          errorMessage: null,
          isLoading: false,
          searchCriteria: criteria
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
      setSearchParams({});
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
      setSearchParams(paramsFromSearchCriteria(criteria));

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
      setSearchParams({});
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