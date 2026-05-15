import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import type { ComponentProps } from "react";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, test, vi } from "vitest";

import type { ReminderResponse } from "../types/application";
import { ReminderSection } from "./ReminderSection";

function buildReminder(): ReminderResponse {
  return {
    id: "rem-001",
    application_id: "app-001",
    title: "Send follow-up",
    due_at: "2026-05-20T09:30:00Z",
    completed: false
  };
}

function renderReminderSection(overrides: Partial<ComponentProps<typeof ReminderSection>> = {}) {
  const props: ComponentProps<typeof ReminderSection> = {
    applicationId: "app-001",
    reminders: [buildReminder()],
    errorMessage: null,
    isCompleting: false,
    isRemoving: false,
    isSubmitting: false,
    onAdd: vi.fn().mockResolvedValue(undefined),
    onComplete: vi.fn().mockResolvedValue(undefined),
    onRemove: vi.fn().mockResolvedValue(undefined),
    ...overrides
  };

  render(
    <MemoryRouter>
      <ReminderSection {...props} />
    </MemoryRouter>
  );

  return props;
}

describe("ReminderSection", () => {
  test("renders the add form and reminder list", () => {
    renderReminderSection();

    expect(screen.getByRole("heading", { level: 2, name: "Schedule reminder" })).toBeInTheDocument();
    expect(screen.getByText("Send follow-up")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Edit" })).toHaveAttribute(
      "href",
      "/applications/app-001/reminders/rem-001/edit"
    );
  });

  test("shows a reminder load error while keeping the add form available", () => {
    renderReminderSection({ errorMessage: "Reminders could not be loaded." });

    expect(screen.getByRole("heading", { level: 2, name: "Schedule reminder" })).toBeInTheDocument();
    expect(screen.getByRole("alert")).toHaveTextContent("Reminders could not be loaded.");
    expect(screen.queryByText("Send follow-up")).not.toBeInTheDocument();
  });

  test("delegates remove actions to the parent page", async () => {
    const onRemove = vi.fn().mockResolvedValue(undefined);

    renderReminderSection({ onRemove });

    fireEvent.click(screen.getByRole("button", { name: "Remove" }));

    await waitFor(() => {
      expect(onRemove).toHaveBeenCalledWith("rem-001");
    });
  });
});