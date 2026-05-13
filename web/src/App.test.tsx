import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, test, vi } from "vitest";

import App from "./App";

vi.mock("./api/applications", () => ({
  getApplications: vi.fn(() => new Promise(() => {}))
}));

/**
 * renderRoute
 *
 * Renders the application at a specific route for route-level component tests.
 * MemoryRouter avoids relying on the real browser URL during tests.
 */
function renderRoute(route: string) {
  return render(
    <MemoryRouter initialEntries={[route]}>
      <App />
    </MemoryRouter>
  );
}

describe("App", () => {
  test("renders the dashboard route", () => {
    renderRoute("/");

    expect(screen.getByRole("heading", { level: 1, name: "Dashboard" })).toBeInTheDocument();
  });

  test("renders the applications route", () => {
    renderRoute("/applications");

    expect(screen.getByRole("heading", { level: 1, name: "Applications" })).toBeInTheDocument();
  });

  test("renders the application detail route", () => {
    renderRoute("/applications/app-001");

    expect(screen.getByRole("heading", { level: 1, name: "Application Detail" })).toBeInTheDocument();
  });

  test("renders the not found route", () => {
    renderRoute("/missing");

    expect(screen.getByRole("heading", { level: 1, name: "Page not found" })).toBeInTheDocument();
  });
});