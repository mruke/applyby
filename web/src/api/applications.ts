import { apiClient } from "./client";
import { endpoints } from "./endpoints";
import type { ApplicationsResponse } from "../types/application";

/**
 * getApplications
 *
 * Loads the current application collection from the backend API.
 * Page components should call this API boundary instead of using fetch directly.
 */
export async function getApplications(): Promise<ApplicationsResponse> {
  return apiClient.request<ApplicationsResponse>(endpoints.applications);
}