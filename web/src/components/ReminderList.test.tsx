import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, test, vi } from "vitest";

import type { ReminderResponse } from "../types/application";
import { ReminderList } from "./ReminderList";

function buildReminder(completed = false): ReminderResponse {
  return {
    id: "rem-001",
    application_id: "app-001",
    title: "Send follow-up",
    due_at: "2026-05-20T09:30:00Z",
    completed
  };
}

function renderReminderList(reminders = [buildReminder()]) {
  const onComplete = vi.fn().mockResolvedValue(undefined);
  const onRemove = vi.fn().mockResolvedValue(undefined);

  render(
    <MemoryRouter>
      <ReminderList
        applicationId="app-001"
        reminders={reminders}
        isCompleting={false}
        isRemoving={false}
        onComplete={onComplete}
        onRemove={onRemove}
      />
    </MemoryRouter>
  );

  return { onComplete, onRemove };
}

describe("ReminderList", () => {
  test("renders an empty reminder state", () => {
    render(
      <MemoryRouter>
        <ReminderList
          applicationId="app-001"
          reminders={[]}
          isCompleting={false}
          onComplete={vi.fn().mockResolvedValue(undefined)}
        />
      </MemoryRouter>
    );

    expect(screen.getByText("No reminders scheduled yet.")).toBeInTheDocument();
  });

  test("renders reminders with edit and complete actions", () => {
    renderReminderList();

    expect(screen.getByText("Send follow-up")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Edit" })).toHaveAttribute(
      "href",
      "/applications/app-001/reminders/rem-001/edit"
    );
    expect(screen.getByRole("button", { name: "Complete" })).toBeInTheDocument();
  });

  test("does not show complete action for completed reminders", () => {
    renderReminderList([buildReminder(true)]);

    expect(screen.queryByRole("button", { name: "Complete" })).not.toBeInTheDocument();
    expect(screen.getByText(/Completed/)).toBeInTheDocument();
  });

  test("delegates complete actions", async () => {
    const { onComplete } = renderReminderList();

    fireEvent.click(screen.getByRole("button", { name: "Complete" }));

    await waitFor(() => {
      expect(onComplete).toHaveBeenCalledWith("rem-001");
    });
  });

  test("delegates remove actions", async () => {
    const { onRemove } = renderReminderList();

    fireEvent.click(screen.getByRole("button", { name: "Remove" }));

    await waitFor(() => {
      expect(onRemove).toHaveBeenCalledWith("rem-001");
    });
  });
});