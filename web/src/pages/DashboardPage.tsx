import { EmptyState } from "../components/EmptyState";

/**
 * DashboardPage
 *
 * Provides the dashboard route placeholder for the frontend foundation.
 * This page will later surface reminders, activity, and summaries that need attention.
 */
export function DashboardPage() {
  return (
    <>
      <header className="page-header">
        <h1>Dashboard</h1>
        <p>Review what needs attention across your job search.</p>
      </header>

      <EmptyState
        title="Dashboard foundation"
        message="Dashboard summaries will be added after the frontend foundation is in place."
      />
    </>
  );
}