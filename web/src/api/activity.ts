import { apiClient } from "./client";
import { endpoints } from "./endpoints";
import type { ActivityEventsResponse } from "../types/application";

/**
 * getActivityEvents
 *
 * Loads activity timeline events for one application through the backend API.
 */
export async function getActivityEvents(applicationId: string): Promise<ActivityEventsResponse> {
  return apiClient.request<ActivityEventsResponse>(endpoints.applicationActivity(applicationId));
}