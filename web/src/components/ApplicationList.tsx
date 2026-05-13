import { Link } from "react-router-dom";

import type { ApplicationResponse } from "../types/application";
import { StatusBadge } from "./StatusBadge";

/**
 * ApplicationListProps
 *
 * Defines the applications rendered by the application list component.
 */
type ApplicationListProps = {
  applications: ApplicationResponse[];
};

/**
 * formatDate
 *
 * Converts an API timestamp into a readable date label for list display.
 */
function formatDate(timestamp: string): string {
  return new Intl.DateTimeFormat(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric"
  }).format(new Date(timestamp));
}

/**
 * ApplicationList
 *
 * Renders applications in a scan-friendly table with clear labels,
 * readable status text, and links to application detail routes.
 */
export function ApplicationList({ applications }: ApplicationListProps) {
  return (
    <section className="application-list-section" aria-labelledby="applications-list-heading">
      <div className="section-heading-row">
        <h2 id="applications-list-heading">Tracked applications</h2>
        <p>{applications.length} total</p>
      </div>

      <div className="table-scroll">
        <table className="application-table">
          <caption>Tracked job applications</caption>
          <thead>
            <tr>
              <th scope="col">Application</th>
              <th scope="col">Company</th>
              <th scope="col">Status</th>
              <th scope="col">Source</th>
              <th scope="col">Created</th>
            </tr>
          </thead>
          <tbody>
            {applications.map((application) => (
              <tr key={application.id}>
                <td>
                  <Link to={`/applications/${application.id}`}>{application.title}</Link>
                </td>
                <td>{application.company_name}</td>
                <td>
                  <StatusBadge status={application.status} />
                </td>
                <td>{application.source || "Not specified"}</td>
                <td>{formatDate(application.created_at)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  );
}