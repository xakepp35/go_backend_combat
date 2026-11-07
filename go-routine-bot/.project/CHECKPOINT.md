# go-routine-bot checkpoint

## Current objective and decisions

Documentation package for go-routine-bot: README, VISION, MVP_SPEC, PILOT, PLAN, ROADMAP and AGENTS, with this checkpoint. Application implementation remains future work. No personal names or named invitation files belong in the repository. Invitations may be prepared privately. Do not recreate the removed message draft.

Pilot: one mentor and initial participants; 2–3 additional learners may join. Composition remains tentative. Three goal-based routes: start in backend / transition to Go / systems design. The bot may be a participant's main Go learning project. Two 45–60 minute sessions plus a review are proposed, not confirmed availability.

One Go app, PostgreSQL/pgxpool v5, embedded Goose v3 and Compose; public object contracts in pkg. Human acceptance, source-grounded AI hints. Course references use actual large chapters and stable IDs rather than assuming a course/lessons.json file exists.

## Latest change

Added PLAN.md: service-contract-first implementation sequence aligned with mentor/0x_PROJECT_FLOW.md, stages 0–8, mentor/participant responsibilities, first three work sessions, schema assessment, transaction guarantees and embedded migration policy. Synchronized README, MVP_SPEC, PILOT, ROADMAP and AGENTS. Schema review retained content versioning/FKs, replaced current-pointer-only progress with enrollment/assignment history, separated submission revisions from operation keys, and clarified rhythm, quota and notification guarantees. Route plan versions are frozen per enrollment; pause does not permit a second unfinished enrollment.

Next implementation starts with one scenario and public GetDashboard/StartStep contracts, not full-schema SQL. No application or migration was created/applied, and no tests or external actions were performed. Personal invitation draft remains absent. Preserve the user's existing .gitignore edit and unrelated course files.
