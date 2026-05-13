import { useEffect, useState } from "react";

import { getApplications } from "../api/applications";
import { ApplicationList } from "../components/ApplicationList";
import { EmptyState } from "../components/EmptyState";
import { ErrorState } from "../components/ErrorState";
import { LoadingState } from "../components/LoadingState";
import type { ApplicationResponse } from "../types/application";

/**
 * ApplicationsPageState
 *
 * Represents the frontend loading state for the applications page.
 */
type ApplicationsPageState = {
  applications: ApplicationResponse[];
  errorMessage: string | null;
  isLoading: boolean;
};

/**
 * ApplicationsPage
 *
 * Loads and displays tracked applications from the backend API.
 * The page owns loading, error, empty, and populated states while delegating
 * list rendering to the ApplicationList component.
 */
export function ApplicationsPage() {
  const [state, setState] = useState<ApplicationsPageState>({
    applications: [],
    errorMessage: null,
    isLoading: true
  });

  useEffect(() => {
    let isCurrentRequest = true;

    async function loadApplications() {
      try {
        const response = await getApplications();

        if (!isCurrentRequest) {
          return;
        }

        setState({
          applications: response.applications,
          errorMessage: null,
          isLoading: false
        });
      } catch {
        if (!isCurrentRequest) {
          return;
        }

        setState({
          applications: [],
          errorMessage: "Applications could not be loaded. Check that the backend is running and try again.",
          isLoading: false
        });
      }
    }

    void loadApplications();

    return () => {
      isCurrentRequest = false;
    };
  }, []);

  return (
    <>
      <header className="page-header">
        <h1>Applications</h1>
        <p>Search, filter, and scan tracked applications.</p>
      </header>

      {state.isLoading ? (
        <LoadingState message="Loading applications..." />
      ) : state.errorMessage ? (
        <ErrorState title="Applications could not be loaded" message={state.errorMessage} />
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