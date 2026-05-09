# ADR-002: Select PostgreSQL for Persistence

## Status

Accepted

## Context

ApplyBy will manage structured job-search data, including applications, companies, contacts, interviews, reminders, documents, status history, and activity events.

The domain is naturally relational. Companies can have many applications, applications can have many events, contacts can be linked to companies or applications, and reminders can be linked to applications or interviews.

The project also has a performance-oriented goal. Search, filtering, reminders, and analytics should be supported by intentional data access patterns rather than only in-memory filtering.

PostgreSQL is being considered because it supports relational modeling, constraints, indexes, full-text search, and analytics-style queries.

## Decision

Use PostgreSQL as the initial persistence layer for ApplyBy.

The database will store the core application data and support indexed query paths for search, filtering, reminders, and analytics.

## Consequences

PostgreSQL fits the relational shape of the ApplyBy domain.

Database constraints and indexes can support data integrity and performance.

The project gains a more industry-standard persistence layer than a purely local file or SQLite-only approach.

PostgreSQL requires more setup than SQLite and will require a clear local development and test database strategy.