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

The first version will focus on a single-user job-search workflow with local or simple hosted persistence. External integrations such as email inbox sync, calendar sync, browser extensions, AI resume rewriting, and multi-user collaboration are future possibilities, not initial requirements.

## Status

Initial repository setup and planning.

## Documentation

Planned documentation:

- `QUICKSTART.md` for setup and run commands
- `ARCHITECTURE.md` for architecture notes
- `docs/adr/` for architecture decision records

## AI Assistance Disclosure

ChatGPT was used during development as a learning, design, and review assistant. It helped with project planning, architecture discussion, documentation drafting, implementation review, and test-suite design. Some tests may be generated with ChatGPT assistance and then reviewed, adapted, and validated as part of this repository.

## License

This project is licensed under the MIT License.