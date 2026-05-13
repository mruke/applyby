import { apiClient } from "./client";
import { endpoints } from "./endpoints";
import type {
  CreateDocumentFormValues,
  CreateDocumentRequest,
  DocumentResponse,
  DocumentsResponse
} from "../types/application";

/**
 * createDocumentId
 *
 * Creates a client-side document identity for the backend add document metadata workflow.
 * The fallback exists for browsers or tests that do not expose crypto.randomUUID.
 */
function createDocumentId(): string {
  if (globalThis.crypto && "randomUUID" in globalThis.crypto) {
    return globalThis.crypto.randomUUID();
  }

  return `document-${Date.now()}`;
}

/**
 * buildCreateDocumentRequest
 *
 * Converts form values into the backend add-document request shape.
 */
function buildCreateDocumentRequest(values: CreateDocumentFormValues): CreateDocumentRequest {
  return {
    id: createDocumentId(),
    name: values.name.trim(),
    kind: values.kind.trim(),
    path: values.path.trim()
  };
}

/**
 * getDocuments
 *
 * Loads document metadata records for one application through the backend API.
 */
export async function getDocuments(applicationId: string): Promise<DocumentsResponse> {
  return apiClient.request<DocumentsResponse>(endpoints.applicationDocuments(applicationId));
}

/**
 * addDocument
 *
 * Adds document metadata to one application through the backend API.
 */
export async function addDocument(
  applicationId: string,
  values: CreateDocumentFormValues
): Promise<DocumentResponse> {
  return apiClient.request<DocumentResponse>(endpoints.applicationDocuments(applicationId), {
    method: "POST",
    body: buildCreateDocumentRequest(values)
  });
}