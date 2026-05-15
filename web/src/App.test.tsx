import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, test, vi } from "vitest";

import App from "./App";

vi.mock("./api/activity", () => ({
  getActivityEvents: vi.fn().mockResolvedValue({ activity_events: [] })
}));

vi.mock("./api/applications", () => ({
  createApplication: vi.fn(),
  getApplicationById: vi.fn().mockResolvedValue({
    id: "app-001",
    title: "Backend Developer",
    company_name: "Example Studio",
    company_website: "https://example.com",
    status: "applied",
    source: "Company site",
    notes: "Applied with backend resume.",
    created_at: "2026-05-10T08:00:00Z"
  }),
  getApplications: vi.fn(() => new Promise(() => {})),
  updateApplicationDetails: vi.fn(),
  updateApplicationStatus: vi.fn()
}));

vi.mock("./api/contacts", () => ({
  addContact: vi.fn(),
  getContacts: vi.fn().mockResolvedValue({
    contacts: [
      {
        id: "contact-001",
        application_id: "app-001",
        name: "Sam Recruiter",
        email: "sam@example.com",
        role: "Recruiter"
      }
    ]
  }),
  removeContact: vi.fn(),
  updateContact: vi.fn()
}));

vi.mock("./api/documents", () => ({
  addDocument: vi.fn(),
  getDocuments: vi.fn().mockResolvedValue({
    documents: [
      {
        id: "doc-001",
        application_id: "app-001",
        name: "Backend Resume",
        kind: "resume",
        path: "documents/backend-resume.pdf"
      }
    ]
  }),
  removeDocument: vi.fn(),
  updateDocument: vi.fn()
}));

vi.mock("./api/reminders", () => ({
  completeReminder: vi.fn(),
  getReminders: vi.fn().mockResolvedValue({
    reminders: [
      {
        id: "rem-001",
        application_id: "app-001",
        title: "Send follow-up",
        due_at: "2026-05-20T09:30:00Z",
        completed: false
      }
    ]
  }),
  removeReminder: vi.fn(),
  scheduleReminder: vi.fn(),
  updateReminder: vi.fn()
}));

/**
 * renderRoute
 *
 * Renders the application at a specific route for route-level component tests.
 * MemoryRouter avoids relying on the real browser URL during tests.
 */
function renderRoute(route: string) {
  return render(
    <MemoryRouter initialEntries={[route]}>
      <App />
    </MemoryRouter>
  );
}

describe("App", () => {
  test("renders the dashboard route", () => {
    renderRoute("/");

    expect(screen.getByRole("heading", { level: 1, name: "Dashboard" })).toBeInTheDocument();
  });

  test("renders the applications route", () => {
    renderRoute("/applications");

    expect(screen.getByRole("heading", { level: 1, name: "Applications" })).toBeInTheDocument();
  });

  test("renders the application detail route", async () => {
    renderRoute("/applications/app-001");

    expect(await screen.findByRole("heading", { level: 1, name: "Backend Developer" })).toBeInTheDocument();
  });

  test("renders the application edit route", async () => {
    renderRoute("/applications/app-001/edit");

    expect(await screen.findByRole("heading", { level: 1, name: "Edit application" })).toBeInTheDocument();
  });

  test("renders the contact edit route", async () => {
    renderRoute("/applications/app-001/contacts/contact-001/edit");

    expect(await screen.findByRole("heading", { level: 1, name: "Edit contact" })).toBeInTheDocument();
  });

  test("renders the document edit route", async () => {
    renderRoute("/applications/app-001/documents/doc-001/edit");

    expect(await screen.findByRole("heading", { level: 1, name: "Edit document metadata" })).toBeInTheDocument();
  });

  test("renders the reminder edit route", async () => {
    renderRoute("/applications/app-001/reminders/rem-001/edit");

    expect(await screen.findByRole("heading", { level: 1, name: "Edit reminder" })).toBeInTheDocument();
  });

  test("renders the not found route", () => {
    renderRoute("/missing");

    expect(screen.getByRole("heading", { level: 1, name: "Page not found" })).toBeInTheDocument();
  });
});
