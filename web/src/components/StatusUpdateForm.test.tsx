import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, test, vi } from "vitest";

import { StatusUpdateForm } from "./StatusUpdateForm";

describe("StatusUpdateForm", () => {
  test("disables submit when the selected status is unchanged", () => {
    const onSubmit = vi.fn();

    render(<StatusUpdateForm currentStatus="applied" isSubmitting={false} onSubmit={onSubmit} />);

    expect(screen.getByRole("button", { name: "Update status" })).toBeDisabled();
  });

  test("submits a changed status", async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined);

    render(<StatusUpdateForm currentStatus="applied" isSubmitting={false} onSubmit={onSubmit} />);

    fireEvent.change(screen.getByLabelText("Status"), {
      target: { value: "interviewing" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Update status" }));

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith("interviewing");
    });
  });

  test("shows pending submit text while submitting", () => {
    const onSubmit = vi.fn();

    render(<StatusUpdateForm currentStatus="applied" isSubmitting={true} onSubmit={onSubmit} />);

    expect(screen.getByRole("button", { name: "Updating..." })).toBeDisabled();
  });
});