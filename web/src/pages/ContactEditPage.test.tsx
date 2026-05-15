import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { beforeEach, describe, expect, test, vi } from "vitest";

import { getApplicationById } from "../api/applications";
import { getContacts, updateContact } from "../api/contacts";
import type { ApplicationResponse, ContactsResponse } from "../types/application";
import { ContactEditPage } from "./ContactEditPage";

vi.mock("../api/applications", () => ({
  getApplicationById: vi.fn()
}));

vi.mock("../api/contacts", () => ({
  getContacts: vi.fn(),
  updateContact: vi.fn()
}));

const mockedGetApplicationById = vi.mocked(getApplicationById);
const mockedGetContacts = vi.mocked(getContacts);
const mockedUpdateContact = vi.mocked(updateContact);

// -----------------------------------------------------------------------------
// buildApplication
//
// Creates an application response for contact edit page tests.
// -----------------------------------------------------------------------------
function buildApplication(): ApplicationResponse {
  return {
    id: "app-001",
    title: "Backend Developer",
    company_name: "Example Studio",
    company_website: "https://example.com",
    status: "applied",
    source: "Company site",
    notes: "Applied with backend resume.",
    created_at: "2026-05-10T08:00:00Z"
  };
}

// -----------------------------------------------------------------------------
// buildContactsResponse
//
// Creates a contacts response for contact edit page tests.
// -----------------------------------------------------------------------------
function buildContactsResponse(): ContactsResponse {
  return {
    contacts: [
      {
        id: "contact-001",
        application_id: "app-001",
        name: "Sam Recruiter",
        email: "sam@example.com",
        role: "Recruiter"
      }
    ]
  };
}

// -----------------------------------------------------------------------------
// renderContactEditPage
//
// Renders the contact edit page at the contact edit route.
// -----------------------------------------------------------------------------
function renderContactEditPage(route = "/applications/app-001/contacts/contact-001/edit") {
  return render(
    <MemoryRouter initialEntries={[route]}>
      <Routes>
        <Route path="/applications/:applicationId/contacts/:contactId/edit" element={<ContactEditPage />} />
        <Route path="/applications/:applicationId" element={<h1>Application detail route</h1>} />
      </Routes>
    </MemoryRouter>
  );
}

beforeEach(() => {
  mockedGetApplicationById.mockReset();
  mockedGetContacts.mockReset();
  mockedUpdateContact.mockReset();
});

describe("ContactEditPage", () => {
  test("shows a loading state while the contact is loading", () => {
    mockedGetApplicationById.mockReturnValue(new Promise(() => {}));
    mockedGetContacts.mockReturnValue(new Promise(() => {}));

    renderContactEditPage();

    expect(screen.getByText("Loading contact...")).toBeInTheDocument();
  });

  test("renders the contact edit form", async () => {
    mockedGetApplicationById.mockResolvedValue(buildApplication());
    mockedGetContacts.mockResolvedValue(buildContactsResponse());

    renderContactEditPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Edit contact" })).toBeInTheDocument();
    expect(screen.getByLabelText("Name")).toHaveValue("Sam Recruiter");
    expect(screen.getByRole("link", { name: "Cancel" })).toHaveAttribute("href", "/applications/app-001");
  });

  test("shows a not found state when the contact does not exist", async () => {
    mockedGetApplicationById.mockResolvedValue(buildApplication());
    mockedGetContacts.mockResolvedValue({ contacts: [] });

    renderContactEditPage();

    expect(await screen.findByRole("heading", { level: 2, name: "Contact not found" })).toBeInTheDocument();
  });

  test("shows an error state when loading fails", async () => {
    mockedGetApplicationById.mockRejectedValue(new Error("network failed"));
    mockedGetContacts.mockResolvedValue(buildContactsResponse());

    renderContactEditPage();

    expect(await screen.findByRole("heading", { level: 2, name: "Contact could not be loaded" })).toBeInTheDocument();
  });

  test("updates contact and returns to the detail route", async () => {
    mockedGetApplicationById.mockResolvedValue(buildApplication());
    mockedGetContacts.mockResolvedValue(buildContactsResponse());
    mockedUpdateContact.mockResolvedValue({
      ...buildContactsResponse().contacts[0],
      name: "Sam Hiring"
    });

    renderContactEditPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Edit contact" })).toBeInTheDocument();

    const nameInput = screen.getByLabelText("Name");

    fireEvent.change(nameInput, {
      target: { value: "Sam Hiring" }
    });

    await waitFor(() => {
      expect(nameInput).toHaveValue("Sam Hiring");
    });

    fireEvent.click(screen.getByRole("button", { name: "Save contact" }));

    await waitFor(() => {
      expect(mockedUpdateContact).toHaveBeenCalledWith("app-001", "contact-001", {
        name: "Sam Hiring",
        email: "sam@example.com",
        role: "Recruiter"
      });
    });

    expect(await screen.findByRole("heading", { level: 1, name: "Application detail route" })).toBeInTheDocument();
  });

  test("shows an error when update fails", async () => {
    mockedGetApplicationById.mockResolvedValue(buildApplication());
    mockedGetContacts.mockResolvedValue(buildContactsResponse());
    mockedUpdateContact.mockRejectedValue(new Error("update failed"));

    renderContactEditPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Edit contact" })).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "Save contact" }));

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "Contact could not be updated. Check the form and try again."
    );
  });
});
