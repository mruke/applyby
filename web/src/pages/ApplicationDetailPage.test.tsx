import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { beforeEach, describe, expect, test, vi } from "vitest";

import { getActivityEvents } from "../api/activity";
import { getApplicationById, updateApplicationStatus } from "../api/applications";
import { addContact, getContacts, removeContact } from "../api/contacts";
import { addDocument, getDocuments } from "../api/documents";
import { completeReminder, getReminders, scheduleReminder } from "../api/reminders";
import type {
  ActivityEventsResponse,
  ApplicationResponse,
  ContactsResponse,
  DocumentsResponse,
  RemindersResponse
} from "../types/application";
import { ApplicationDetailPage } from "./ApplicationDetailPage";

vi.mock("../api/activity", () => ({
  getActivityEvents: vi.fn()
}));

vi.mock("../api/applications", () => ({
  getApplicationById: vi.fn(),
  updateApplicationStatus: vi.fn()
}));

vi.mock("../api/contacts", () => ({
  addContact: vi.fn(),
  getContacts: vi.fn(),
  removeContact: vi.fn()
}));

vi.mock("../api/documents", () => ({
  addDocument: vi.fn(),
  getDocuments: vi.fn()
}));

vi.mock("../api/reminders", () => ({
  completeReminder: vi.fn(),
  getReminders: vi.fn(),
  scheduleReminder: vi.fn()
}));

/**
 * mockedGetActivityEvents
 *
 * Provides typed access to the mocked activity API function.
 */
const mockedGetActivityEvents = vi.mocked(getActivityEvents);

/**
 * mockedGetApplicationById
 *
 * Provides typed access to the mocked application detail API function.
 */
const mockedGetApplicationById = vi.mocked(getApplicationById);

/**
 * mockedUpdateApplicationStatus
 *
 * Provides typed access to the mocked status update API function.
 */
const mockedUpdateApplicationStatus = vi.mocked(updateApplicationStatus);

/**
 * mockedAddContact
 *
 * Provides typed access to the mocked add contact API function.
 */
const mockedAddContact = vi.mocked(addContact);

/**
 * mockedGetContacts
 *
 * Provides typed access to the mocked contact list API function.
 */
const mockedGetContacts = vi.mocked(getContacts);

// -----------------------------------------------------------------------------
// mockedRemoveContact
//
// Provides typed access to the mocked remove contact API function.
// -----------------------------------------------------------------------------
const mockedRemoveContact = vi.mocked(removeContact);

/**
 * mockedAddDocument
 *
 * Provides typed access to the mocked add document API function.
 */
const mockedAddDocument = vi.mocked(addDocument);

/**
 * mockedGetDocuments
 *
 * Provides typed access to the mocked document list API function.
 */
const mockedGetDocuments = vi.mocked(getDocuments);

/**
 * mockedCompleteReminder
 *
 * Provides typed access to the mocked complete reminder API function.
 */
const mockedCompleteReminder = vi.mocked(completeReminder);

/**
 * mockedGetReminders
 *
 * Provides typed access to the mocked reminder list API function.
 */
const mockedGetReminders = vi.mocked(getReminders);

/**
 * mockedScheduleReminder
 *
 * Provides typed access to the mocked schedule reminder API function.
 */
const mockedScheduleReminder = vi.mocked(scheduleReminder);

/**
 * buildApplication
 *
 * Creates an application response for detail page tests.
 */
function buildApplication(status: ApplicationResponse["status"] = "applied"): ApplicationResponse {
  return {
    id: "app-001",
    title: "Backend Developer",
    company_name: "Example Studio",
    company_website: "https://example.com",
    status,
    source: "Company site",
    notes: "Applied with backend resume.",
    created_at: "2026-05-10T08:00:00Z"
  };
}

/**
 * buildRemindersResponse
 *
 * Creates a reminders response for detail page tests.
 */
function buildRemindersResponse(completed = false): RemindersResponse {
  return {
    reminders: [
      {
        id: "rem-001",
        application_id: "app-001",
        title: "Follow up with recruiter",
        due_at: "2026-05-14T09:00:00Z",
        completed
      }
    ]
  };
}

