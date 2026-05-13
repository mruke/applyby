import { render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { beforeEach, describe, expect, test, vi } from "vitest";

import { getApplications } from "../api/applications";
import type { ApplicationsResponse } from "../types/application";
import { ApplicationsPage } from "./ApplicationsPage";

vi.mock("../api/applications", () => ({
  getApplications: vi.fn()
}));

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

beforeEach(() => {
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
    expect(screen.getByText("Applied")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Frontend Developer" })).toBeInTheDocument();
    expect(screen.getByText("Interface Labs")).toBeInTheDocument();
    expect(screen.getByText("Interviewing")).toBeInTheDocument();
  });

  test("shows an error state when applications fail to load", async () => {
    mockedGetApplications.mockRejectedValue(new Error("network failed"));

    renderApplicationsPage();

    expect(
      await screen.findByRole("heading", { level: 2, name: "Applications could not be loaded" })
    ).toBeInTheDocument();

    await waitFor(() => {
      expect(mockedGetApplications).toHaveBeenCalledTimes(1);
    });
  });
});