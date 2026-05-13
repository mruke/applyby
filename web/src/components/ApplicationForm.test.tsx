import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, test, vi } from "vitest";

import { ApplicationForm } from "./ApplicationForm";

/**
 * fillRequiredFields
 *
 * Fills the minimum required fields for create application form tests.
 */
function fillRequiredFields() {
  fireEvent.change(screen.getByLabelText("Application title"), {
    target: { value: "Backend Developer" }
  });

  fireEvent.change(screen.getByLabelText("Company name"), {
    target: { value: "Example Studio" }
  });
}

describe("ApplicationForm", () => {
  test("shows validation when the title is missing", async () => {
    const onSubmit = vi.fn();

    render(<ApplicationForm isSubmitting={false} onSubmit={onSubmit} />);

    fireEvent.change(screen.getByLabelText("Company name"), {
      target: { value: "Example Studio" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Add application" }));

    expect(await screen.findByRole("alert")).toHaveTextContent("Application title is required.");
    expect(onSubmit).not.toHaveBeenCalled();
  });

  test("shows validation when the company is missing", async () => {
    const onSubmit = vi.fn();

    render(<ApplicationForm isSubmitting={false} onSubmit={onSubmit} />);

    fireEvent.change(screen.getByLabelText("Application title"), {
      target: { value: "Backend Developer" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Add application" }));

    expect(await screen.findByRole("alert")).toHaveTextContent("Company name is required.");
    expect(onSubmit).not.toHaveBeenCalled();
  });

  test("submits entered values", async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined);

    render(<ApplicationForm isSubmitting={false} onSubmit={onSubmit} />);

    fillRequiredFields();

    fireEvent.change(screen.getByLabelText("Source"), {
      target: { value: "Company site" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Add application" }));

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          title: "Backend Developer",
          companyName: "Example Studio",
          source: "Company site",
          status: "applied"
        })
      );
    });
  });

  test("disables the submit button while submitting", () => {
    const onSubmit = vi.fn();

    render(<ApplicationForm isSubmitting={true} onSubmit={onSubmit} />);

    expect(screen.getByRole("button", { name: "Adding..." })).toBeDisabled();
  });
});