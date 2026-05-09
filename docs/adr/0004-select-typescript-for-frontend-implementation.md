# ADR-004: Select TypeScript for Frontend Implementation

## Status

Accepted

## Context

ApplyBy will handle structured frontend data such as applications, companies, contacts, reminders, interviews, documents, statuses, filters, and analytics summaries.

Plain JavaScript would be possible, but explicit types can make frontend contracts easier to understand and maintain.

TypeScript is being considered because it adds static typing to frontend code and can help catch mismatches between API responses, UI components, form data, and frontend state.

## Decision

Use TypeScript for the ApplyBy frontend implementation.

Frontend data structures, component props, API response shapes, form values, and status-related values should be typed where doing so improves clarity.

## Consequences

TypeScript makes frontend data contracts more explicit.

Typed component props and API response shapes can improve maintainability.

The project gains a useful industry-relevant frontend skill.

TypeScript adds configuration and build-tooling complexity.

Care will be needed to avoid unnecessary type-level complexity.