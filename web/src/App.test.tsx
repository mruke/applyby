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
  updateApplicationStatus: vi.fn()
}));

vi.mock("./api/contacts", () => ({
  addContact: vi.fn(),
  getContacts: vi.fn().mockResolvedValue({ contacts: [] })
}));

vi.mock("./api/documents", () => ({
  addDocument: vi.fn(),
  getDocuments: vi.fn().mockResolvedValue({ documents: [] })
}));

vi.mock("./api/reminders", () => ({
  completeReminder: vi.fn(),
  getReminders: vi.fn().mockResolvedValue({ reminders: [] }),
  scheduleReminder: vi.fn()
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

  test("renders the not found route", () => {
    renderRoute("/missing");

    expect(screen.getByRole("heading", { level: 1, name: "Page not found" })).toBeInTheDocument();
  });
});