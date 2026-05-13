import { render, screen } from "@testing-library/react";
import { describe, expect, test } from "vitest";

import type { DocumentResponse } from "../types/application";
import { DocumentList } from "./DocumentList";

/**
 * buildDocument
 *
 * Creates a document metadata response for document list tests.
 */
function buildDocument(): DocumentResponse {
  return {
    id: "doc-001",
    application_id: "app-001",
    name: "Backend Resume",
    kind: "resume",
    path: "documents/backend-resume.pdf"
  };
}

describe("DocumentList", () => {
  test("renders an empty document state", () => {
    render(<DocumentList documents={[]} />);

    expect(screen.getByText("No document metadata added yet.")).toBeInTheDocument();
  });

  test("renders documents", () => {
    render(<DocumentList documents={[buildDocument()]} />);

    expect(screen.getByText("Backend Resume")).toBeInTheDocument();
    expect(screen.getByText(/resume/)).toBeInTheDocument();
    expect(screen.getByText(/documents\/backend-resume.pdf/)).toBeInTheDocument();
  });
});