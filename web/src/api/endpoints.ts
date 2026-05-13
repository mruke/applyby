/**
 * endpoints
 *
 * Defines backend API paths in one location so page and component code does
 * not hard-code route strings throughout the frontend.
 */
export const endpoints = {
  applications: "/applications",
  applicationSearch: "/applications/search",
  applicationStatus: (applicationId: string) => `/applications/${applicationId}/status`,
  applicationActivity: (applicationId: string) => `/applications/${applicationId}/activity`,
  applicationReminders: (applicationId: string) => `/applications/${applicationId}/reminders`,
  applicationContacts: (applicationId: string) => `/applications/${applicationId}/contacts`,
  applicationDocuments: (applicationId: string) => `/applications/${applicationId}/documents`,
  completeReminder: (reminderId: string) => `/reminders/${reminderId}/complete`
} as const;