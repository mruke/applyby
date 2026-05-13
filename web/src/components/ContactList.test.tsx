import { render, screen } from "@testing-library/react";
import { describe, expect, test } from "vitest";

import type { ContactResponse } from "../types/application";
import { ContactList } from "./ContactList";

/**
 * buildContact
 *
 * Creates a contact response for contact list tests.
 */
function buildContact(): ContactResponse {
  return {
    id: "contact-001",
    application_id: "app-001",
    name: "Sam Recruiter",
    email: "sam@example.com",
    role: "Recruiter"
  };
}

describe("ContactList", () => {
  test("renders an empty contact state", () => {
    render(<ContactList contacts={[]} />);

    expect(screen.getByText("No contacts added yet.")).toBeInTheDocument();
  });

  test("renders contacts", () => {
    render(<ContactList contacts={[buildContact()]} />);

    expect(screen.getByText("Sam Recruiter")).toBeInTheDocument();
    expect(screen.getByText("Recruiter · sam@example.com")).toBeInTheDocument();
    
  });
});