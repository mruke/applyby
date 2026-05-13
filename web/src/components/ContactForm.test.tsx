import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, test, vi } from "vitest";

import { ContactForm } from "./ContactForm";

describe("ContactForm", () => {
  test("shows validation when the name is missing", async () => {
    const onSubmit = vi.fn();

    render(<ContactForm isSubmitting={false} onSubmit={onSubmit} />);

    fireEvent.click(screen.getByRole("button", { name: "Add contact" }));

    expect(await screen.findByRole("alert")).toHaveTextContent("Contact name is required.");
    expect(onSubmit).not.toHaveBeenCalled();
  });

  test("submits entered values", async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined);

    render(<ContactForm isSubmitting={false} onSubmit={onSubmit} />);

    fireEvent.change(screen.getByLabelText("Name"), {
      target: { value: "Sam Recruiter" }
    });

    fireEvent.change(screen.getByLabelText("Email"), {
      target: { value: "sam@example.com" }
    });

    fireEvent.change(screen.getByLabelText("Role"), {
      target: { value: "Recruiter" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Add contact" }));

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        name: "Sam Recruiter",
        email: "sam@example.com",
        role: "Recruiter"
      });
    });
  });

  test("disables submit while submitting", () => {
    render(<ContactForm isSubmitting={true} onSubmit={vi.fn()} />);

    expect(screen.getByRole("button", { name: "Adding..." })).toBeDisabled();
  });
});