import { apiClient } from "./client";
import { endpoints } from "./endpoints";
import type { ContactResponse, ContactsResponse, CreateContactFormValues, CreateContactRequest } from "../types/application";

/**
 * createContactId
 *
 * Creates a client-side contact identity for the backend add contact workflow.
 * The fallback exists for browsers or tests that do not expose crypto.randomUUID.
 */
function createContactId(): string {
  if (globalThis.crypto && "randomUUID" in globalThis.crypto) {
    return globalThis.crypto.randomUUID();
  }

  return `contact-${Date.now()}`;
}

/**
 * buildCreateContactRequest
 *
 * Converts form values into the backend add-contact request shape.
 */
function buildCreateContactRequest(values: CreateContactFormValues): CreateContactRequest {
  return {
    id: createContactId(),
    name: values.name.trim(),
    email: values.email.trim(),
    role: values.role.trim()
  };
}

/**
 * getContacts
 *
 * Loads contacts for one application through the backend API.
 */
export async function getContacts(applicationId: string): Promise<ContactsResponse> {
  return apiClient.request<ContactsResponse>(endpoints.applicationContacts(applicationId));
}

/**
 * addContact
 *
 * Adds a contact to one application through the backend API.
 */
export async function addContact(applicationId: string, values: CreateContactFormValues): Promise<ContactResponse> {
  return apiClient.request<ContactResponse>(endpoints.applicationContacts(applicationId), {
    method: "POST",
    body: buildCreateContactRequest(values)
  });
}