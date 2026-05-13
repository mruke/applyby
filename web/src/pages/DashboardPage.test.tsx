import { render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, test, vi } from "vitest";

import { getApplications } from "../api/applications";
import type { ApplicationsResponse } from "../types/application";
import { DashboardPage } from "./DashboardPage";

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
 * buildApplicationsResponse
 *
 * Creates an applications API response for dashboard tests.
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
      },
      {
        id: "app-003",
        title: "Platform Developer",
        company_name: "Systems Co",
        company_website: "https://systems.example.com",
        status: "offer",
        source: "Recruiter",
        notes: "Offer received.",
        created_at: "2026-05-12T08:00:00Z"
      }
    ]
  };
}

beforeEach(() => {
  mockedGetApplications.mockReset();
});

describe("DashboardPage", () => {
  test("shows a loading state while dashboard data loads", () => {
    mockedGetApplications.mockReturnValue(new Promise(() => {}));

    render(<DashboardPage />);

    expect(screen.getByText("Loading dashboard...")).toBeInTheDocument();
  });

  test("renders dashboard metrics", async () => {
    mockedGetApplications.mockResolvedValue(buildApplicationsResponse());

    render(<DashboardPage />);

    expect(await screen.findByRole("heading", { level: 2, name: "Total applications" })).toBeInTheDocument();
    expect(screen.getByRole("heading", { level: 2, name: "Active applications" })).toBeInTheDocument();
    expect(screen.getByRole("heading", { level: 2, name: "Interviewing" })).toBeInTheDocument();
    expect(screen.getByRole("heading", { level: 2, name: "Offers" })).toBeInTheDocument();
  });

  test("shows an error state when dashboard data fails to load", async () => {
    mockedGetApplications.mockRejectedValue(new Error("network failed"));

    render(<DashboardPage />);

    expect(
      await screen.findByRole("heading", { level: 2, name: "Dashboard could not be loaded" })
    ).toBeInTheDocument();
  });
});