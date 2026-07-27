# ApplyBy

ApplyBy is a personal job application CRM for tracking applications, deadlines, companies, contacts, interview status, follow-ups, documents, and activity history.

It's built as both a practical job-search tool and a portfolio project focused on disciplined design: clear domain modeling, validated workflows, searchable records, and reminder prioritization.

## At a Glance

- **Type:** Local-first, single-user job application CRM
- **Stack:** Go backend, React + TypeScript frontend, PostgreSQL
- **Status:** First prototype, complete and usable

Implementation decisions (why Go, why PostgreSQL, why this test strategy, etc.) are recorded individually in [`docs/adr/`](docs/adr/).

## Demo

ApplyBy is demonstrated with fictional demo data.

- [Watch the demo video](https://youtu.be/ObvVGlxBuSs)
- See [QUICKSTART.md](./QUICKSTART.md) for local setup

### Dashboard filtering

Filter applications directly from the dashboard using prebuilt summary buttons for quick triage.

<img src="./docs/assets/dashboard-filtering.gif" alt="Dashboard filtering" width="500">

### Create application

The applications workbench is the main management view for adding new job applications.

<img src="./docs/assets/create-application.gif" alt="Create application" width="350">

### Workbench filtering

The workbench supports more granular filtering by application details and workflow state.

<img src="./docs/assets/workshop-filter.gif" alt="Workbench filtering" width="550">

### Application activity tracking

Each application has its own workflow area for status updates, reminders, contacts, and document metadata (URI-style references, not file uploads). Changes are tracked in a per-application activity feed.

<img src="./docs/assets/activity-tracking.gif" alt="Activity tracking" width="700">

## Project Goals

ApplyBy helps answer questions like:

- Which jobs have I applied to, and which need follow-up?
- Which applications have upcoming deadlines?
- Which companies and contacts am I in touch with?
- Which interview stages are active?
- Which resume or document version did I use?
- Which sources and strategies are actually producing responses?

## Implementation Status

| Step | Focus Area | Status |
|---|---|---|
| 1 | Backend foundation (Go module, package structure) | Complete |
| 2 | Domain model (applications, companies, contacts, reminders, documents, activity) | Complete |
| 3 | Application lifecycle (status validation, status history) | Complete |
| 4 | Application workflows (create, update, list, search) | Complete |
| 5 | PostgreSQL persistence (schema, repositories, integration tests) | Complete |
| 6 | Search and reminders | Complete |
| 7 | HTTP API | Complete |
| 8 | Frontend foundation (React/TypeScript structure, routes, API client) | Complete |
| 9 | User interface (dashboard, workbench, detail views, forms) | Complete |
| 10 | CRUD completion (editing, removal across all record types) | Complete |
| 11 | UX hardening (accessibility, layout, dashboard filtering) | Complete |
| 12 | Documentation and engineering audit | Complete |
| 13 | Analytics and benchmarks | Deferred |
| 14 | Deployment and packaging | Deferred |

The prototype is finished and usable. Remaining items are enhancements, not blockers.

## Current CRUD Coverage

| Area | Create | Read | Update | Delete | Notes |
|---|---:|---:|---:|---:|---|
| Applications | Yes | Yes | Yes | Yes | Removing an application also removes its related records. |
| Reminders | Yes | Yes | Yes | Yes | Can be scheduled, edited, completed, and removed. |
| Contacts | Yes | Yes | Yes | Yes | Can be added, edited, and removed. |
| Document metadata | Yes | Yes | Yes | Yes | Metadata only, file upload/storage is deferred. |
| Activity history | System-generated | Yes | No | Cascades with application | Append-only during normal use. |

## Not Yet Implemented

Deferred beyond the current single-user prototype:

- Authentication, authorization, and multi-user data ownership
- Hosted deployment and packaging for non-technical users
- File upload and file storage (document metadata only, for now)
- Calendar integrations and notification delivery
- Job-search analytics (response rates by source, time-to-response, interview conversion)
- Generated data and benchmark coverage
- End-to-end browser automation

## Documentation

- [`QUICKSTART.md`](./QUICKSTART.md) — setup and run commands
- [`ARCHITECTURE.md`](./ARCHITECTURE.md) — system boundaries, design strategy, and tradeoffs
- [`docs/FRONTEND_UX.md`](./docs/FRONTEND_UX.md) — frontend UX and accessibility principles
- [`docs/adr/`](./docs/adr/) — individual architecture decisions

## AI Assistance Disclosure

ChatGPT was used during development as a learning, design, and review assistant, for project planning, architecture discussion, documentation drafting, implementation review, and test-suite design. Some tests were generated with ChatGPT assistance, then reviewed, adapted, and validated as part of this repository.

## License

MIT.