/**
 * buildActivityResponse
 *
 * Creates an activity response for detail page tests.
 */
function buildActivityResponse(): ActivityEventsResponse {
  return {
    activity_events: [
      {
        application_id: "app-001",
        type: "status_changed",
        occurred_at: "2026-05-14T09:00:00Z",
        description: "Status changed from applied to interviewing."
      }
    ]
  };
}

/**
 * buildContactsResponse
 *
 * Creates a contacts response for detail page tests.
 */
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

/**
 * buildDocumentsResponse
 *
 * Creates a documents response for detail page tests.
 */
function buildDocumentsResponse(): DocumentsResponse {
  return {
    documents: [
      {
        id: "doc-001",
        application_id: "app-001",
        name: "Backend Resume",
        kind: "resume",
        path: "documents/backend-resume.pdf"
      }
    ]
  };
}

/**
 * renderDetailPage
 *
 * Renders the application detail page at a route with an application id.
 */
function renderDetailPage(route = "/applications/app-001") {
  return render(
    <MemoryRouter initialEntries={[route]}>
      <Routes>
        <Route path="/applications/:applicationId" element={<ApplicationDetailPage />} />
      </Routes>
    </MemoryRouter>
  );
}

/**
 * mockSuccessfulDetailLoad
 *
 * Configures API mocks for a successful detail page load.
 */
function mockSuccessfulDetailLoad() {
  mockedGetApplicationById.mockResolvedValue(buildApplication());
  mockedGetReminders.mockResolvedValue(buildRemindersResponse());
  mockedGetActivityEvents.mockResolvedValue(buildActivityResponse());
  mockedGetContacts.mockResolvedValue(buildContactsResponse());
  mockedGetDocuments.mockResolvedValue(buildDocumentsResponse());
}

beforeEach(() => {
  mockedAddContact.mockReset();
  mockedAddDocument.mockReset();
  mockedCompleteReminder.mockReset();
  mockedGetActivityEvents.mockReset();
  mockedGetApplicationById.mockReset();
  mockedGetContacts.mockReset();
  mockedRemoveContact.mockReset();
  mockedGetDocuments.mockReset();
  mockedGetReminders.mockReset();
  mockedScheduleReminder.mockReset();
  mockedUpdateApplicationStatus.mockReset();
});

