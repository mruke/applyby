import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { beforeEach, describe, expect, test, vi } from "vitest";

import { createApplication, getApplications } from "../api/applications";
import type { ApplicationsResponse } from "../types/application";
import { ApplicationsPage } from "./ApplicationsPage";

vi.mock("../api/applications", () => ({
  createApplication: vi.fn(),
  getApplications: vi.fn()
}));

/**
 * mockedCreateApplication
 *
 * Provides typed access to the mocked create application API function.
 */
const mockedCreateApplication = vi.mocked(createApplication);

/**
 * mockedGetApplications
 *
 * Provides typed access to the mocked applications API function.
 */
const mockedGetApplications = vi.mocked(getApplications);

/**
 * renderApplicationsPage
 *
 * Renders the applications page inside a router for component tests.
 */
function renderApplicationsPage() {
  return render(
    <MemoryRouter>
      <ApplicationsPage />
    </MemoryRouter>
  );
}

/**
 * buildApplicationsResponse
 *
 * Creates an applications API response for applications page tests.
 */
function buildApplicationsResponse(): ApplicationsResponse {
  return {
    applications: [
      {
        id: "app-001",
        title: "Backend Developer",
        company_name: "Example Studio",
        company_website: "https://example.com",
        status: "applied",
        source: "Company site",
        notes: "Applied with backend resume.",
        created_at: "2026-05-10T08:00:00Z"
      },
      {
        id: "app-002",
        title: "Frontend Developer",
        company_name: "Interface Labs",
        company_website: "https://interface.example.com",
        status: "interviewing",
        source: "Referral",
        notes: "Recruiter screen completed.",
        created_at: "2026-05-11T08:00:00Z"
      }
    ]
  };
}

/**
 * fillCreateApplicationForm
 *
 * Fills the required create application form fields for page tests.
 */
function fillCreateApplicationForm() {
  fireEvent.change(screen.getByLabelText("Application title"), {
    target: { value: "Backend Developer" }
  });

  fireEvent.change(screen.getByLabelText("Company name"), {
    target: { value: "Example Studio" }
  });
}

beforeEach(() => {
  mockedCreateApplication.mockReset();
  mockedGetApplications.mockReset();
});

describe("ApplicationsPage", () => {
  test("shows a loading state while applications are loading", () => {
    mockedGetApplications.mockReturnValue(new Promise(() => {}));

    renderApplicationsPage();

    expect(screen.getByText("Loading applications...")).toBeInTheDocument();
  });

  test("shows an empty state when no applications exist", async () => {
    mockedGetApplications.mockResolvedValue({ applications: [] });

    renderApplicationsPage();

    expect(await screen.findByRole("heading", { level: 2, name: "No applications yet" })).toBeInTheDocument();
  });

  test("renders loaded applications", async () => {
    mockedGetApplications.mockResolvedValue(buildApplicationsResponse());

    renderApplicationsPage();

    expect(await screen.findByRole("link", { name: "Backend Developer" })).toBeInTheDocument();
    expect(screen.getByText("Example Studio")).toBeInTheDocument();
    expect(screen.getAllByText("Applied").length).toBeGreaterThan(0);
    expect(screen.getByRole("link", { name: "Frontend Developer" })).toBeInTheDocument();
    expect(screen.getByText("Interface Labs")).toBeInTheDocument();
    expect(screen.getAllByText("Interviewing").length).toBeGreaterThan(0);
  });

  test("shows an error state when applications fail to load", async () => {
    mockedGetApplications.mockRejectedValue(new Error("network failed"));

    renderApplicationsPage();

    expect(
      await screen.findByRole("heading", { level: 2, name: "Applications need attention" })
    ).toBeInTheDocument();

    await waitFor(() => {
      expect(mockedGetApplications).toHaveBeenCalledTimes(1);
    });
  });

  test("creates an application and refreshes the list", async () => {
    mockedGetApplications
      .mockResolvedValueOnce({ applications: [] })
      .mockResolvedValueOnce(buildApplicationsResponse());
    mockedCreateApplication.mockResolvedValue(buildApplicationsResponse().applications[0]);

    renderApplicationsPage();

    expect(await screen.findByRole("heading", { level: 2, name: "No applications yet" })).toBeInTheDocument();

    fillCreateApplicationForm();
    fireEvent.click(screen.getByRole("button", { name: "Add application" }));

    await waitFor(() => {
      expect(mockedCreateApplication).toHaveBeenCalledWith(
        expect.objectContaining({
          title: "Backend Developer",
          companyName: "Example Studio"
        })
      );
    });

    expect(await screen.findByRole("status")).toHaveTextContent("Application added.");
    expect(screen.getByRole("link", { name: "Backend Developer" })).toBeInTheDocument();
  });

  test("shows an error when create application fails", async () => {
    mockedGetApplications.mockResolvedValue({ applications: [] });
    mockedCreateApplication.mockRejectedValue(new Error("save failed"));

    renderApplicationsPage();

    expect(await screen.findByRole("heading", { level: 2, name: "No applications yet" })).toBeInTheDocument();

    fillCreateApplicationForm();
    fireEvent.click(screen.getByRole("button", { name: "Add application" }));

    expect(await screen.findByRole("heading", { level: 2, name: "Applications need attention" })).toBeInTheDocument();
  });
});