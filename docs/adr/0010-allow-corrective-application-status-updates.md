# 0010. Allow Corrective Application Status Updates

## Status

Accepted

## Context

ApplyBy originally modeled application status as a directional lifecycle. That approach made some transitions invalid, such as moving from `rejected` back to `interviewing`.

Manual testing showed that this was too strict for a personal tracker. Users may enter the wrong status, receive new information, reopen an opportunity, or need to correct historical data. A tracker should support correcting state unless there is a strong business reason not to.

ApplyBy still needs to reject unsupported status values, but it does not need to enforce a one-way pipeline for valid statuses.

## Decision

ApplyBy will allow applications to move between any two different valid statuses.

The backend remains authoritative for status validation. Unknown statuses are rejected. No-op transitions from a status to the same status are rejected.

## Consequences

Positive consequences:

- Users can correct mistakes.
- The UI can present all valid statuses without conflicting with backend lifecycle rules.
- The status model better matches a personal job-search tracker.
- Tests can focus on valid status values rather than rigid pipeline movement.

Negative consequences:

- Status history may contain non-linear movement.
- The model is less suitable for enforcing a formal business process.
- Reporting must not assume that status movement is strictly forward.

Guidance:

- Preserve status history and activity events for visibility.
- Continue rejecting invalid status values.
- Do not reintroduce strict transition rules unless the product becomes a formal workflow system rather than a personal tracker.