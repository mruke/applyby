# ApplyBy

ApplyBy is a personal job application CRM for tracking applications, deadlines, companies, contacts, interviews, follow-ups, documents, and outcomes.

The project is intended to be a practical job-search tool and a portfolio project focused on disciplined application design. It emphasizes clear domain modeling, validated workflows, performance-conscious data structures, searchable records, reminder prioritization, activity history, and job-search analytics.

## Project Goals

ApplyBy is designed to help an individual answer questions such as:

- Which jobs have I applied to?
- Which applications have upcoming deadlines?
- Which applications need follow-up?
- Which companies and contacts am I interacting with?
- Which interview stages are active?
- Which resume or document version did I use?
- Which sources and application strategies are producing responses?

## Planned Product Features

- Application tracking
- Deadline tracking
- Company and contact management
- Interview tracking
- Follow-up reminders
- Resume and document version tracking
- Activity timeline
- Search, filtering, and sorting
- Job-search analytics dashboard

## Planned Engineering Features

- Finite-state machine for application status transitions
- Append-only activity event history
- Priority queue for reminders
- Indexed search and filtering
- Precomputed analytics summaries
- Relationship traversal across applications, companies, contacts, and documents
- Clear separation between domain, application, storage, search, reminders, analytics, and interface layers
- Unit, integration, and end-to-end tests
- Benchmark suite with generated data

## Scope

ApplyBy is not intended to be a production SaaS platform in its first version.

The first version focuses on a single-user job-search workflow with local or simple hosted persistence. Authentication, employer-facing workflows, file upload and file storage, calendar integrations, and production deployment wiring are future possibilities, not initial requirements.

## Selected Implementation Direction

ApplyBy will be implemented as a full-stack application with:

- Go for the backend service
- PostgreSQL for persistence
- React for the frontend UI
- TypeScript for frontend implementation
- A layered testing strategy with unit, integration, end-to-end, and helper test areas

These decisions are recorded in `docs/adr/`.

## Implementation Steps

The project is planned as a sequence of focused implementation steps.

| Step | Focus Area | Goal | Status |
|---|---|---|---|
| 1 | Backend foundation | Initialize the Go backend module, package structure, and first test workflow. | Complete |
| 2 | Domain model | Define the core job application, company, contact, reminder, document, and activity concepts. | Complete |
| 3 | Application lifecycle | Implement validated application status transitions using an explicit lifecycle model. | Complete |
| 4 | Application workflows | Add use cases for creating, updating, listing, and managing applications while keeping workflows separate from storage and API code. | Complete |
| 5 | PostgreSQL persistence | Add relational schema, migrations, repositories, constraints, and integration tests for durable storage. | Complete |
| 6 | Search and reminders | Add indexed search/filter behavior and reminder prioritization for performance-conscious job-search workflows. | Complete |
| 7 | HTTP API | Expose backend workflows through thin request/response handlers without placing business rules in the API layer. | Complete |
| 8 | Frontend foundation | Initialize the React and TypeScript frontend structure, route layout, API client boundary, and frontend test setup. | Complete |
| 9 | User interface | Add application list, creation, detail, status, reminder, activity, contact, document metadata, search, filtering, and dashboard views. | Complete |
| 10 | Analytics and benchmarks | Add generated data and benchmark coverage for search, reminders, analytics, and relationship traversal. | Planned |
| 11 | Documentation and polish | Finalize setup instructions, architecture notes, ADRs, validation commands, and portfolio framing. | In progress |
| 12 | Application detail editing | Add non-status application detail editing for title, company, website, source, and notes. Keep lifecycle status updates separate. | Planned |
| 13 | Contact maintenance workflows | Add contact editing and removal workflows. | Planned |
| 14 | Document metadata maintenance workflows | Add document metadata editing and removal workflows. File upload/storage remains deferred. | Planned |
| 15 | Reminder maintenance workflows | Add reminder editing and canceling workflows. Keep reminder completion as a separate workflow. | Planned |
| 16 | Application removal policy | Decide and implement the intended application removal behavior, including how it should affect reminders, contacts, documents, and activity history. | Planned |

## Current CRUD Coverage

ApplyBy currently supports the core local CRM workflow, but it is not yet full CRUD for every record type.

| Area | Create | Read | Update | Delete / Remove | Current gap |
|---|---:|---:|---:|---:|---|
| Applications | Yes | Yes | Partial | Partial | Details can be created/read; status can be updated; general detail editing is planned. Application removal behavior still needs an explicit product decision. |
| Reminders | Yes | Yes | Partial | No | Reminders can be scheduled/read/completed; editing and canceling are planned. |
| Contacts | Yes | Yes | No | No | Contacts can be added/read; editing and removal are planned. |
| Document metadata | Yes | Yes | No | No | Document metadata can be added/read; editing and removal are planned. |
| Activity history | System-generated | Yes | No | No | Activity is append-only by design. |

## Status

The current single-user version includes the Go backend API server, PostgreSQL persistence, React and TypeScript frontend, application workflows, reminders, contacts, document metadata, activity history, search/filtering, and a basic dashboard.

Generated data, benchmark coverage, and expanded analytics validation remain planned.
## Future Plans

The following areas are intentionally deferred beyond the current single-user version:

- Authentication and authorization
- Employer-facing workflows
- File upload and file storage
- Calendar integrations
- Production deployment wiring
## Documentation

Planned documentation:

- `QUICKSTART.md` for setup and run commands
- `ARCHITECTURE.md` for architecture notes
- `docs/adr/` for architecture decision records

## AI Assistance Disclosure

ChatGPT was used during development as a learning, design, and review assistant. It helped with project planning, architecture discussion, documentation drafting, implementation review, and test-suite design. Some tests may be generated with ChatGPT assistance and then reviewed, adapted, and validated as part of this repository.

## License

This project is licensed under the MIT License.

