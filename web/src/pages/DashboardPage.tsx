import { useEffect, useState } from "react";

import { getApplications } from "../api/applications";
import { ErrorState } from "../components/ErrorState";
import { LoadingState } from "../components/LoadingState";
import type { ApplicationResponse, ApplicationStatus } from "../types/application";

/**
 * DashboardPageState
 *
 * Represents the loading state for dashboard application summaries.
 */
type DashboardPageState = {
  applications: ApplicationResponse[];
  errorMessage: string | null;
  isLoading: boolean;
};

/**
 * activeStatuses
 *
 * Defines statuses counted as active work on the dashboard.
 */
const activeStatuses: ApplicationStatus[] = ["draft", "interested", "applied", "interviewing", "offer"];

/**
 * DashboardPage
 *
 * Provides a basic dashboard overview using the existing applications API.
 */
export function DashboardPage() {
  const [state, setState] = useState<DashboardPageState>({
    applications: [],
    errorMessage: null,
    isLoading: true
  });

  useEffect(() => {
    let isCurrentRequest = true;

    async function loadDashboard() {
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
          errorMessage: "Dashboard could not be loaded. Check that the backend is running and try again.",
          isLoading: false
        });
      }
    }

    void loadDashboard();

    return () => {
      isCurrentRequest = false;
    };
  }, []);

  const activeCount = state.applications.filter((application) => activeStatuses.includes(application.status)).length;
  const interviewingCount = state.applications.filter((application) => application.status === "interviewing").length;
  const offerCount = state.applications.filter((application) => application.status === "offer").length;

  return (
    <>
      <header className="page-header">
        <h1>Dashboard</h1>
        <p>Review what needs attention across your job search.</p>
      </header>

      {state.isLoading ? (
        <LoadingState message="Loading dashboard..." />
      ) : state.errorMessage ? (
        <ErrorState title="Dashboard could not be loaded" message={state.errorMessage} />
      ) : (
        <section className="dashboard-grid" aria-label="Application overview">
          <article className="metric-card">
            <h2>Total applications</h2>
            <p>{state.applications.length}</p>
          </article>

          <article className="metric-card">
            <h2>Active applications</h2>
            <p>{activeCount}</p>
          </article>

          <article className="metric-card">
            <h2>Interviewing</h2>
            <p>{interviewingCount}</p>
          </article>

          <article className="metric-card">
            <h2>Offers</h2>
            <p>{offerCount}</p>
          </article>
        </section>
      )}
    </>
  );
}