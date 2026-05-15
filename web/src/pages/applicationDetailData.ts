import { getActivityEvents } from "../api/activity";
import { getApplicationById } from "../api/applications";
import { getContacts } from "../api/contacts";
import { getDocuments } from "../api/documents";
import { getReminders } from "../api/reminders";
import type {
  ActivityEventResponse,
  ApplicationResponse,
  ContactResponse,
  DocumentResponse,
  ReminderResponse
} from "../types/application";

/**
 * SectionErrorMessages
 *
 * Represents independently recoverable load failures for detail-page sections.
 */
export type SectionErrorMessages = {
  activity: string | null;
  contacts: string | null;
  documents: string | null;
  reminders: string | null;
};

/**
 * ApplicationDetailData
 *
 * Represents the data needed to render the application detail page.
 */
export type ApplicationDetailData = {
  activityEvents: ActivityEventResponse[];
  application: ApplicationResponse | null;
  contacts: ContactResponse[];
  documents: DocumentResponse[];
  reminders: ReminderResponse[];
  sectionErrors: SectionErrorMessages;
};

/**
 * emptySectionErrors
 *
 * Provides the default section-level error state.
 */
export function emptySectionErrors(): SectionErrorMessages {
  return {
    activity: null,
    contacts: null,
    documents: null,
    reminders: null
  };
}

/**
 * emptyApplicationDetailData
 *
 * Provides empty detail data for a missing application.
 */
function emptyApplicationDetailData(): ApplicationDetailData {
  return {
    activityEvents: [],
    application: null,
    contacts: [],
    documents: [],
    reminders: [],
    sectionErrors: emptySectionErrors()
  };
}

/**
 * sectionErrorMessage
 *
 * Provides readable section-level load error messages.
 */
function sectionErrorMessage(section: keyof SectionErrorMessages): string {
  switch (section) {
    case "activity":
      return "Activity could not be loaded.";
    case "contacts":
      return "Contacts could not be loaded.";
    case "documents":
      return "Documents could not be loaded.";
    case "reminders":
      return "Reminders could not be loaded.";
    default:
      return "Section data could not be loaded.";
  }
}

/**
 * resultValueOrDefault
 *
 * Returns the fulfilled result value or a fallback when a section request fails.
 */
function resultValueOrDefault<TValue>(result: PromiseSettledResult<TValue>, fallback: TValue): TValue {
  return result.status === "fulfilled" ? result.value : fallback;
}

/**
 * sectionErrorFromResult
 *
 * Returns a section-specific error message when a section request fails.
 */
function sectionErrorFromResult<TValue>(
  section: keyof SectionErrorMessages,
  result: PromiseSettledResult<TValue>
): string | null {
  return result.status === "rejected" ? sectionErrorMessage(section) : null;
}

/**
 * fetchApplicationDetailData
 *
 * Loads one application and its related detail resources.
 * The core application must load successfully. Related sections are allowed
 * to fail independently so one supporting feature does not block the page.
 */
export async function fetchApplicationDetailData(applicationId: string): Promise<ApplicationDetailData> {
  const application = await getApplicationById(applicationId);

  if (!application) {
    return emptyApplicationDetailData();
  }

  const [remindersResult, activityResult, contactsResult, documentsResult] = await Promise.allSettled([
    getReminders(applicationId),
    getActivityEvents(applicationId),
    getContacts(applicationId),
    getDocuments(applicationId)
  ]);

  const remindersResponse = resultValueOrDefault(remindersResult, { reminders: [] });
  const activityResponse = resultValueOrDefault(activityResult, { activity_events: [] });
  const contactsResponse = resultValueOrDefault(contactsResult, { contacts: [] });
  const documentsResponse = resultValueOrDefault(documentsResult, { documents: [] });

  return {
    activityEvents: activityResponse.activity_events,
    application,
    contacts: contactsResponse.contacts,
    documents: documentsResponse.documents,
    reminders: remindersResponse.reminders,
    sectionErrors: {
      activity: sectionErrorFromResult("activity", activityResult),
      contacts: sectionErrorFromResult("contacts", contactsResult),
      documents: sectionErrorFromResult("documents", documentsResult),
      reminders: sectionErrorFromResult("reminders", remindersResult)
    }
  };
}