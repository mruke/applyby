import { ApiError, apiClient } from "./client";
import { endpoints } from "./endpoints";
import type {
  ApplicationResponse,
  ApplicationsResponse,
  ApplicationSearchCriteria,
  ApplicationStatus,
  CreateApplicationFormValues,
  CreateApplicationRequest
} from "../types/application";

/**
 * createApplicationId
 *
 * Creates a client-side application identity for the backend create workflow.
 * The fallback exists for browsers or tests that do not expose crypto.randomUUID.
 */
function createApplicationId(): string {
  if (globalThis.crypto && "randomUUID" in globalThis.crypto) {
    return globalThis.crypto.randomUUID();
  }

  return `app-${Date.now()}`;
}

/**
 * buildCreateApplicationRequest
 *
 * Converts form values into the backend create-application request shape.
 */
function buildCreateApplicationRequest(values: CreateApplicationFormValues): CreateApplicationRequest {
  return {
    id: createApplicationId(),
    title: values.title.trim(),
    company_name: values.companyName.trim(),
    company_website: values.companyWebsite.trim(),
    status: values.status,
    source: values.source.trim(),
    notes: values.notes.trim(),
    created_at: new Date().toISOString()
  };
}

/**
 * buildSearchPath
 *
 * Converts frontend search criteria into the backend query-parameter contract.
 */
function buildSearchPath(criteria: ApplicationSearchCriteria): string {
  const params = new URLSearchParams();

  if (criteria.companyName.trim() !== "") {
    params.set("company_name", criteria.companyName.trim());
  }

  if (criteria.source.trim() !== "") {
    params.set("source", criteria.source.trim());
  }

  if (criteria.text.trim() !== "") {
    params.set("text", criteria.text.trim());
  }

  for (const status of criteria.statuses) {
    params.append("status", status);
  }

  const queryString = params.toString();

  return queryString === "" ? endpoints.applicationSearch : `${endpoints.applicationSearch}?${queryString}`;
}

/**
 * getApplications
 *
 * Loads the current application collection from the backend API.
 * Page components should call this API boundary instead of using fetch directly.
 */
export async function getApplications(): Promise<ApplicationsResponse> {
  return apiClient.request<ApplicationsResponse>(endpoints.applications);
}

/**
 * searchApplications
 *
 * Searches applications using the backend search query-parameter contract.
 */
export async function searchApplications(criteria: ApplicationSearchCriteria): Promise<ApplicationsResponse> {
  return apiClient.request<ApplicationsResponse>(buildSearchPath(criteria));
}

/**
 * getApplicationById
 *
 * Loads one application through the backend detail endpoint.
 */
export async function getApplicationById(applicationId: string): Promise<ApplicationResponse | null> {
  try {
    return await apiClient.request<ApplicationResponse>(endpoints.applicationDetail(applicationId));
  } catch (error) {
    if (error instanceof ApiError && error.statusCode === 404) {
      return null;
    }

    throw error;
  }
}

/**
 * createApplication
 *
 * Creates a new application through the backend API.
 */
export async function createApplication(values: CreateApplicationFormValues): Promise<ApplicationResponse> {
  return apiClient.request<ApplicationResponse>(endpoints.applications, {
    method: "POST",
    body: buildCreateApplicationRequest(values)
  });
}

/**
 * updateApplicationStatus
 *
 * Updates one application's status through the backend API.
 */
export async function updateApplicationStatus(
  applicationId: string,
  status: ApplicationStatus
): Promise<ApplicationResponse> {
  return apiClient.request<ApplicationResponse>(endpoints.applicationStatus(applicationId), {
    method: "PATCH",
    body: { status }
  });
}