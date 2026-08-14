# Contributing

ApplyBy is a solo portfolio project, but it`s developed in such a way that as a template I might re-use processes for future team-based project plans. Every change requires a branch, pull request and documentation update. CI must pass before anything merges into main. 

## Workflow 

1. Open or claim an issue before starting work. Every future change should trace back to some issue. 
2. Branch off of `main` using a prefix that describes change type"
    - `step/` For feature slice
    - `fix/` For bug fix
    - `patch/` For refactor
    - `docs/` For doc only changes
    - `qa/` For testing
    - `style/` for UI changes
3. Commit in small steps and describe only what changed
4. Open PR to `main` using `PULL_REQUEST_TEMPLATE.md` referencing the issue with `Closes #<number>`
5. CI must pass before merging 
6. Merge pull request by default matching branch commits 

## Important notes

### Architectural Decisions 

If a change introduces a new dependency, changes to persistence or API contract, or reverses an earlier decision, add or update an ADR and readme index. 

### Tests

- Backend: Table-driven tests alongside the code they cover (*_test.go). Repository interfaces get an in-memory fake for application-layer tests; Postgres-backed tests are integration tests gated behind `APPLYBY_INTEGRATION_TESTS=1`
- Frontend: component and page tests alongside the component (*.test.tsx), run with npm test.

### Scope discipline

This project deliberately avoids generic CRUD abstractions and speculative interfaces. Preference is given to explicit, readable code over abstraction even if it means some repetition (within reason, use your better judgement).

