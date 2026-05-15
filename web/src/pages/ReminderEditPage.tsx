import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";

import { getApplicationById } from "../api/applications";
import { getReminders, updateReminder } from "../api/reminders";
import { ReminderEditForm } from "../components/ReminderEditForm";
import { EmptyState } from "../components/EmptyState";
import { ErrorState } from "../components/ErrorState";
import { LoadingState } from "../components/LoadingState";
import type { ApplicationResponse, ReminderResponse, UpdateReminderFormValues } from "../types/application";

// -----------------------------------------------------------------------------
// ReminderEditPageState
//
// Represents loading and submit state for editing one reminder.
// -----------------------------------------------------------------------------
type ReminderEditPageState = {
  application: ApplicationResponse | null;
  reminder: ReminderResponse | null;
  errorMessage: string | null;
  isLoading: boolean;
  isSubmitting: boolean;
};

// -----------------------------------------------------------------------------
// findReminderById
//
// Finds one reminder in an application reminder collection.
// -----------------------------------------------------------------------------
function findReminderById(reminders: ReminderResponse[], reminderId: string): ReminderResponse | null {
  return reminders.find((reminder) => reminder.id === reminderId) ?? null;
}

// -----------------------------------------------------------------------------
// ReminderEditPage
//
// Loads an application reminder and renders the reminder edit workflow.
// -----------------------------------------------------------------------------
export function ReminderEditPage() {
  const navigate = useNavigate();
  const { applicationId, reminderId } = useParams<{ applicationId: string; reminderId: string }>();

  const [state, setState] = useState<ReminderEditPageState>({
    application: null,
    reminder: null,
    errorMessage: null,
    isLoading: true,
    isSubmitting: false
  });

  useEffect(() => {
    let isCurrentRequest = true;

    async function loadReminder() {
      try {
        if (!applicationId || !reminderId) {
          throw new Error("missing route ids");
        }

        const [application, remindersResponse] = await Promise.all([
          getApplicationById(applicationId),
          getReminders(applicationId)
        ]);

        if (!isCurrentRequest) {
          return;
        }

        setState((currentState) => ({
          ...currentState,
          application,
          reminder: findReminderById(remindersResponse.reminders, reminderId),
          errorMessage: null,
          isLoading: false
        }));
      } catch {
        if (!isCurrentRequest) {
          return;
        }

        setState((currentState) => ({
          ...currentState,
          application: null,
          reminder: null,
          errorMessage: "Reminder could not be loaded. Check that the backend is running and try again.",
          isLoading: false
        }));
      }
    }

    void loadReminder();

    return () => {
      isCurrentRequest = false;
    };
  }, [applicationId, reminderId]);

  async function handleSubmit(values: UpdateReminderFormValues) {
    if (!applicationId || !reminderId) {
      return;
    }

    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isSubmitting: true
    }));

    try {
      await updateReminder(reminderId, values);
      void navigate(`/applications/${applicationId}`);
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Reminder could not be updated. Check the form and try again.",
        isSubmitting: false
      }));
    }
  }

  if (state.isLoading) {
    return <LoadingState message="Loading reminder..." />;
  }

  if (state.errorMessage && !state.application && !state.reminder) {
    return <ErrorState title="Reminder could not be loaded" message={state.errorMessage} />;
  }

  if (!state.application) {
    return (
      <EmptyState
        title="Application not found"
        message="No application matched this route. Return to the applications list and choose an existing application."
      />
    );
  }

  if (!state.reminder) {
    return (
      <EmptyState
        title="Reminder not found"
        message="No reminder matched this route. Return to the application detail page and choose an existing reminder."
      />
    );
  }

  return (
    <>
      <header className="page-header">
        <Link to={`/applications/${state.application.id}`}>Back to application</Link>
        <h1>Edit reminder</h1>
        <p>
          {state.reminder.title} · {state.application.title}
        </p>
      </header>

      {state.errorMessage ? (
        <p className="form-message form-message--error" role="alert">
          {state.errorMessage}
        </p>
      ) : null}

      <ReminderEditForm reminder={state.reminder} isSubmitting={state.isSubmitting} onSubmit={handleSubmit} />

      <p>
        <Link className="secondary-button" to={`/applications/${state.application.id}`}>
          Cancel
        </Link>
      </p>
    </>
  );
}