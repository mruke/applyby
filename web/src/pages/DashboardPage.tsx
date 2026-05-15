import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";

import { getApplications } from "../api/applications";
import { ErrorState } from "../components/ErrorState";
import { LoadingState } from "../components/LoadingState";
import { StatusBadge } from "../components/StatusBadge";
import type { ApplicationResponse, ApplicationStatus } from "../types/application";
import { formatShortDate } from "../utils/dateFormatting";

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
 * DashboardFilter
 *
 * Represents the active application summary filter on the dashboard.
 */
type DashboardFilter = "all" | "active" | "interviewing" | "offer";

/**
 * DashboardFilterCard
 *
 * Defines one dashboard summary card and its filter behavior.
 */
type DashboardFilterCard = {
  key: DashboardFilter;
  title: string;
  description: string;
  count: number;
};

/**
 * activeStatuses
 *
 * Defines statuses counted as active work on the dashboard.
 */
const activeStatuses: ApplicationStatus[] = ["draft", "interested", "applied", "interviewing", "offer"];

/**
 * filterApplications
 *
 * Returns applications matching the selected dashboard summary filter.
 */
function filterApplications(applications: ApplicationResponse[], filter: DashboardFilter): ApplicationResponse[] {
  switch (filter) {
    case "active":
      return applications.filter((application) => activeStatuses.includes(application.status));
    case "interviewing":
      return applications.filter((application) => application.status === "interviewing");
    case "offer":
      return applications.filter((application) => application.status === "offer");
    case "all":
    default:
      return applications;
  }
}

/**
 * sortNewestFirst
 *
 * Sorts applications by creation date descending for dashboard access.
 */
function sortNewestFirst(applications: ApplicationResponse[]): ApplicationResponse[] {
  return [...applications].sort(
    (left, right) => new Date(right.created_at).getTime() - new Date(left.created_at).getTime()
  );
}

/**
 * filterLabel
 *
 * Returns human-readable copy for the selected dashboard filter.
 */
function filterLabel(filter: DashboardFilter): string {
  switch (filter) {
    case "active":
      return "Active applications";
    case "interviewing":
      return "Interviewing applications";
    case "offer":
      return "Offer applications";
    case "all":
    default:
      return "All applications";
  }
}

/**
 * DashboardPage
 *
 * Provides a summary-first landing page for filtering and opening applications.
 */
export function DashboardPage() {
  const [selectedFilter, setSelectedFilter] = useState<DashboardFilter>("all");
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

  const filteredApplications = useMemo(
    () => sortNewestFirst(filterApplications(state.applications, selectedFilter)),
    [selectedFilter, state.applications]
  );

  const activeCount = state.applications.filter((application) => activeStatuses.includes(application.status)).length;
  const interviewingCount = state.applications.filter((application) => application.status === "interviewing").length;
  const offerCount = state.applications.filter((application) => application.status === "offer").length;

  const filterCards: DashboardFilterCard[] = [
    {
      key: "all",
      title: "Total applications",
      description: "Show every tracked application",
      count: state.applications.length
    },
    {
      key: "active",
      title: "Active applications",
      description: "Show applications still in motion",
      count: activeCount
    },
    {
      key: "interviewing",
      title: "Interviewing",
      description: "Show interview-stage applications",
      count: interviewingCount
    },
    {
      key: "offer",
      title: "Offers",
      description: "Show applications with offers",
      count: offerCount
    }
  ];

  return (
    <>
      <header className="page-header">
        <h1>Dashboard</h1>
        <p>Use the summary cards to filter quick-access applications, or open the applications page for full workflows.</p>
      </header>

      {state.isLoading ? (
        <LoadingState message="Loading dashboard..." />
      ) : state.errorMessage ? (
        <ErrorState title="Dashboard could not be loaded" message={state.errorMessage} />
      ) : (
        <div className="dashboard-layout">
          <section className="dashboard-grid" aria-label="Application summary filters">
            {filterCards.map((card) => (
              <button
                key={card.key}
                type="button"
                className="metric-card metric-card--button"
                aria-pressed={selectedFilter === card.key}
                onClick={() => setSelectedFilter(card.key)}
              >
                <h2>{card.title}</h2>
                <p>{card.count}</p>
                <span>{card.description}</span>
              </button>
            ))}
          </section>

          <section className="state-card" aria-labelledby="dashboard-applications-heading">
            <div className="section-heading-row">
              <div>
                <h2 id="dashboard-applications-heading">Applications</h2>
                <p>
                  {filterLabel(selectedFilter)} · {filteredApplications.length} shown
                </p>
              </div>
              <Link className="secondary-button" to="/applications">
                Open applications workbench
              </Link>
            </div>

            {filteredApplications.length === 0 ? (
              <p>No applications match this summary view.</p>
            ) : (
              <ul className="stack-list dashboard-application-list">
                {filteredApplications.map((application) => (
                  <li key={application.id} className="stack-list__item">
                    <div>
                      <strong>
                        <Link to={`/applications/${application.id}`}>{application.title}</Link>
                      </strong>
                      <p>
                        {application.company_name} · Added {formatShortDate(application.created_at)}
                      </p>
                    </div>
                    <StatusBadge status={application.status} />
                  </li>
                ))}
              </ul>
            )}
          </section>
        </div>
      )}
    </>
  );
}