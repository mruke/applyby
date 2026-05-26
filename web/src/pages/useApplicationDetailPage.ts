import { useCallback, useEffect, useState } from "react";

import { removeApplication, updateApplicationStatus } from "../api/applications";
import { addContact, removeContact } from "../api/contacts";
import { addDocument, removeDocument } from "../api/documents";
import { completeReminder, removeReminder, scheduleReminder } from "../api/reminders";
import type {
  ApplicationStatus,
  CreateContactFormValues,
  CreateDocumentFormValues,
  CreateReminderFormValues
} from "../types/application";
import { emptySectionErrors, fetchApplicationDetailData } from "./applicationDetailData";
import type { ApplicationDetailData } from "./applicationDetailData";

type UseApplicationDetailPageOptions = {
  applicationId: string | undefined;
  onApplicationRemoved: () => void;
};

export type ApplicationDetailPageState = ApplicationDetailData & {
  errorMessage: string | null;
  isAddingContact: boolean;
  isAddingDocument: boolean;
  isCompletingReminder: boolean;
  isLoading: boolean;
  isRemovingContact: boolean;
  isRemovingDocument: boolean;
  isRemovingReminder: boolean;
  isSubmittingReminder: boolean;
  isSubmittingStatus: boolean;
  successMessage: string | null;
};

function applyDetailData(
  currentState: ApplicationDetailPageState,
  detailData: ApplicationDetailData
): ApplicationDetailPageState {
  return {
    ...currentState,
    activityEvents: detailData.activityEvents,
    application: detailData.application,
    contacts: detailData.contacts,
    documents: detailData.documents,
    errorMessage: null,
    isLoading: false,
    reminders: detailData.reminders,
    sectionErrors: detailData.sectionErrors
  };
}

