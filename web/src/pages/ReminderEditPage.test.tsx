import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { beforeEach, describe, expect, test, vi } from "vitest";

import { getApplicationById } from "../api/applications";
import { getReminders, updateReminder } from "../api/reminders";
import type { ApplicationResponse, RemindersResponse } from "../types/application";
import { ReminderEditPage } from "./ReminderEditPage";

vi.mock("../api/applications", () => ({
  getApplicationById: vi.fn()
}));

vi.mock("../api/reminders", () => ({
  getReminders: vi.fn(),
  updateReminder: vi.fn()
}));

const mockedGetApplicationById = vi.mocked(getApplicationById);
const mockedGetReminders = vi.mocked(getReminders);
const mockedUpdateReminder = vi.mocked(updateReminder);

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

function buildRemindersResponse(): RemindersResponse {
  return {
    reminders: [
      {
        id: "rem-001",
        application_id: "app-001",
        title: "Send follow-up",
        due_at: "2026-05-20T09:30:00Z",
        completed: false
      }
    ]
  };
}

function renderReminderEditPage(route = "/applications/app-001/reminders/rem-001/edit") {
  return render(
    <MemoryRouter initialEntries={[route]}>
      <Routes>
        <Route path="/applications/:applicationId/reminders/:reminderId/edit" element={<ReminderEditPage />} />
        <Route path="/applications/:applicationId" element={<h1>Application detail route</h1>} />
      </Routes>
    </MemoryRouter>
  );
}

beforeEach(() => {
  mockedGetApplicationById.mockReset();
  mockedGetReminders.mockReset();
  mockedUpdateReminder.mockReset();
});

describe("ReminderEditPage", () => {
  test("shows a loading state while the reminder is loading", () => {
    mockedGetApplicationById.mockReturnValue(new Promise(() => {}));
    mockedGetReminders.mockReturnValue(new Promise(() => {}));

    renderReminderEditPage();

    expect(screen.getByText("Loading reminder...")).toBeInTheDocument();
  });

  test("renders the reminder edit form", async () => {
    mockedGetApplicationById.mockResolvedValue(buildApplication());
    mockedGetReminders.mockResolvedValue(buildRemindersResponse());

    renderReminderEditPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Edit reminder" })).toBeInTheDocument();
    expect(screen.getByLabelText("Reminder title")).toHaveValue("Send follow-up");
    expect(screen.getByRole("link", { name: "Cancel" })).toHaveAttribute("href", "/applications/app-001");
  });

  test("shows a not found state when the reminder does not exist", async () => {
    mockedGetApplicationById.mockResolvedValue(buildApplication());
    mockedGetReminders.mockResolvedValue({ reminders: [] });

    renderReminderEditPage();

    expect(await screen.findByRole("heading", { level: 2, name: "Reminder not found" })).toBeInTheDocument();
  });

  test("shows an error state when loading fails", async () => {
    mockedGetApplicationById.mockRejectedValue(new Error("network failed"));
    mockedGetReminders.mockResolvedValue(buildRemindersResponse());

    renderReminderEditPage();

    expect(await screen.findByRole("heading", { level: 2, name: "Reminder could not be loaded" })).toBeInTheDocument();
  });

  test("updates reminder and returns to the detail route", async () => {
    mockedGetApplicationById.mockResolvedValue(buildApplication());
    mockedGetReminders.mockResolvedValue(buildRemindersResponse());
    mockedUpdateReminder.mockResolvedValue({
      ...buildRemindersResponse().reminders[0],
      title: "Send updated follow-up"
    });

    renderReminderEditPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Edit reminder" })).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText("Reminder title"), {
      target: { value: "Send updated follow-up" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Save reminder" }));

    await waitFor(() => {
      expect(mockedUpdateReminder).toHaveBeenCalledWith(
        "rem-001",
        expect.objectContaining({
          title: "Send updated follow-up"
        })
      );
    });

    expect(await screen.findByRole("heading", { level: 1, name: "Application detail route" })).toBeInTheDocument();
  });

  test("shows an error when update fails", async () => {
    mockedGetApplicationById.mockResolvedValue(buildApplication());
    mockedGetReminders.mockResolvedValue(buildRemindersResponse());
    mockedUpdateReminder.mockRejectedValue(new Error("update failed"));

    renderReminderEditPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Edit reminder" })).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "Save reminder" }));

    expect(await screen.findByRole("alert")).toHaveTextContent("Reminder could not be updated. Check the form and try again.");
  });
});