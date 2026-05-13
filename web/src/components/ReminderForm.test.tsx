import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, test, vi } from "vitest";

import { ReminderForm } from "./ReminderForm";

/**
 * fillReminderForm
 *
 * Fills the required reminder form fields for component tests.
 */
function fillReminderForm() {
  fireEvent.change(screen.getByLabelText("Reminder title"), {
    target: { value: "Follow up with recruiter" }
  });

  fireEvent.change(screen.getByLabelText("Due date and time"), {
    target: { value: "2026-05-14T09:00" }
  });
}

describe("ReminderForm", () => {
  test("shows validation when the title is missing", async () => {
    const onSubmit = vi.fn();

    render(<ReminderForm isSubmitting={false} onSubmit={onSubmit} />);

    fireEvent.change(screen.getByLabelText("Due date and time"), {
      target: { value: "2026-05-14T09:00" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Schedule reminder" }));

    expect(await screen.findByRole("alert")).toHaveTextContent("Reminder title is required.");
    expect(onSubmit).not.toHaveBeenCalled();
  });

  test("shows validation when the due date is missing", async () => {
    const onSubmit = vi.fn();

    render(<ReminderForm isSubmitting={false} onSubmit={onSubmit} />);

    fireEvent.change(screen.getByLabelText("Reminder title"), {
      target: { value: "Follow up with recruiter" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Schedule reminder" }));

    expect(await screen.findByRole("alert")).toHaveTextContent("Reminder due date is required.");
    expect(onSubmit).not.toHaveBeenCalled();
  });

  test("submits entered values", async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined);

    render(<ReminderForm isSubmitting={false} onSubmit={onSubmit} />);

    fillReminderForm();
    fireEvent.click(screen.getByRole("button", { name: "Schedule reminder" }));

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        title: "Follow up with recruiter",
        dueAt: "2026-05-14T09:00"
      });
    });
  });

  test("disables the submit button while submitting", () => {
    const onSubmit = vi.fn();

    render(<ReminderForm isSubmitting={true} onSubmit={onSubmit} />);

    expect(screen.getByRole("button", { name: "Scheduling..." })).toBeDisabled();
  });
});