export function useApplicationDetailPage({
  applicationId,
  onApplicationRemoved
}: UseApplicationDetailPageOptions) {
  const [state, setState] = useState<ApplicationDetailPageState>({
    activityEvents: [],
    application: null,
    contacts: [],
    documents: [],
    errorMessage: null,
    isAddingContact: false,
    isAddingDocument: false,
    isCompletingReminder: false,
    isLoading: true,
    isRemovingContact: false,
    isRemovingDocument: false,
    isRemovingReminder: false,
    isSubmittingReminder: false,
    isSubmittingStatus: false,
    reminders: [],
    sectionErrors: emptySectionErrors(),
    successMessage: null
  });

  const loadDetailData = useCallback(async () => {
    if (!applicationId) {
      setState((currentState) => ({
        ...currentState,
        application: null,
        errorMessage: "Application id is missing from the route.",
        isLoading: false
      }));
      return;
    }

    const detailData = await fetchApplicationDetailData(applicationId);

    setState((currentState) => applyDetailData(currentState, detailData));
  }, [applicationId]);

  useEffect(() => {
    let isCurrentRequest = true;

    async function loadInitialDetailData() {
      try {
        if (!applicationId) {
          throw new Error("missing application id");
        }

        const detailData = await fetchApplicationDetailData(applicationId);

        if (!isCurrentRequest) {
          return;
        }

        setState((currentState) => applyDetailData(currentState, detailData));
      } catch {
        if (!isCurrentRequest) {
          return;
        }

        setState((currentState) => ({
          ...currentState,
          application: null,
          errorMessage: "Application could not be loaded. Check that the backend is running and try again.",
          isLoading: false,
          sectionErrors: emptySectionErrors()
        }));
      }
    }

    void loadInitialDetailData();

    return () => {
      isCurrentRequest = false;
    };
  }, [applicationId]);

  async function handleStatusUpdate(status: ApplicationStatus) {
    if (!applicationId) {
      return;
    }

    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isSubmittingStatus: true,
      successMessage: null
    }));

    try {
      await updateApplicationStatus(applicationId, status);
      await loadDetailData();

      setState((currentState) => ({
        ...currentState,
        isSubmittingStatus: false,
        successMessage: "Status updated."
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Status could not be updated. Check the selected status and try again.",
        isSubmittingStatus: false,
        successMessage: null
      }));
    }
  }

  async function handleRemoveApplication() {
    if (!state.application) {
      return;
    }

    const confirmed = window.confirm(
      "Remove this application? This also removes related reminders, contacts, documents, and activity history."
    );

    if (!confirmed) {
      return;
    }

    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      successMessage: null
    }));

    try {
      await removeApplication(state.application.id);
      onApplicationRemoved();
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Application could not be removed. Try again.",
        successMessage: null
      }));
    }
  }

  async function handleRemoveReminder(reminderId: string) {
    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isRemovingReminder: true,
      successMessage: null
    }));

    try {
      await removeReminder(reminderId);
      await loadDetailData();

      setState((currentState) => ({
        ...currentState,
        isRemovingReminder: false,
        successMessage: "Reminder removed."
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Reminder could not be removed. Try again.",
        isRemovingReminder: false,
        successMessage: null
      }));
    }
  }

  async function handleScheduleReminder(values: CreateReminderFormValues) {
    if (!applicationId) {
      return;
    }

    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isSubmittingReminder: true,
      successMessage: null
    }));

    try {
      await scheduleReminder(applicationId, values);
      await loadDetailData();

      setState((currentState) => ({
        ...currentState,
        isSubmittingReminder: false,
        successMessage: "Reminder scheduled."
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Reminder could not be scheduled. Check the form and try again.",
        isSubmittingReminder: false,
        successMessage: null
      }));
    }
  }

  async function handleCompleteReminder(reminderId: string) {
    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isCompletingReminder: true,
      successMessage: null
    }));

    try {
      await completeReminder(reminderId);
      await loadDetailData();

      setState((currentState) => ({
        ...currentState,
        isCompletingReminder: false,
        successMessage: "Reminder completed."
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Reminder could not be completed. Try again.",
        isCompletingReminder: false,
        successMessage: null
      }));
    }
  }

  async function handleAddContact(values: CreateContactFormValues) {
    if (!applicationId) {
      return;
    }

    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isAddingContact: true,
      successMessage: null
    }));

    try {
      await addContact(applicationId, values);
      await loadDetailData();

      setState((currentState) => ({
        ...currentState,
        isAddingContact: false,
        successMessage: "Contact added."
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Contact could not be added. Check the form and try again.",
        isAddingContact: false,
        successMessage: null
      }));
    }
  }

  async function handleRemoveContact(contactId: string) {
    if (!applicationId) {
      return;
    }

    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isRemovingContact: true,
      successMessage: null
    }));

    try {
      await removeContact(applicationId, contactId);
      await loadDetailData();

      setState((currentState) => ({
        ...currentState,
        isRemovingContact: false,
        successMessage: "Contact removed."
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Contact could not be removed. Try again.",
        isRemovingContact: false,
        successMessage: null
      }));
    }
  }

  async function handleRemoveDocument(documentId: string) {
    if (!applicationId) {
      return;
    }

    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isRemovingDocument: true,
      successMessage: null
    }));

    try {
      await removeDocument(applicationId, documentId);
      await loadDetailData();

      setState((currentState) => ({
        ...currentState,
        isRemovingDocument: false,
        successMessage: "Document metadata removed."
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Document metadata could not be removed. Try again.",
        isRemovingDocument: false,
        successMessage: null
      }));
    }
  }

  async function handleAddDocument(values: CreateDocumentFormValues) {
    if (!applicationId) {
      return;
    }

    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isAddingDocument: true,
      successMessage: null
    }));

    try {
      await addDocument(applicationId, values);
      await loadDetailData();

      setState((currentState) => ({
        ...currentState,
        isAddingDocument: false,
        successMessage: "Document metadata added."
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Document metadata could not be added. Check the form and try again.",
        isAddingDocument: false,
        successMessage: null
      }));
    }
  }

  return {
    state,
    actions: {
      handleAddContact,
      handleAddDocument,
      handleCompleteReminder,
      handleRemoveApplication,
      handleRemoveContact,
      handleRemoveDocument,
      handleRemoveReminder,
      handleScheduleReminder,
      handleStatusUpdate
    }
  };
}