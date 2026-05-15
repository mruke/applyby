import { fireEvent, render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
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
 * renderDashboardPage
 *
 * Renders the dashboard inside a router for link tests.
 */
function renderDashboardPage() {
  return render(
    <MemoryRouter>
      <DashboardPage />
    </MemoryRouter>
  );
}

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
      },
      {
        id: "app-004",
        title: "Data Developer",
        company_name: "Archive Ltd",
        company_website: "https://archive.example.com",
        status: "archived",
        source: "Job board",
        notes: "Archived older opportunity.",
        created_at: "2026-05-09T08:00:00Z"
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

    renderDashboardPage();

    expect(screen.getByText("Loading dashboard...")).toBeInTheDocument();
  });

  test("renders dashboard summary filter buttons and applications", async () => {
    mockedGetApplications.mockResolvedValue(buildApplicationsResponse());

    renderDashboardPage();

    expect(await screen.findByRole("button", { name: /Total applications/i })).toHaveAttribute("aria-pressed", "true");
    expect(screen.getByRole("button", { name: /Active applications/i })).toHaveAttribute("aria-pressed", "false");
    expect(screen.getByRole("heading", { level: 2, name: "Applications" })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Open applications workbench" })).toHaveAttribute("href", "/applications");
    expect(screen.getByRole("link", { name: "Platform Developer" })).toHaveAttribute("href", "/applications/app-003");
    expect(screen.getByRole("link", { name: "Data Developer" })).toHaveAttribute("href", "/applications/app-004");
  });

  test("filters dashboard applications when a summary card is clicked", async () => {
    mockedGetApplications.mockResolvedValue(buildApplicationsResponse());

    renderDashboardPage();

    expect(await screen.findByRole("link", { name: "Backend Developer" })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Data Developer" })).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: /Interviewing/i }));

    expect(screen.getByRole("button", { name: /Interviewing/i })).toHaveAttribute("aria-pressed", "true");
    expect(screen.getByRole("link", { name: "Frontend Developer" })).toBeInTheDocument();
    expect(screen.queryByRole("link", { name: "Backend Developer" })).not.toBeInTheDocument();
    expect(screen.queryByRole("link", { name: "Platform Developer" })).not.toBeInTheDocument();
    expect(screen.queryByRole("link", { name: "Data Developer" })).not.toBeInTheDocument();
  });

  test("shows an empty summary view when no applications match the selected card", async () => {
    mockedGetApplications.mockResolvedValue({
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
        }
      ]
    });

    renderDashboardPage();

    expect(await screen.findByRole("link", { name: "Backend Developer" })).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: /Offers/i }));

    expect(screen.getByText("No applications match this summary view.")).toBeInTheDocument();
  });

  test("shows an error state when dashboard data fails to load", async () => {
    mockedGetApplications.mockRejectedValue(new Error("network failed"));

    renderDashboardPage();

    expect(
      await screen.findByRole("heading", { level: 2, name: "Dashboard could not be loaded" })
    ).toBeInTheDocument();
  });
});