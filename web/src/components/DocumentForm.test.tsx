import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, test, vi } from "vitest";

import { DocumentForm } from "./DocumentForm";

describe("DocumentForm", () => {
  test("shows validation when the document name is missing", async () => {
    const onSubmit = vi.fn();

    render(<DocumentForm isSubmitting={false} onSubmit={onSubmit} />);

    fireEvent.change(screen.getByLabelText("Kind"), {
      target: { value: "resume" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Add document" }));

    expect(await screen.findByRole("alert")).toHaveTextContent("Document name is required.");
    expect(onSubmit).not.toHaveBeenCalled();
  });

  test("shows validation when the kind is missing", async () => {
    const onSubmit = vi.fn();

    render(<DocumentForm isSubmitting={false} onSubmit={onSubmit} />);

    fireEvent.change(screen.getByLabelText("Document name"), {
      target: { value: "Backend Resume" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Add document" }));

    expect(await screen.findByRole("alert")).toHaveTextContent("Document kind is required.");
    expect(onSubmit).not.toHaveBeenCalled();
  });

  test("submits entered values", async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined);

    render(<DocumentForm isSubmitting={false} onSubmit={onSubmit} />);

    fireEvent.change(screen.getByLabelText("Document name"), {
      target: { value: "Backend Resume" }
    });

    fireEvent.change(screen.getByLabelText("Kind"), {
      target: { value: "resume" }
    });

    fireEvent.change(screen.getByLabelText("Path or reference"), {
      target: { value: "documents/backend-resume.pdf" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Add document" }));

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        name: "Backend Resume",
        kind: "resume",
        path: "documents/backend-resume.pdf"
      });
    });
  });

  test("disables submit while submitting", () => {
    render(<DocumentForm isSubmitting={true} onSubmit={vi.fn()} />);

    expect(screen.getByRole("button", { name: "Adding..." })).toBeDisabled();
  });
});