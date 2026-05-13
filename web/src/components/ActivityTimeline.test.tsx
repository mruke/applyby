import { render, screen } from "@testing-library/react";
import { describe, expect, test } from "vitest";

import type { ActivityEventResponse } from "../types/application";
import { ActivityTimeline } from "./ActivityTimeline";

/**
 * buildActivityEvent
 *
 * Creates an activity event response for timeline component tests.
 */
function buildActivityEvent(): ActivityEventResponse {
  return {
    application_id: "app-001",
    type: "status_changed",
    occurred_at: "2026-05-14T09:00:00Z",
    description: "Status changed from applied to interviewing."
  };
}

describe("ActivityTimeline", () => {
  test("renders an empty activity state", () => {
    render(<ActivityTimeline events={[]} />);

    expect(screen.getByText("No activity recorded yet.")).toBeInTheDocument();
  });

  test("renders activity events", () => {
    render(<ActivityTimeline events={[buildActivityEvent()]} />);

    expect(screen.getByText("Status changed from applied to interviewing.")).toBeInTheDocument();
    expect(screen.getByText(/status_changed/)).toBeInTheDocument();
  });
});