import { apiClient } from "./client";
import { endpoints } from "./endpoints";
import type {
  CreateReminderFormValues,
  CreateReminderRequest,
  ReminderResponse,
  RemindersResponse
} from "../types/application";

/**
 * createReminderId
 *
 * Creates a client-side reminder identity for the backend schedule reminder workflow.
 * The fallback exists for browsers or tests that do not expose crypto.randomUUID.
 */
function createReminderId(): string {
  if (globalThis.crypto && "randomUUID" in globalThis.crypto) {
    return globalThis.crypto.randomUUID();
  }

  return `reminder-${Date.now()}`;
}

/**
 * buildCreateReminderRequest
 *
 * Converts form values into the backend schedule-reminder request shape.
 */
function buildCreateReminderRequest(values: CreateReminderFormValues): CreateReminderRequest {
  return {
    id: createReminderId(),
    title: values.title.trim(),
    due_at: new Date(values.dueAt).toISOString()
  };
}

/**
 * getReminders
 *
 * Loads reminders for one application through the backend API.
 */
export async function getReminders(applicationId: string): Promise<RemindersResponse> {
  return apiClient.request<RemindersResponse>(endpoints.applicationReminders(applicationId));
}

/**
 * scheduleReminder
 *
 * Schedules a reminder for one application through the backend API.
 */
export async function scheduleReminder(
  applicationId: string,
  values: CreateReminderFormValues
): Promise<ReminderResponse> {
  return apiClient.request<ReminderResponse>(endpoints.applicationReminders(applicationId), {
    method: "POST",
    body: buildCreateReminderRequest(values)
  });
}

/**
 * completeReminder
 *
 * Marks one reminder as complete through the backend API.
 */
export async function completeReminder(reminderId: string): Promise<ReminderResponse> {
  return apiClient.request<ReminderResponse>(endpoints.completeReminder(reminderId), {
    method: "PATCH"
  });
}