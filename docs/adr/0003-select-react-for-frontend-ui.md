# ADR-003: Select React for the Frontend UI

## Status

Accepted

## Context

ApplyBy will need an interactive user interface for managing applications, deadlines, contacts, interviews, follow-ups, documents, and analytics.

The frontend should support workflows such as application entry, pipeline/status management, application detail views, reminder review, search/filtering, and dashboard summaries.

React is being considered because it is widely used for interactive web interfaces and supports component-based UI design.

## Decision

Use React for the ApplyBy frontend UI.

React will be used to build the user-facing application interface.

## Consequences

React supports reusable UI components for forms, lists, cards, status displays, reminders, and dashboard panels.

The project gains a practical frontend technology commonly used in industry.

Care will be needed to keep core business rules out of UI components.

The project will need frontend build tooling and frontend test tooling once implementation begins.