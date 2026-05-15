import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import type { ComponentProps } from "react";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, test, vi } from "vitest";

import type { ContactResponse } from "../types/application";
import { ContactSection } from "./ContactSection";

// -----------------------------------------------------------------------------
// buildContact
//
// Creates a contact response for contact section tests.
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

// -----------------------------------------------------------------------------
// renderContactSection
//
// Renders the contact section with default callbacks and overridable props.
// -----------------------------------------------------------------------------
function renderContactSection(overrides: Partial<ComponentProps<typeof ContactSection>> = {}) {
  const props: ComponentProps<typeof ContactSection> = {
    applicationId: "app-001",
    contacts: [buildContact()],
    errorMessage: null,
    isAdding: false,
    isRemoving: false,
    onAdd: vi.fn().mockResolvedValue(undefined),
    onRemove: vi.fn().mockResolvedValue(undefined),
    ...overrides
  };

  render(
    <MemoryRouter>
      <ContactSection {...props} />
    </MemoryRouter>
  );

  return props;
}

describe("ContactSection", () => {
  test("renders the add form and contact list", () => {
    renderContactSection();

    expect(screen.getByRole("heading", { level: 2, name: "Add contact" })).toBeInTheDocument();
    expect(screen.getByText("Sam Recruiter")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Edit" })).toHaveAttribute(
      "href",
      "/applications/app-001/contacts/contact-001/edit"
    );
  });

  test("shows a contact load error while keeping the add form available", () => {
    renderContactSection({ errorMessage: "Contacts could not be loaded." });

    expect(screen.getByRole("heading", { level: 2, name: "Add contact" })).toBeInTheDocument();
    expect(screen.getByRole("alert")).toHaveTextContent("Contacts could not be loaded.");
    expect(screen.queryByText("Sam Recruiter")).not.toBeInTheDocument();
  });

  test("delegates remove actions to the parent page", async () => {
    const onRemove = vi.fn().mockResolvedValue(undefined);

    renderContactSection({ onRemove });

    fireEvent.click(screen.getByRole("button", { name: "Remove" }));

    await waitFor(() => {
      expect(onRemove).toHaveBeenCalledWith("contact-001");
    });
  });
});
