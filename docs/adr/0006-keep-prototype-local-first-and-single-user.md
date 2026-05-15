# 0006. Keep Prototype Local-First and Single-User

## Status

Accepted

## Context

ApplyBy began as a personal job application tracker and portfolio project. The immediate goal was to build a complete, usable local prototype for one user rather than a hosted multi-user platform.

The application now supports the core job-search workflow: applications, status updates, contacts, reminders, document metadata, activity history, search, dashboard summary filtering, and application removal.

Adding authentication, authorization, hosted deployment, account ownership, and multi-user tenancy would substantially increase system complexity. Those concerns are important for a production service, but they are not required for the current prototype goal.

## Decision

ApplyBy will remain local-first and single-user for the completed prototype.

The current system will not include authentication, authorization, account ownership, or multi-tenant data separation. Those concerns are deferred until the project direction explicitly shifts toward hosted or shared use.

## Consequences

Positive consequences:

- The prototype remains focused on the core job-search tracking workflow.
- The architecture avoids premature account and tenancy complexity.
- Local development and manual testing remain simple.
- The app remains easier to reason about as a portfolio project.

Negative consequences:

- The app is not currently safe to host as a shared public application.
- Multiple users cannot safely use the same deployment with separated data.
- Future hosted deployment will require a deliberate authentication and ownership design.

Follow-up work if this decision changes:

- Add user identity and ownership to the data model.
- Add authentication and authorization.
- Add tests for access control and data isolation.
- Review every repository query and API route for ownership filtering.