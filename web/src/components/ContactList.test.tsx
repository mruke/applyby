import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, test, vi } from "vitest";

import type { ContactResponse } from "../types/application";
import { ContactList } from "./ContactList";

// -----------------------------------------------------------------------------
// buildContact
//
// Creates a contact response for contact list tests.
// -----------------------------------------------------------------------------
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
    render(
      <MemoryRouter>
        <ContactList applicationId="app-001" contacts={[]} />
      </MemoryRouter>
    );

    expect(screen.getByText("No contacts added yet.")).toBeInTheDocument();
  });

  test("renders contacts", () => {
    render(
      <MemoryRouter>
        <ContactList applicationId="app-001" contacts={[buildContact()]} />
      </MemoryRouter>
    );

    expect(screen.getByText("Sam Recruiter")).toBeInTheDocument();
    expect(screen.getByText("Recruiter · sam@example.com")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Edit" })).toHaveAttribute(
      "href",
      "/applications/app-001/contacts/contact-001/edit"
    );
  });

  test("removes a contact", async () => {
    const onRemove = vi.fn().mockResolvedValue(undefined);

    render(
      <MemoryRouter>
        <ContactList applicationId="app-001" contacts={[buildContact()]} onRemove={onRemove} />
      </MemoryRouter>
    );

    fireEvent.click(screen.getByRole("button", { name: "Remove" }));

    await waitFor(() => {
      expect(onRemove).toHaveBeenCalledWith("contact-001");
    });
  });
});
