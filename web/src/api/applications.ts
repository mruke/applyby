import { apiClient } from "./client";
import { endpoints } from "./endpoints";
import type {
  ApplicationResponse,
  ApplicationsResponse,
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
 * getApplications
 *
 * Loads the current application collection from the backend API.
 * Page components should call this API boundary instead of using fetch directly.
 */
export async function getApplications(): Promise<ApplicationsResponse> {
  return apiClient.request<ApplicationsResponse>(endpoints.applications);
}

/**
 * getApplicationById
 *
 * Loads one application by using the existing applications collection endpoint.
 * This avoids adding a backend detail endpoint before the backend contract exists.
 */
export async function getApplicationById(applicationId: string): Promise<ApplicationResponse | null> {
  const response = await getApplications();

  return response.applications.find((application) => application.id === applicationId) ?? null;
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