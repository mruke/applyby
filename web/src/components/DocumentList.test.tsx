import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, test, vi } from "vitest";

import type { DocumentResponse } from "../types/application";
import { DocumentList } from "./DocumentList";

function buildDocument(): DocumentResponse {
  return {
    id: "doc-001",
    application_id: "app-001",
    name: "Backend Resume",
    kind: "resume",
    path: "documents/backend-resume.pdf"
  };
}

function renderDocumentList(onRemove = vi.fn().mockResolvedValue(undefined)) {
  render(
    <MemoryRouter>
      <DocumentList applicationId="app-001" documents={[buildDocument()]} onRemove={onRemove} />
    </MemoryRouter>
  );

  return onRemove;
}

describe("DocumentList", () => {
  test("renders an empty document state", () => {
    render(
      <MemoryRouter>
        <DocumentList applicationId="app-001" documents={[]} />
      </MemoryRouter>
    );

    expect(screen.getByText("No document metadata added yet.")).toBeInTheDocument();
  });

  test("renders documents with edit action", () => {
    renderDocumentList();

    expect(screen.getByText("Backend Resume")).toBeInTheDocument();
    expect(screen.getByText(/resume/)).toBeInTheDocument();
    expect(screen.getByText(/documents\/backend-resume.pdf/)).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Edit" })).toHaveAttribute(
      "href",
      "/applications/app-001/documents/doc-001/edit"
    );
  });

  test("delegates remove actions", async () => {
    const onRemove = renderDocumentList();

    fireEvent.click(screen.getByRole("button", { name: "Remove" }));

    await waitFor(() => {
      expect(onRemove).toHaveBeenCalledWith("doc-001");
    });
  });
});