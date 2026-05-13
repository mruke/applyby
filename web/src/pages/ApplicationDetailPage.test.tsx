import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { beforeEach, describe, expect, test, vi } from "vitest";

import { getApplicationById, updateApplicationStatus } from "../api/applications";
import type { ApplicationResponse } from "../types/application";
import { ApplicationDetailPage } from "./ApplicationDetailPage";

vi.mock("../api/applications", () => ({
  getApplicationById: vi.fn(),
  updateApplicationStatus: vi.fn()
}));

/**
 * mockedGetApplicationById
 *
 * Provides typed access to the mocked application detail API function.
 */
const mockedGetApplicationById = vi.mocked(getApplicationById);

/**
 * mockedUpdateApplicationStatus
 *
 * Provides typed access to the mocked status update API function.
 */
const mockedUpdateApplicationStatus = vi.mocked(updateApplicationStatus);

/**
 * buildApplication
 *
 * Creates an application response for detail page tests.
 */
function buildApplication(status: ApplicationResponse["status"] = "applied"): ApplicationResponse {
  return {
    id: "app-001",
    title: "Backend Developer",
    company_name: "Example Studio",
    company_website: "https://example.com",
    status,
    source: "Company site",
    notes: "Applied with backend resume.",
    created_at: "2026-05-10T08:00:00Z"
  };
}

/**
 * renderDetailPage
 *
 * Renders the application detail page at a route with an application id.
 */
function renderDetailPage(route = "/applications/app-001") {
  return render(
    <MemoryRouter initialEntries={[route]}>
      <Routes>
        <Route path="/applications/:applicationId" element={<ApplicationDetailPage />} />
      </Routes>
    </MemoryRouter>
  );
}

beforeEach(() => {
  mockedGetApplicationById.mockReset();
  mockedUpdateApplicationStatus.mockReset();
});

describe("ApplicationDetailPage", () => {
  test("shows a loading state while the application is loading", () => {
    mockedGetApplicationById.mockReturnValue(new Promise(() => {}));

    renderDetailPage();

    expect(screen.getByText("Loading application...")).toBeInTheDocument();
  });

  test("renders application details", async () => {
    mockedGetApplicationById.mockResolvedValue(buildApplication());

    renderDetailPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Backend Developer" })).toBeInTheDocument();
    expect(screen.getByText("Example Studio")).toBeInTheDocument();
    expect(screen.getAllByText("Applied").length).toBeGreaterThan(0);
    expect(screen.getByText("Applied with backend resume.")).toBeInTheDocument();
  });

  test("shows a not found state when the application does not exist", async () => {
    mockedGetApplicationById.mockResolvedValue(null);

    renderDetailPage();

    expect(await screen.findByRole("heading", { level: 2, name: "Application not found" })).toBeInTheDocument();
  });

  test("shows an error state when loading fails", async () => {
    mockedGetApplicationById.mockRejectedValue(new Error("network failed"));

    renderDetailPage();

    expect(
      await screen.findByRole("heading", { level: 2, name: "Application could not be loaded" })
    ).toBeInTheDocument();
  });

  test("updates application status", async () => {
    mockedGetApplicationById.mockResolvedValue(buildApplication());
    mockedUpdateApplicationStatus.mockResolvedValue(buildApplication("interviewing"));

    renderDetailPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Backend Developer" })).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText("Status"), {
      target: { value: "interviewing" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Update status" }));

    await waitFor(() => {
      expect(mockedUpdateApplicationStatus).toHaveBeenCalledWith("app-001", "interviewing");
    });

    expect(await screen.findByRole("status")).toHaveTextContent("Status updated.");
  });

  test("shows an error when status update fails", async () => {
    mockedGetApplicationById.mockResolvedValue(buildApplication());
    mockedUpdateApplicationStatus.mockRejectedValue(new Error("invalid transition"));

    renderDetailPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Backend Developer" })).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText("Status"), {
      target: { value: "interviewing" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Update status" }));

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "Status could not be updated. Check the selected status and try again."
    );
  });
});