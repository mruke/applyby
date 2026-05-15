import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { beforeEach, describe, expect, test, vi } from "vitest";

import { getApplicationById } from "../api/applications";
import { getDocuments, updateDocument } from "../api/documents";
import type { ApplicationResponse, DocumentsResponse } from "../types/application";
import { DocumentEditPage } from "./DocumentEditPage";

vi.mock("../api/applications", () => ({
  getApplicationById: vi.fn()
}));

vi.mock("../api/documents", () => ({
  getDocuments: vi.fn(),
  updateDocument: vi.fn()
}));

const mockedGetApplicationById = vi.mocked(getApplicationById);
const mockedGetDocuments = vi.mocked(getDocuments);
const mockedUpdateDocument = vi.mocked(updateDocument);

function buildApplication(): ApplicationResponse {
  return {
    id: "app-001",
    title: "Backend Developer",
    company_name: "Example Studio",
    company_website: "https://example.com",
    status: "applied",
    source: "Company site",
    notes: "Applied with backend resume.",
    created_at: "2026-05-10T08:00:00Z"
  };
}

function buildDocumentsResponse(): DocumentsResponse {
  return {
    documents: [
      {
        id: "doc-001",
        application_id: "app-001",
        name: "Backend Resume",
        kind: "resume",
        path: "documents/backend-resume.pdf"
      }
    ]
  };
}

function renderDocumentEditPage(route = "/applications/app-001/documents/doc-001/edit") {
  return render(
    <MemoryRouter initialEntries={[route]}>
      <Routes>
        <Route path="/applications/:applicationId/documents/:documentId/edit" element={<DocumentEditPage />} />
        <Route path="/applications/:applicationId" element={<h1>Application detail route</h1>} />
      </Routes>
    </MemoryRouter>
  );
}

beforeEach(() => {
  mockedGetApplicationById.mockReset();
  mockedGetDocuments.mockReset();
  mockedUpdateDocument.mockReset();
});

describe("DocumentEditPage", () => {
  test("shows a loading state while the document is loading", () => {
    mockedGetApplicationById.mockReturnValue(new Promise(() => {}));
    mockedGetDocuments.mockReturnValue(new Promise(() => {}));

    renderDocumentEditPage();

    expect(screen.getByText("Loading document metadata...")).toBeInTheDocument();
  });

  test("renders the document edit form", async () => {
    mockedGetApplicationById.mockResolvedValue(buildApplication());
    mockedGetDocuments.mockResolvedValue(buildDocumentsResponse());

    renderDocumentEditPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Edit document metadata" })).toBeInTheDocument();
    expect(screen.getByLabelText("Document name")).toHaveValue("Backend Resume");
    expect(screen.getByRole("link", { name: "Cancel" })).toHaveAttribute("href", "/applications/app-001");
  });

  test("shows a not found state when the document does not exist", async () => {
    mockedGetApplicationById.mockResolvedValue(buildApplication());
    mockedGetDocuments.mockResolvedValue({ documents: [] });

    renderDocumentEditPage();

    expect(await screen.findByRole("heading", { level: 2, name: "Document metadata not found" })).toBeInTheDocument();
  });

  test("shows an error state when loading fails", async () => {
    mockedGetApplicationById.mockRejectedValue(new Error("network failed"));
    mockedGetDocuments.mockResolvedValue(buildDocumentsResponse());

    renderDocumentEditPage();

    expect(
      await screen.findByRole("heading", { level: 2, name: "Document metadata could not be loaded" })
    ).toBeInTheDocument();
  });

  test("updates document metadata and returns to the detail route", async () => {
    mockedGetApplicationById.mockResolvedValue(buildApplication());
    mockedGetDocuments.mockResolvedValue(buildDocumentsResponse());
    mockedUpdateDocument.mockResolvedValue({
      ...buildDocumentsResponse().documents[0],
      name: "Backend Resume v2"
    });

    renderDocumentEditPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Edit document metadata" })).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText("Document name"), {
      target: { value: "Backend Resume v2" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Save document" }));

    await waitFor(() => {
      expect(mockedUpdateDocument).toHaveBeenCalledWith("app-001", "doc-001", {
        name: "Backend Resume v2",
        kind: "resume",
        path: "documents/backend-resume.pdf"
      });
    });

    expect(await screen.findByRole("heading", { level: 1, name: "Application detail route" })).toBeInTheDocument();
  });

  test("shows an error when update fails", async () => {
    mockedGetApplicationById.mockResolvedValue(buildApplication());
    mockedGetDocuments.mockResolvedValue(buildDocumentsResponse());
    mockedUpdateDocument.mockRejectedValue(new Error("update failed"));

    renderDocumentEditPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Edit document metadata" })).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "Save document" }));

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "Document metadata could not be updated. Check the form and try again."
    );
  });
});