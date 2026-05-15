import { apiClient } from "./client";
import { endpoints } from "./endpoints";
import type {
  CreateDocumentFormValues,
  CreateDocumentRequest,
  DocumentResponse,
  DocumentsResponse,
  UpdateDocumentFormValues,
  UpdateDocumentRequest
} from "../types/application";
import { createClientId } from "../utils/clientIds";

/**
 * buildCreateDocumentRequest
 *
 * Converts form values into the backend add-document request shape.
 */
function buildCreateDocumentRequest(values: CreateDocumentFormValues): CreateDocumentRequest {
  return {
    id: createClientId("document"),
    name: values.name.trim(),
    kind: values.kind.trim(),
    path: values.path.trim()
  };
}

// -----------------------------------------------------------------------------
// buildUpdateDocumentRequest
//
// Converts edit form values into the backend update-document request shape.
// -----------------------------------------------------------------------------
function buildUpdateDocumentRequest(values: UpdateDocumentFormValues): UpdateDocumentRequest {
  return {
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

// -----------------------------------------------------------------------------
// updateDocument
//
// Updates one application document metadata record through the backend API.
// -----------------------------------------------------------------------------
export async function updateDocument(
  applicationId: string,
  documentId: string,
  values: UpdateDocumentFormValues
): Promise<DocumentResponse> {
  return apiClient.request<DocumentResponse>(endpoints.applicationDocument(applicationId, documentId), {
    method: "PATCH",
    body: buildUpdateDocumentRequest(values)
  });
}

// -----------------------------------------------------------------------------
// removeDocument
//
// Removes one application document metadata record through the backend API.
// -----------------------------------------------------------------------------
export async function removeDocument(applicationId: string, documentId: string): Promise<void> {
  return apiClient.request<void>(endpoints.applicationDocument(applicationId, documentId), {
    method: "DELETE"
  });
}