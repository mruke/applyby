# 0008. Separate Dashboard Summary from Applications Workbench

## Status

Accepted

## Context

ApplyBy has two high-level frontend needs:

1. A fast landing page that gives a user quick access to relevant applications.
2. A full workbench for creating, searching, filtering, and scanning application records.

During development, the dashboard evolved from static metrics into an interactive summary view. The applications page evolved into a split layout with application creation on one side and the searchable application table on the other.

Duplicating all application-management controls on the dashboard would make the dashboard compete with the applications page and increase UI maintenance cost.

## Decision

The dashboard will remain a summary and navigation screen.

Dashboard summary cards filter an application list in place for fast access. Application titles link to detail pages. The dashboard may link to the applications page, but it should not duplicate the full application creation/search/edit workbench.

The applications page remains the primary workbench for application creation, search, filtering, and list scanning.

## Consequences

Positive consequences:

- The dashboard stays focused and lightweight.
- The applications page has a clear responsibility as the full workbench.
- Users can quickly open relevant applications without losing the dedicated management screen.
- Future dashboard features have a clear boundary.

Negative consequences:

- Some application information appears in both dashboard and applications views.
- The dashboard needs frontend filtering logic separate from backend-backed application search.
- Users must navigate to the applications page for creation and full search workflows.

Guidance:

- Add summary and fast-access features to the dashboard.
- Add creation, full search, bulk workflows, and table-level management to the applications page.
- Avoid turning the dashboard into a second application workbench.