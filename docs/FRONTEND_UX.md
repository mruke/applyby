# Frontend UX Principles

ApplyBy is a job application tracker. The interface should be easy to understand.

This document defines frontend UX and accessibility principles for the UI implementation. It is intended to be practical. It should guide layout, components, forms, state feedback, and page structure without becoming a full visual design system.

## Goals

The UI should help the user quickly answer:

- What applications exist?
- What status is each application in?
- What needs attention next?
- What reminders are due or overdue?
- What changed recently?
- What contacts and documents are attached to an application?
- What action can I take from this screen?

## Core Principles

| Principle | ApplyBy Rule |
|---|---|
| Observability | Important state should be visible without forcing the user to dig through the interface. |
| Readability | Pages should be easy to scan, with clear headings, labels, spacing, and grouping. |
| Usability | Common actions should be obvious and reachable from the page where the user naturally needs them. |
| Accessibility | Semantic HTML, labels, keyboard access, visible focus, and readable status text are required from the start. |
| Consistency | Status labels, form patterns, actions, tables, cards, and empty states should behave the same way across the app. |
| Recognition over recall | The UI should show available options and actions instead of requiring the user to remember valid states or workflows. |
| Feedback | User actions should produce clear visible feedback, such as saved, loading, error, empty, or updated states. |
| Error recovery | Validation errors should explain what needs fixing and preserve entered form data where possible. |

## Accessibility Baseline

ApplyBy should follow these accessibility defaults during frontend implementation:

- Use semantic elements such as `header`, `nav`, `main`, `section`, `form`, `label`, `button`, and `table` where appropriate.
- Use one clear `h1` per page and a logical heading order after it.
- Give every form control a visible or programmatic label.
- Make all interactive controls keyboard accessible.
- Preserve visible focus states for links, buttons, inputs, and custom controls.
- Do not rely on color alone to communicate status.
- Use readable text labels for statuses, errors, filters, and actions.
- Provide clear loading, error, empty, and success states.
- Keep validation messages close to the fields they describe.
- Prefer native HTML controls unless a custom component has a clear reason to exist.

## Observability Rules

The user should be able to understand the system state quickly.

Important states to make visible:

| State | UI Treatment |
|---|---|
| Application status | Text badge such as `Applied`, `Interviewing`, `Offer`, `Rejected`, `Withdrawn`, or `Archived`. |
| Next reminder | Show the nearest incomplete reminder when available. |
| Overdue reminder | Use text plus visual emphasis; do not rely only on color. |
| Active filters | Show applied filters near the application list. |
| Recent activity | Show recent status changes and events on the application detail page. |
| Empty sections | Explain what is missing and how to add it. |
| Loading state | Show that data is being fetched. |
| Error state | Explain what failed and what the user can do next. |
| Save result | Confirm when an action succeeds or fails. |

## Readability Rules

ApplyBy is scan-heavy. The UI should support quick reading.

Use clear section labels:

- Applications
- Status
- Reminders
- Contacts
- Documents
- Activity
- Search and Filters

Avoid vague labels:

- Info
- Data
- Misc
- Items
- Stuff

Prefer short, concrete field labels:

| Good | Avoid |
|---|---|
| Company | Company info |
| Status | Current state |
| Source | Origin |
| Next reminder | Reminder data |
| Last activity | Activity info |

## Layout Direction

The first version should use a simple layout.

Recommended page structure:

| Page | Purpose |
|---|---|
| Dashboard | Summarize items needing attention. |
| Applications | Provide the main application workbench. |
| Application detail | Show one application's related records and activity. |
| Forms | Create or update structured records with clear validation. |

Ideal app shell:

```text
Header or top navigation
Main content area
Page heading
Primary page action
Primary content
Secondary sections
```

The layout should be functional before it is decorative!

## Application List Expectations

The applications list must support scanning.

Recommended visible fields:

- Title
- Company
- Status
- Source
- Next reminder
- Last activity or updated date
- Primary action to view details

The list can be a table or structured list. The choice should depend on readability and responsive behavior.

## Application Detail Expectations

The application detail page should show the full context of one application.

Recommended sections:

1. Summary
2. Status
3. Reminders
4. Contacts
5. Documents
6. Activity

The most important actions should be easy to find:

- Edit application details
- Update status
- Schedule reminder
- Edit or remove reminder
- Add contact
- Edit or remove contact
- Add document metadata
- Edit or remove document metadata
- Remove application

## Forms

Forms should be predictable and forgiving.

Form rules:

- Use clear labels.
- Mark required fields.
- Use constrained controls for known values, such as status.
- Keep user-entered values when validation fails.
- Place validation messages near the relevant field.
- Use clear submit and cancel/back actions.
- Avoid long forms when a smaller focused form would be clearer.

## Status Labels

Statuses should always be text-first.

Supported application statuses:

- Draft
- Interested
- Applied
- Interviewing
- Offer
- Rejected
- Withdrawn
- Archived

Color may support status recognition, but text must carry the meaning.

## Empty States

Empty states should explain both the current state and the next action.

Examples:

| Area | Empty State |
|---|---|
| Applications | No applications yet. Add your first application to start tracking. |
| Reminders | No reminders scheduled for this application. |
| Contacts | No contacts added yet. |
| Documents | No document metadata added yet. |
| Activity | No activity recorded yet. |

## Feedback and Error States

The user should not have to guess what happened.

Use clear feedback for:

- application created
- application edited
- application removed
- status updated
- reminder scheduled
- reminder edited
- reminder completed
- reminder removed
- contact added
- contact edited
- contact removed
- document metadata added
- document metadata edited
- document metadata removed
- dashboard filter changed
- search returned no results
- save failed
- network or API failure

Errors should be direct and recoverable. Avoid vague messages such as `Something went wrong` unless paired with a useful next action.

## Frontend Architecture Rules

The frontend should preserve the backend boundaries.

| Rule | Reason |
|---|---|
| Keep API calls in an API client layer. | Components should not know request details. |
| Keep shared TypeScript types in a dedicated area. | API data shapes should be easy to find. |
| Keep business rules out of React components. | Backend domain/application layers own workflow correctness. |
| Keep route-level pages separate from reusable components. | Page orchestration and component display have different responsibilities. |
| Keep UI state local unless it needs to be shared. | Avoid unnecessary global state. |

## Frontend Implementation Notes

The frontend includes the route-level pages, API client boundary, shared types, feedback states, application list/detail views, dashboard overview, and record-management sections needed for the current prototype.

