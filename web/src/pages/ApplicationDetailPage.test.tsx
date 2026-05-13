import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { beforeEach, describe, expect, test, vi } from "vitest";

import { getActivityEvents } from "../api/activity";
import { getApplicationById, updateApplicationStatus } from "../api/applications";
import { completeReminder, getReminders, scheduleReminder } from "../api/reminders";
import type { ActivityEventsResponse, ApplicationResponse, RemindersResponse } from "../types/application";
import { ApplicationDetailPage } from "./ApplicationDetailPage";

vi.mock("../api/activity", () => ({
  getActivityEvents: vi.fn()
}));

vi.mock("../api/applications", () => ({
  getApplicationById: vi.fn(),
  updateApplicationStatus: vi.fn()
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
}

beforeEach(() => {
  mockedCompleteReminder.mockReset();
  mockedGetActivityEvents.mockReset();
  mockedGetApplicationById.mockReset();
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

  test("renders application details, reminders, and activity", async () => {
    mockSuccessfulDetailLoad();

    renderDetailPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Backend Developer" })).toBeInTheDocument();
    expect(screen.getByText("Example Studio")).toBeInTheDocument();
    expect(screen.getAllByText("Applied").length).toBeGreaterThan(0);
    expect(screen.getByText("Applied with backend resume.")).toBeInTheDocument();
    expect(screen.getByText("Follow up with recruiter")).toBeInTheDocument();
    expect(screen.getByText("Status changed from applied to interviewing.")).toBeInTheDocument();
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
});