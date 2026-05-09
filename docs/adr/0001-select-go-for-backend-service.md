# ADR-001: Select Go for the Backend Service

## Status

Accepted

## Context

ApplyBy is intended to be a personal job application CRM with an emphasis on clear backend boundaries, performance-conscious design, and practical use of data structures and algorithms.

The existing portfolio projects are Python-based. ApplyBy should add a different signal by using an industry-relevant backend language outside the author's current primary project experience.

The backend is expected to own core application behavior such as application lifecycle rules, reminder handling, activity history, persistence coordination, search-related operations, and analytics preparation.

Go is being considered because it is commonly used for backend services and performance-conscious systems. It also provides an opportunity to practice static typing, compiled deployment, explicit error handling, concurrency, and built-in benchmarking.

## Decision

Use Go for the ApplyBy backend service.

The backend will own core application behavior and expose the API used by the frontend.

## Consequences

Go adds a new backend language to the portfolio and supports the project's performance-oriented goals.

The project can use Go's standard testing and benchmarking tools to measure application behavior.

Development will be slower at first because Go is new to the author.

The project will require care around Go package structure, error handling, dependency management, and frontend-backend integration.