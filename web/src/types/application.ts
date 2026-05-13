/**
 * ApplicationStatus
 *
 * Represents the application lifecycle statuses returned by the backend API.
 */
export type ApplicationStatus =
  | "draft"
  | "interested"
  | "applied"
  | "interviewing"
  | "offer"
  | "rejected"
  | "withdrawn"
  | "archived";

/**
 * ApplicationResponse
 *
 * Represents one application returned by the backend API.
 */
export type ApplicationResponse = {
  id: string;
  title: string;
  company_name: string;
  company_website: string;
  status: ApplicationStatus;
  source: string;
  notes: string;
  created_at: string;
  applied_at?: string;
};

/**
 * ApplicationsResponse
 *
 * Represents a collection of applications returned by the backend API.
 */
export type ApplicationsResponse = {
  applications: ApplicationResponse[];
};

/**
 * CreateApplicationFormValues
 *
 * Represents user-entered values from the create application form.
 */
export type CreateApplicationFormValues = {
  title: string;
  companyName: string;
  companyWebsite: string;
  status: ApplicationStatus;
  source: string;
  notes: string;
};

/**
 * CreateApplicationRequest
 *
 * Represents the backend request body for creating an application.
 */
export type CreateApplicationRequest = {
  id: string;
  title: string;
  company_name: string;
  company_website: string;
  status: ApplicationStatus;
  source: string;
  notes: string;
  created_at: string;
};

/**
 * ReminderResponse
 *
 * Represents one reminder returned by the backend API.
 */
export type ReminderResponse = {
  id: string;
  application_id: string;
  title: string;
  due_at: string;
  completed: boolean;
};

/**
 * RemindersResponse
 *
 * Represents a collection of reminders returned by the backend API.
 */
export type RemindersResponse = {
  reminders: ReminderResponse[];
};

/**
 * ActivityEventResponse
 *
 * Represents one activity timeline event returned by the backend API.
 */
export type ActivityEventResponse = {
  application_id: string;
  type: string;
  occurred_at: string;
  description: string;
};

/**
 * ActivityEventsResponse
 *
 * Represents a collection of activity events returned by the backend API.
 */
export type ActivityEventsResponse = {
  activity_events: ActivityEventResponse[];
};

/**
 * ContactResponse
 *
 * Represents one application contact returned by the backend API.
 */
export type ContactResponse = {
  id: string;
  application_id: string;
  name: string;
  email: string;
  role: string;
};

/**
 * ContactsResponse
 *
 * Represents a collection of contacts returned by the backend API.
 */
export type ContactsResponse = {
  contacts: ContactResponse[];
};

/**
 * DocumentResponse
 *
 * Represents one document metadata record returned by the backend API.
 */
export type DocumentResponse = {
  id: string;
  application_id: string;
  name: string;
  kind: string;
  path: string;
};

/**
 * DocumentsResponse
 *
 * Represents a collection of document metadata records returned by the backend API.
 */
export type DocumentsResponse = {
  documents: DocumentResponse[];
};