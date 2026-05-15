# 0007. Use Explicit Workflow Services

## Status

Accepted

## Context

ApplyBy has workflows that are more specific than basic table mutations. Examples include:

- creating an application
- updating application details
- updating application status
- scheduling, editing, completing, and removing reminders
- adding, editing, and removing contacts
- adding, editing, and removing document metadata
- removing an application and relying on persistence rules for related data
- recording activity events and status history

A generic CRUD service or generic handler layer would reduce file count, but it would also hide workflow intent and make side effects less explicit.

ApplyBy favors clear boundaries and readable workflows over a minimal number of abstractions.

## Decision

ApplyBy will use explicit application-layer workflow services for meaningful user actions.

Each workflow should have a focused service, input model, repository interface dependency, validation path, and tests. API handlers should call these services instead of owning workflow rules directly.

## Consequences

Positive consequences:

- Workflow intent is clear from file and type names.
- Tests can target behavior at the application-service boundary.
- API handlers remain thin.
- Activity/history side effects are easier to see and test.
- Future changes can be made to one workflow without changing a generic mutation engine.

Negative consequences:

- The codebase contains more files.
- Contact, document, and reminder workflows can look repetitive.
- Some workflow services may appear small.

Guidance:

- Do not introduce generic CRUD abstractions unless repetition is causing real defects or confusing changes.
- Prefer small explicit services over clever indirection.
- Refactor repeated workflow patterns only when the abstraction preserves business intent.