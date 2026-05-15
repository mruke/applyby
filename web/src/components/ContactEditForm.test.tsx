import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, test, vi } from "vitest";

import type { ContactResponse } from "../types/application";
import { ContactEditForm } from "./ContactEditForm";

// -----------------------------------------------------------------------------
// buildContact
//
// Creates a contact response for contact edit form tests.
// -----------------------------------------------------------------------------
function buildContact(): ContactResponse {
  return {
    id: "contact-001",
    application_id: "app-001",
    name: "Sam Recruiter",
    email: "sam@example.com",
    role: "Recruiter"
  };
}

describe("ContactEditForm", () => {
  test("renders current contact values", () => {
    render(<ContactEditForm contact={buildContact()} isSubmitting={false} onSubmit={vi.fn()} />);

    expect(screen.getByLabelText("Name")).toHaveValue("Sam Recruiter");
    expect(screen.getByLabelText("Email")).toHaveValue("sam@example.com");
    expect(screen.getByLabelText("Role")).toHaveValue("Recruiter");
  });

  test("shows validation when the name is missing", async () => {
    const onSubmit = vi.fn();

    render(<ContactEditForm contact={buildContact()} isSubmitting={false} onSubmit={onSubmit} />);

    fireEvent.change(screen.getByLabelText("Name"), {
      target: { value: "" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Save contact" }));

    expect(await screen.findByRole("alert")).toHaveTextContent("Contact name is required.");
    expect(onSubmit).not.toHaveBeenCalled();
  });

  test("submits edited values", async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined);

    render(<ContactEditForm contact={buildContact()} isSubmitting={false} onSubmit={onSubmit} />);

    fireEvent.change(screen.getByLabelText("Name"), {
      target: { value: "Sam Hiring" }
    });

    fireEvent.change(screen.getByLabelText("Role"), {
      target: { value: "Hiring Manager" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Save contact" }));

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        name: "Sam Hiring",
        email: "sam@example.com",
        role: "Hiring Manager"
      });
    });
  });

  test("disables submit while submitting", () => {
    render(<ContactEditForm contact={buildContact()} isSubmitting={true} onSubmit={vi.fn()} />);

    expect(screen.getByRole("button", { name: "Saving..." })).toBeDisabled();
  });
});
