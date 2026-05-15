# 0009. Track Document Metadata Without File Storage

## Status

Accepted

## Context

ApplyBy tracks document records associated with job applications. These records represent metadata about resumes, cover letters, portfolios, or other job-search documents.

The current prototype does not need to upload, store, download, virus-scan, or permission-check actual files. Adding file storage would require additional architectural decisions around local paths, object storage, upload limits, security, backups, and future multi-user access.

The useful prototype behavior is knowing which document metadata was associated with an application, not managing binary files.

## Decision

ApplyBy will track document metadata only in the completed local prototype.

Document records may store descriptive information needed by the job-search workflow, but actual file upload and file storage are deferred.

## Consequences

Positive consequences:

- The prototype supports useful document tracking without file-storage complexity.
- The data model remains relational and straightforward.
- The app avoids premature decisions about local filesystem paths or object storage.
- Future hosted deployment is not constrained by an early file-storage choice.

Negative consequences:

- Users cannot upload or retrieve files from ApplyBy.
- ApplyBy cannot guarantee that an external file still exists.
- File versioning and file download workflows are not available.

Follow-up work if file storage is added:

- Decide between local filesystem storage and object storage.
- Add upload and download APIs.
- Add file size/type validation.
- Add security and retention rules.
- Add tests for file persistence and cleanup behavior.