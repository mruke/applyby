import { useParams } from "react-router-dom";

import { EmptyState } from "../components/EmptyState";

/**
 * ApplicationDetailPage
 *
 * Provides the application detail route placeholder for the frontend foundation.
 * This page will later show status, reminders, contacts, documents, and activity.
 */
export function ApplicationDetailPage() {
  const { applicationId } = useParams();

  return (
    <>
      <header className="page-header">
        <h1>Application Detail</h1>
        <p>View status, reminders, contacts, documents, and activity for one application.</p>
      </header>

      <EmptyState
        title="Application detail foundation"
        message={`Detail sections will be implemented after the frontend foundation is validated. Application: ${
          applicationId ?? "unknown"
        }.`}
      />
    </>
  );
}