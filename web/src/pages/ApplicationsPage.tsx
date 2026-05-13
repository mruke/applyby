import { EmptyState } from "../components/EmptyState";

/**
 * ApplicationsPage
 *
 * Provides the applications route placeholder for the frontend foundation.
 * This page will later contain the searchable application list.
 */
export function ApplicationsPage() {
  return (
    <>
      <header className="page-header">
        <h1>Applications</h1>
        <p>Search, filter, and scan tracked applications.</p>
      </header>

      <EmptyState
        title="Applications foundation"
        message="The application list and filters will be implemented in the UI step."
      />
    </>
  );
}