describe("ApplicationDetailPage", () => {
  test("shows a loading state while the application is loading", () => {
    mockedGetApplicationById.mockReturnValue(new Promise(() => {}));

    renderDetailPage();

    expect(screen.getByText("Loading application...")).toBeInTheDocument();
  });

  test("renders application details, reminders, activity, contacts, and documents", async () => {
    mockSuccessfulDetailLoad();

    renderDetailPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Backend Developer" })).toBeInTheDocument();
    expect(screen.getByText("Example Studio")).toBeInTheDocument();
    expect(screen.getAllByText("Applied").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Applied with backend resume.").length).toBeGreaterThan(0);
    expect(screen.getByText("Follow up with recruiter")).toBeInTheDocument();
    expect(screen.getByText("Status changed from applied to interviewing.")).toBeInTheDocument();
    expect(screen.getByText("Sam Recruiter")).toBeInTheDocument();
    expect(screen.getByText("Backend Resume")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Edit application" })).toHaveAttribute(
      "href",
      "/applications/app-001/edit"
    );
  });

  test("shows a not found state when the application does not exist", async () => {
    mockedGetApplicationById.mockResolvedValue(null);

    renderDetailPage();

    expect(await screen.findByRole("heading", { level: 2, name: "Application not found" })).toBeInTheDocument();
  });

  test("shows an error state when loading fails", async () => {
    mockedGetApplicationById.mockRejectedValue(new Error("network failed"));

    renderDetailPage();

    expect(
      await screen.findByRole("heading", { level: 2, name: "Application could not be loaded" })
    ).toBeInTheDocument();
  });

  test("updates application status", async () => {
    mockSuccessfulDetailLoad();
    mockedUpdateApplicationStatus.mockResolvedValue(buildApplication("interviewing"));

    renderDetailPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Backend Developer" })).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText("Status"), {
      target: { value: "interviewing" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Update status" }));

    await waitFor(() => {
      expect(mockedUpdateApplicationStatus).toHaveBeenCalledWith("app-001", "interviewing");
    });

    expect(await screen.findByRole("status")).toHaveTextContent("Status updated.");
  });

  test("shows an error when status update fails", async () => {
    mockSuccessfulDetailLoad();
    mockedUpdateApplicationStatus.mockRejectedValue(new Error("invalid transition"));

    renderDetailPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Backend Developer" })).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText("Status"), {
      target: { value: "interviewing" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Update status" }));

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "Status could not be updated. Check the selected status and try again."
    );
  });

  test("schedules a reminder", async () => {
    mockSuccessfulDetailLoad();
    mockedScheduleReminder.mockResolvedValue(buildRemindersResponse().reminders[0]);

    renderDetailPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Backend Developer" })).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText("Reminder title"), {
      target: { value: "Follow up with recruiter" }
    });

    fireEvent.change(screen.getByLabelText("Due date and time"), {
      target: { value: "2026-05-14T09:00" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Schedule reminder" }));

    await waitFor(() => {
      expect(mockedScheduleReminder).toHaveBeenCalledWith("app-001", {
        title: "Follow up with recruiter",
        dueAt: "2026-05-14T09:00"
      });
    });

    expect(await screen.findByRole("status")).toHaveTextContent("Reminder scheduled.");
  });

  test("completes a reminder", async () => {
    mockSuccessfulDetailLoad();
    mockedCompleteReminder.mockResolvedValue(buildRemindersResponse(true).reminders[0]);

    renderDetailPage();

    expect(await screen.findByText("Follow up with recruiter")).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "Complete" }));

    await waitFor(() => {
      expect(mockedCompleteReminder).toHaveBeenCalledWith("rem-001");
    });

    expect(await screen.findByRole("status")).toHaveTextContent("Reminder completed.");
  });

  test("adds a contact", async () => {
    mockSuccessfulDetailLoad();
    mockedAddContact.mockResolvedValue(buildContactsResponse().contacts[0]);

    renderDetailPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Backend Developer" })).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText("Name"), {
      target: { value: "Sam Recruiter" }
    });

    fireEvent.change(screen.getByLabelText("Email"), {
      target: { value: "sam@example.com" }
    });

    fireEvent.change(screen.getByLabelText("Role"), {
      target: { value: "Recruiter" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Add contact" }));

    await waitFor(() => {
      expect(mockedAddContact).toHaveBeenCalledWith("app-001", {
        name: "Sam Recruiter",
        email: "sam@example.com",
        role: "Recruiter"
      });
    });

    expect(await screen.findByRole("status")).toHaveTextContent("Contact added.");
  });

  test("removes a contact", async () => {
    mockSuccessfulDetailLoad();
    mockedRemoveContact.mockResolvedValue(undefined);

    renderDetailPage();

    expect(await screen.findByText("Sam Recruiter")).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "Remove" }));

    await waitFor(() => {
      expect(mockedRemoveContact).toHaveBeenCalledWith("app-001", "contact-001");
    });

    expect(await screen.findByRole("status")).toHaveTextContent("Contact removed.");
  });

  test("adds document metadata", async () => {
    mockSuccessfulDetailLoad();
    mockedAddDocument.mockResolvedValue(buildDocumentsResponse().documents[0]);

    renderDetailPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Backend Developer" })).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText("Document name"), {
      target: { value: "Backend Resume" }
    });

    fireEvent.change(screen.getByLabelText("Kind"), {
      target: { value: "resume" }
    });

    fireEvent.change(screen.getByLabelText("Path or reference"), {
      target: { value: "documents/backend-resume.pdf" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Add document" }));

    await waitFor(() => {
      expect(mockedAddDocument).toHaveBeenCalledWith("app-001", {
        name: "Backend Resume",
        kind: "resume",
        path: "documents/backend-resume.pdf"
      });
    });

    expect(await screen.findByRole("status")).toHaveTextContent("Document metadata added.");
  });
});
