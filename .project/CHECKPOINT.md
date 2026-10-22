# Course refactor checkpoint — 2026-10-20

Status: completed third bounded phase (HTTP/API/relations/interfaces). Original user request remains phased course improvement. No overall goal, subagents, root actions, staging/resets/commits/pushes or behavior test runs. Preserve the user's heavily pre-staged index.

## Stable contract

- Russian warm «ты»; small Python/QA background.24 weekly meetings,50 min +10 reserve. Employment as first Go junior/junior+ is the core goal.
- Original LOGO.png,500 original interview questions, legacy/2026-10-19 and ignored tmp/FRAMING.md preserved. Advanced distributed track is separate;128-byte offsets/BLAKE3 unchanged. No Kafka/Redis prerequisite in core.
- Public contracts in pkg: types, enums, constants, errors, DTOs, interfaces. pkg/ports.go for repository. Service/port operations (ctx,req)->(res,err), object results including list wrapper; api.New takes pkg.TaskService. Private helpers may return scalars/slices. gRPC still needs protobuf/adapter.
- Tests must not be added or run unless requested. Existing assertions may be maintained for changed contracts; do not imply they passed. Ordinary builds/rendering/diff checks are allowed.

## Completed phases

1. Rewrote student/mentor/combat01–04, first principles/onboarding, deep Git labs, project map, six SVGs, examples/hexagonal. START_HERE uses winget/MSI/snap and go mod init.
2. Rewrote05–09 and12/14/21 with HW/mentor, eight SVGs.09 is ENVIRONMENT before SQL;21 revisits releases/readiness/logs/shutdown. SQL12–13 uses jobflow_lab; app uses jobflow. Dockerfile/compose.yaml/env files in reference09–24; removed competing old compose21–24. Core14–24 pgxpool/v5.11.0, Goose/v3.29.0 embedded SQL, Provider + advisory lock, same-pool stdlib bridge. Migration before recovery/workers/HTTP, startup failure nonzero, no separate migrator/runtime SQL files. Legacy course_migrations rejected preserving data; docs/POSTGRES_UPGRADE.md gives explicit transfer.
3. Rewrote10/11/13/15 student/mentor/combat, four SVGs, docs/HTTP.md, targeted architecture/14 notes, metadata/navigation, interview route (not original bank). HW2–3h /4–6h with scaffold /3–5h /3–5h.
   -10: client/server/request/response, query hello, about/health, container address, connection vs HTTP failures, HEAD.
   -11: supplied pkg/app/memory/Decode scaffold, JSON→DTO→service→port, exact validation/error guarantees, PUT replacement, state invariants.
   -13:3 tasks/3 events, JOIN3 vs LEFT4, ON vs WHERE, COUNT, CASCADE/audit, integrity vs atomicity. reference13/sql/relations.sql is rollback lab; sequence values not rolled back.
   -15: revisit existing interfaces, consumer-driven contracts, method sets/compile assertions vs behavior, composition, two graphs, interface cost.
   - Task models/errors now pkg in reference06–24. Object TaskService/TaskRepository in11–24. API depends only on pkg; app owns validation, memory/postgres implement outgoing port in repository.go; storage helpers private. Insert differs from inbound Create. HTTP v1 still flat Task and bare list array from res.Tasks. Health DTO in pkg. Job20+ remains internal and explicitly awaits its phase.
   - Existing checks/reference assertions11–24 updated to public object port, no new scenarios. Existing API scenarios11–14 build app.New(store). No test suite was run or test-file compilation claimed.

## Current evidence and limits

- This phase: all15 reference10–24 modules go build ./... succeeded, Go1.27.2. Earlier phase built06–24 and first01–04/hex example.
- Docker images built this phase11/24: go-course-local:part3-lesson11 /part3-lesson24. Logs /tmp/go-course-part3-image-{11,24}.log. Prior images09/24 go-course-local:lesson9 /lesson24.
- Compose config --quiet passed11/15/24, previously09/14/24. No containers started in these phases.
- Four new SVGs rendered/viewed in /tmp/go-course-part3-preview; prior previews /tmp/go-course-part2-preview.18 new core posters total including Git/hex.
- Navigation generator only, no full test/validation tools. git diff --check clean. No unstaged changes to preserved logo/bank/archive/input. docs/VALIDATION.md separates current builds from historical full-suite report.
- No behavior/integration/SQL runtime test claims. Native build excludes _test.go.

## Next bounded phase

Deepen16 transactions,17 goroutines,18 workers,19 context,20 background lifecycle with the same humane progression and meaningful practice/SVGs. Migrate Job20+ public types/DTOs/interfaces into pkg and object contracts during that phase; keep HTTP wire format, same embedded Goose and startup order.22–24 performance/employment prose/HW still await review. Do not redo completed01–15, rename24 weeks or inflate novice homework.

Tools: /tmp/jobflow-course-tools/go/bin/go (1.27.2); modules go1.27.0. No root install needed. Native librsvg/Cairo available for SVG rendering. Authored temporary scripts live /tmp; do not blindly rerun mutation scripts (some are not idempotent).
