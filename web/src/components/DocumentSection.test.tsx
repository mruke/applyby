import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import type { ComponentProps } from "react";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, test, vi } from "vitest";

import type { DocumentResponse } from "../types/application";
import { DocumentSection } from "./DocumentSection";

function buildDocument(): DocumentResponse {
  return {
    id: "doc-001",
    application_id: "app-001",
    name: "Backend Resume",
    kind: "resume",
    path: "documents/backend-resume.pdf"
  };
}

function renderDocumentSection(overrides: Partial<ComponentProps<typeof DocumentSection>> = {}) {
  const props: ComponentProps<typeof DocumentSection> = {
    applicationId: "app-001",
    documents: [buildDocument()],
    errorMessage: null,
    isAdding: false,
    isRemoving: false,
    onAdd: vi.fn().mockResolvedValue(undefined),
    onRemove: vi.fn().mockResolvedValue(undefined),
    ...overrides
  };

  render(
    <MemoryRouter>
      <DocumentSection {...props} />
    </MemoryRouter>
  );

  return props;
}

describe("DocumentSection", () => {
  test("renders the add form and document list", () => {
    renderDocumentSection();

    expect(screen.getByRole("heading", { level: 2, name: "Add document metadata" })).toBeInTheDocument();
    expect(screen.getByText("Backend Resume")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Edit" })).toHaveAttribute(
      "href",
      "/applications/app-001/documents/doc-001/edit"
    );
  });

  test("shows a document load error while keeping the add form available", () => {
    renderDocumentSection({ errorMessage: "Documents could not be loaded." });

    expect(screen.getByRole("heading", { level: 2, name: "Add document metadata" })).toBeInTheDocument();
    expect(screen.getByRole("alert")).toHaveTextContent("Documents could not be loaded.");
    expect(screen.queryByText("Backend Resume")).not.toBeInTheDocument();
  });

  test("delegates remove actions to the parent page", async () => {
    const onRemove = vi.fn().mockResolvedValue(undefined);

    renderDocumentSection({ onRemove });

    fireEvent.click(screen.getByRole("button", { name: "Remove" }));

    await waitFor(() => {
      expect(onRemove).toHaveBeenCalledWith("doc-001");
    });
  });
});