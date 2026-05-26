import { Link, useNavigate, useParams } from "react-router-dom";

import { ActivityTimeline } from "../components/ActivityTimeline";
import { ContactSection } from "../components/ContactSection";
import { DocumentSection } from "../components/DocumentSection";
import { EmptyState } from "../components/EmptyState";
import { ErrorState } from "../components/ErrorState";
import { LoadingState } from "../components/LoadingState";
import { ReminderSection } from "../components/ReminderSection";
import { StatusBadge } from "../components/StatusBadge";
import { StatusUpdateForm } from "../components/StatusUpdateForm";
import { formatLongDate } from "../utils/dateFormatting";
import { useApplicationDetailPage } from "./useApplicationDetailPage";

/**
 * ApplicationDetailPage
 *
 * Loads one application and renders detail workflows for status, reminders,
 * contacts, document metadata, and activity history.
 */
export function ApplicationDetailPage() {
  const navigate = useNavigate();
  const { applicationId } = useParams<{ applicationId: string }>();

  const { state, actions } = useApplicationDetailPage({
    applicationId,
    onApplicationRemoved: () => {
      void navigate("/applications");
    }
  });

  if (state.isLoading) {
    return <LoadingState message="Loading application..." />;
  }

  if (state.errorMessage && !state.application) {
    return <ErrorState title="Application could not be loaded" message={state.errorMessage} />;
  }

  if (!state.application) {
    return (
      <EmptyState
        title="Application not found"
        message="No application matched this route. Return to the applications list and choose an existing application."
      />
    );
  }

  return (
    <>
      <header className="page-header">
        <Link to="/applications">Back to applications</Link>
        <h1>{state.application.title}</h1>
        <p>{state.application.company_name}</p>
      </header>

      {state.successMessage ? (
        <p className="form-message form-message--success" role="status">
          {state.successMessage}
        </p>
      ) : null}

      {state.errorMessage ? (
        <p className="form-message form-message--error" role="alert">
          {state.errorMessage}
        </p>
      ) : null}

      <section className="detail-grid" aria-label="Application details">
        <article className="state-card">
          <h2>Summary</h2>
          <dl className="detail-list">
            <div>
              <dt>Status</dt>
              <dd>
                <StatusBadge status={state.application.status} />
              </dd>
            </div>
            <div>
              <dt>Company website</dt>
              <dd>
                {state.application.company_website ? (
                  <a href={state.application.company_website}>{state.application.company_website}</a>
                ) : (
                  "Not specified"
                )}
              </dd>
            </div>
            <div>
              <dt>Source</dt>
              <dd>{state.application.source || "Not specified"}</dd>
            </div>
            <div>
              <dt>Created</dt>
              <dd>{formatLongDate(state.application.created_at)}</dd>
            </div>
            <div>
              <dt>Notes</dt>
              <dd>{state.application.notes || "No notes added yet."}</dd>
            </div>
          </dl>
        </article>

        <div className="form-actions">
          <Link className="secondary-button" to={`/applications/${state.application.id}/edit`}>
            Edit application
          </Link>

          <button type="button" onClick={() => void actions.handleRemoveApplication()}>
            Remove application
          </button>
        </div>

        <StatusUpdateForm
          currentStatus={state.application.status}
          isSubmitting={state.isSubmittingStatus}
          onSubmit={actions.handleStatusUpdate}
        />

        <ReminderSection
          applicationId={state.application.id}
          reminders={state.reminders}
          errorMessage={state.sectionErrors.reminders}
          isCompleting={state.isCompletingReminder}
          isRemoving={state.isRemovingReminder}
          isSubmitting={state.isSubmittingReminder}
          onAdd={actions.handleScheduleReminder}
          onComplete={actions.handleCompleteReminder}
          onRemove={actions.handleRemoveReminder}
        />

        {state.sectionErrors.activity ? (
          <section className="state-card" aria-labelledby="activity-heading">
            <h2 id="activity-heading">Activity</h2>
            <p className="form-message form-message--error" role="alert">
              {state.sectionErrors.activity}
            </p>
          </section>
        ) : (
          <ActivityTimeline events={state.activityEvents} />
        )}

        <ContactSection
          applicationId={state.application.id}
          contacts={state.contacts}
          errorMessage={state.sectionErrors.contacts}
          isAdding={state.isAddingContact}
          isRemoving={state.isRemovingContact}
          onAdd={actions.handleAddContact}
          onRemove={actions.handleRemoveContact}
        />

        <DocumentSection
          applicationId={state.application.id}
          documents={state.documents}
          errorMessage={state.sectionErrors.documents}
          isAdding={state.isAddingDocument}
          isRemoving={state.isRemovingDocument}
          onAdd={actions.handleAddDocument}
          onRemove={actions.handleRemoveDocument}
        />
      </section>
    </>
  );
}
