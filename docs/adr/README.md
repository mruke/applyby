## Preliminary Technical Direction

ApplyBy is currently evaluating an initial stack that supports a performance-conscious job application CRM.

The proposed direction is:

- Go for the backend service
- PostgreSQL for persistence
- React for the frontend UI
- TypeScript for frontend implementation
- A layered testing strategy with unit, integration, end-to-end, and helper test areas

These decisions are not final until the related ADRs are accepted.

See:

- `docs/adr/0001-select-go-for-backend-service.md`
- `docs/adr/0002-select-postgresql-for-persistence.md`
- `docs/adr/0003-select-react-for-frontend-ui.md`
- `docs/adr/0004-select-typescript-for-frontend-implementation.md`
- `docs/adr/0005-select-layered-testing-strategy.md`