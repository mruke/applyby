import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { beforeEach, describe, expect, test, vi } from "vitest";

import { getApplicationById, updateApplicationDetails } from "../api/applications";
import type { ApplicationResponse } from "../types/application";
import { ApplicationEditPage } from "./ApplicationEditPage";

vi.mock("../api/applications", () => ({
  getApplicationById: vi.fn(),
  updateApplicationDetails: vi.fn()
}));

const mockedGetApplicationById = vi.mocked(getApplicationById);
const mockedUpdateApplicationDetails = vi.mocked(updateApplicationDetails);

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

function renderEditPage() {
  return render(
    <MemoryRouter initialEntries={["/applications/app-001/edit"]}>
      <Routes>
        <Route path="/applications/:applicationId/edit" element={<ApplicationEditPage />} />
        <Route path="/applications/:applicationId" element={<h1>Application detail route</h1>} />
      </Routes>
    </MemoryRouter>
  );
}

beforeEach(() => {
  mockedGetApplicationById.mockReset();
  mockedUpdateApplicationDetails.mockReset();
});

describe("ApplicationEditPage", () => {
  test("shows a loading state while the application is loading", () => {
    mockedGetApplicationById.mockReturnValue(new Promise(() => {}));

    renderEditPage();

    expect(screen.getByText("Loading application...")).toBeInTheDocument();
  });

  test("renders the application edit form", async () => {
    mockedGetApplicationById.mockResolvedValue(buildApplication());

    renderEditPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Edit application" })).toBeInTheDocument();
    expect(screen.getByLabelText("Application title")).toHaveValue("Backend Developer");
    expect(screen.getByLabelText("Company name")).toHaveValue("Example Studio");
    expect(screen.getByRole("link", { name: "Cancel" })).toHaveAttribute("href", "/applications/app-001");
  });

  test("shows a not found state when the application does not exist", async () => {
    mockedGetApplicationById.mockResolvedValue(null);

    renderEditPage();

    expect(await screen.findByRole("heading", { level: 2, name: "Application not found" })).toBeInTheDocument();
  });

  test("shows an error state when loading fails", async () => {
    mockedGetApplicationById.mockRejectedValue(new Error("network failed"));

    renderEditPage();

    expect(
      await screen.findByRole("heading", { level: 2, name: "Application could not be loaded" })
    ).toBeInTheDocument();
  });

  test("updates application details and returns to the detail route", async () => {
    mockedGetApplicationById.mockResolvedValue(buildApplication());
    mockedUpdateApplicationDetails.mockResolvedValue({
      ...buildApplication(),
      title: "Senior Backend Developer"
    });

    renderEditPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Edit application" })).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText("Application title"), {
      target: { value: "Senior Backend Developer" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Save application details" }));

    await waitFor(() => {
      expect(mockedUpdateApplicationDetails).toHaveBeenCalledWith("app-001", {
        title: "Senior Backend Developer",
        companyName: "Example Studio",
        companyWebsite: "https://example.com",
        source: "Company site",
        notes: "Applied with backend resume."
      });
    });

    expect(await screen.findByRole("heading", { level: 1, name: "Application detail route" })).toBeInTheDocument();
  });

  test("shows an error when update fails", async () => {
    mockedGetApplicationById.mockResolvedValue(buildApplication());
    mockedUpdateApplicationDetails.mockRejectedValue(new Error("update failed"));

    renderEditPage();

    expect(await screen.findByRole("heading", { level: 1, name: "Edit application" })).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "Save application details" }));

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "Application details could not be updated. Check the form and try again."
    );
  });
});