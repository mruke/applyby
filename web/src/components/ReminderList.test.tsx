import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, test, vi } from "vitest";

import type { ReminderResponse } from "../types/application";
import { ReminderList } from "./ReminderList";

/**
 * buildReminder
 *
 * Creates a reminder response for reminder list component tests.
 */
function buildReminder(completed = false): ReminderResponse {
  return {
    id: "rem-001",
    application_id: "app-001",
    title: "Follow up with recruiter",
    due_at: "2026-05-14T09:00:00Z",
    completed
  };
}

describe("ReminderList", () => {
  test("renders an empty reminder state", () => {
    render(<ReminderList isCompleting={false} onComplete={vi.fn()} reminders={[]} />);

    expect(screen.getByText("No reminders scheduled for this application.")).toBeInTheDocument();
  });

  test("renders reminders", () => {
    render(<ReminderList isCompleting={false} onComplete={vi.fn()} reminders={[buildReminder()]} />);

    expect(screen.getByText("Follow up with recruiter")).toBeInTheDocument();
    expect(screen.getByText(/Incomplete/)).toBeInTheDocument();
  });

  test("completes an incomplete reminder", async () => {
    const onComplete = vi.fn().mockResolvedValue(undefined);

    render(<ReminderList isCompleting={false} onComplete={onComplete} reminders={[buildReminder()]} />);

    fireEvent.click(screen.getByRole("button", { name: "Complete" }));

    await waitFor(() => {
      expect(onComplete).toHaveBeenCalledWith("rem-001");
    });
  });

  test("does not render complete button for completed reminders", () => {
    render(<ReminderList isCompleting={false} onComplete={vi.fn()} reminders={[buildReminder(true)]} />);

    expect(screen.queryByRole("button", { name: "Complete" })).not.toBeInTheDocument();
  });
});