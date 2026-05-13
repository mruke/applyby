import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, test, vi } from "vitest";

import type { ApplicationSearchCriteria } from "../types/application";
import { ApplicationSearchForm } from "./ApplicationSearchForm";

/**
 * emptyCriteria
 *
 * Provides default empty search criteria for component tests.
 */
const emptyCriteria: ApplicationSearchCriteria = {
  companyName: "",
  source: "",
  statuses: [],
  text: ""
};

describe("ApplicationSearchForm", () => {
  test("submits search criteria", async () => {
    const onSearch = vi.fn().mockResolvedValue(undefined);

    render(
      <ApplicationSearchForm
        criteria={emptyCriteria}
        isSearching={false}
        onClear={vi.fn()}
        onSearch={onSearch}
      />
    );

    fireEvent.change(screen.getByLabelText("Text"), {
      target: { value: "backend" }
    });

    fireEvent.change(screen.getByLabelText("Company"), {
      target: { value: "Example Studio" }
    });

    fireEvent.change(screen.getByLabelText("Source"), {
      target: { value: "Company site" }
    });

    fireEvent.change(screen.getByLabelText("Status"), {
      target: { value: "applied" }
    });

    fireEvent.click(screen.getByRole("button", { name: "Search" }));

    await waitFor(() => {
      expect(onSearch).toHaveBeenCalledWith({
        companyName: "Example Studio",
        source: "Company site",
        statuses: ["applied"],
        text: "backend"
      });
    });
  });

  test("clears search criteria", async () => {
    const onClear = vi.fn().mockResolvedValue(undefined);

    render(
      <ApplicationSearchForm
        criteria={{
          companyName: "Example Studio",
          source: "Company site",
          statuses: ["applied"],
          text: "backend"
        }}
        isSearching={false}
        onClear={onClear}
        onSearch={vi.fn()}
      />
    );

    fireEvent.click(screen.getByRole("button", { name: "Clear search" }));

    await waitFor(() => {
      expect(onClear).toHaveBeenCalledTimes(1);
    });
  });

  test("disables actions while searching", () => {
    render(
      <ApplicationSearchForm
        criteria={emptyCriteria}
        isSearching={true}
        onClear={vi.fn()}
        onSearch={vi.fn()}
      />
    );

    expect(screen.getByRole("button", { name: "Searching..." })).toBeDisabled();
    expect(screen.getByRole("button", { name: "Clear search" })).toBeDisabled();
  });
});