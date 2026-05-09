# ADR-005: Select a Layered Testing Strategy

## Status

Accepted

## Context

ApplyBy is intended to demonstrate disciplined software engineering, not just feature completion.

The project will contain domain rules, backend use cases, persistence behavior, search/filter behavior, reminder prioritization, frontend interactions, and end-to-end user workflows.

The author's previous portfolio projects use a clear testing structure with separate unit, integration, end-to-end, and helper test areas. ApplyBy should preserve that general structure while adapting the specific tools to the selected stack.

## Decision

Use a layered testing strategy for ApplyBy.

The repository will keep separate areas for unit tests, integration tests, end-to-end tests, and shared test helpers.

Unit tests will focus on domain rules, lifecycle transitions, reminder prioritization, search helpers, analytics calculations, and small services.

Integration tests will focus on persistence, repository behavior, API boundaries, and database-backed workflows.

End-to-end tests will focus on a small number of complete user workflows once a runnable UI exists.

## Consequences

The testing structure stays consistent with the author's existing portfolio style.

Domain rules and data structures can be tested without requiring the full application to run.

Persistence and API behavior can be tested separately from the UI.

End-to-end tests can validate critical workflows without becoming the main source of coverage.

The test suite will require more organization than a minimal project, and integration tests will require a database setup strategy.