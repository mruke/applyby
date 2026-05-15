import { apiClient } from "./client";
import { endpoints } from "./endpoints";
import type {
  CreateReminderFormValues,
  CreateReminderRequest,
  ReminderResponse,
  RemindersResponse,
  UpdateReminderFormValues,
  UpdateReminderRequest
} from "../types/application";
import { createClientId } from "../utils/clientIds";

// -----------------------------------------------------------------------------
// buildCreateReminderRequest
//
// Converts form values into the backend schedule-reminder request shape.
// -----------------------------------------------------------------------------
function buildCreateReminderRequest(values: CreateReminderFormValues): CreateReminderRequest {
  return {
    id: createClientId("reminder"),
    title: values.title.trim(),
    due_at: values.dueAt
  };
}

// -----------------------------------------------------------------------------
// buildUpdateReminderRequest
//
// Converts edit form values into the backend update-reminder request shape.
// -----------------------------------------------------------------------------
function buildUpdateReminderRequest(values: UpdateReminderFormValues): UpdateReminderRequest {
  return {
    title: values.title.trim(),
    due_at: values.dueAt
  };
}

// -----------------------------------------------------------------------------
// getReminders
//
// Loads reminders for one application through the backend API.
// -----------------------------------------------------------------------------
export async function getReminders(applicationId: string): Promise<RemindersResponse> {
  return apiClient.request<RemindersResponse>(endpoints.applicationReminders(applicationId));
}

// -----------------------------------------------------------------------------
// scheduleReminder
//
// Schedules a reminder for one application through the backend API.
// -----------------------------------------------------------------------------
export async function scheduleReminder(
  applicationId: string,
  values: CreateReminderFormValues
): Promise<ReminderResponse> {
  return apiClient.request<ReminderResponse>(endpoints.applicationReminders(applicationId), {
    method: "POST",
    body: buildCreateReminderRequest(values)
  });
}

// -----------------------------------------------------------------------------
// completeReminder
//
// Marks one reminder complete through the backend API.
// -----------------------------------------------------------------------------
export async function completeReminder(reminderId: string): Promise<ReminderResponse> {
  return apiClient.request<ReminderResponse>(endpoints.completeReminder(reminderId), {
    method: "PATCH"
  });
}

// -----------------------------------------------------------------------------
// updateReminder
//
// Updates one reminder through the backend API.
// -----------------------------------------------------------------------------
export async function updateReminder(
  reminderId: string,
  values: UpdateReminderFormValues
): Promise<ReminderResponse> {
  return apiClient.request<ReminderResponse>(endpoints.reminder(reminderId), {
    method: "PATCH",
    body: buildUpdateReminderRequest(values)
  });
}

// -----------------------------------------------------------------------------
// removeReminder
//
// Removes one reminder through the backend API.
// -----------------------------------------------------------------------------
export async function removeReminder(reminderId: string): Promise<void> {
  return apiClient.request<void>(endpoints.reminder(reminderId), {
    method: "DELETE"
  });
}