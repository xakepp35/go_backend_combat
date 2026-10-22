# Go Backend Combat

## Approved course contract

- Rebuild the course as 24 weekly lessons, six months, Russian language.
- Students may have only two months of Python/QA exposure. Explain programming, algorithms, functions and tools from first principles; do not assume fluent command of another language.
- Each lesson has 50 minutes of planned work and 10 minutes of reserve; final mock uses the last 10 minutes for feedback.
- Homework estimates follow the actual work, including reading/debugging. Revised lessons 01–04 use 1.5–2.5, 2–3, 2–4 and 3–4 hours; later estimates require review during their phase. Every increment develops JobFlow and has an optional creative ★ exercise. Never pad a tiny task to meet a fixed duration.
- Use warm Russian prose addressed to «ты», everyday examples, questions before answers, and occasional accurate developer humor. Introduce terms after their motivating example.
- The core goal is preparation for a first Go junior/junior+ job. Restore and retain the original logo, brand and all 500 interview questions.
- A mentor reviews every homework; open tests supplement explanation and review.
- Core: one Go application, HTTP, PostgreSQL, bounded background work. Distributed JobFlow is an advanced continuation, not a prerequisite for completing core.
- Keep student notes, compact mentor guides, homework, stage references and acceptance checks synchronized.
- Preserve the advanced 128-byte field offsets in `advanced/jobflow/docs/FRAMING.md`. Do not replace BLAKE3 with SHA-256.
- In advanced JobFlow, PostgreSQL owns durable state/bodies and Redis projections/cache/Streams are recoverable. Core requires no Redis or Kafka.
- `pub`, `sub`, `act`, `ret` belong to the advanced continuation. Only registered cooperative handlers execute.
- Handwritten Markdown is authoritative. Generators may update marked navigation blocks, never lesson prose, the selling README or the interview bank.
- Tests and reference solutions are open. A failing mandatory check must produce a nonzero exit; infrastructure omissions cannot count as full PASS.

## Current incremental refactor (2026-10-20)

- Revise in bounded parts. Completed phases cover 01–15 plus the early environment/release pass on21. Task contracts are migrated in11–24; later Job contracts and late lesson prose still need their own phase. Do not claim the whole course is finished.
- Introduce requirements, explicit invariants, boundary examples and a repeatable development procedure from lesson 01. Keep a warm voice and add meaningful Basic+/Middle+ experiments, tradeoffs and failure analysis.
- Each revised lesson needs a linked, readable SVG poster; repo-native SVG is edited as code. Preserve original posters and logo.
- Use normal automation: Windows winget/MSI, Ubuntu `sudo snap install go --classic`, project creation through `go mod init jobflow`; do not teach hand-written initial go.mod.
- New public contracts live in pkg: requests/responses, types/enums/constants, shared errors and interfaces. Repository port is separate in pkg/ports.go. API constructors accept the public service interface, not concrete app implementations.
- Public service operations use (ctx, req) -> (res, err) with object requests/responses, including list wrappers. Pure private helpers may return scalars. gRPC still requires protobuf and a transport adapter.
- docs/PROJECT_FLOW.md and examples/hexagonal define the new architecture convention. reference/06–24 Task models/errors live in pkg; reference/11–24 Task service and repository use object contracts. Job 20+ remains for its separate phase; label its scope honestly.
- Teach Git with separate disposable laboratories: branching/merge/fast-forward, conflict diagnosis, dirty worktree vs detached HEAD, reset recovery and local amend vs shared published history. Never run destructive demonstration commands in the user's course repository.

## Working rules

- Read `.project/CHECKPOINT.md` before resuming a long implementation phase; keep it concise and current.
- Preserve `legacy/2026-10-19` and its SHA256 manifest. It includes the owner's two SQL drafts.
- Leave `tmp/FRAMING.md` unchanged; it is the original ignored input.
- No root installation is needed. Isolate test infrastructure with a unique Compose project; do not modify unrelated containers.
- Use the owner's workstation permission rules from the session. Never ask for passwords in chat.
- Do not delegate to subagents unless the user authorizes it.

## Practical environment revision (2026-10-20)

- Docker/Compose basics are lesson 09 ENVIRONMENT, before HTTP10 and SQL12. Lesson21 OPERATIONS is the second pass on release, readiness, compatible schema evolution, logs and shutdown.
- PostgreSQL always runs in Docker; psql runs via compose exec. SQL12–13 experiments use jobflow_lab, app schema uses jobflow.
- Core14–24 use native pgxpool/v5 and embedded Goose/v3 migrations during application startup. No cmd/migrate, separate migrator service or runtime SQL directory. Goose uses stdlib.OpenDBFromPool for the same pool, Provider and session advisory locking.
- Apply migrations successfully before jobs recovery/workers/HTTP. Startup failure exits nonzero. Applied SQL is immutable, new changes get new versions; no automatic Down on stop/rollback.
- Legacy course_migrations is detected and rejected without deleting data; docs/POSTGRES_UPGRADE.md describes an explicit preserving transfer.
- reference09–24 have one compose.yaml each. Avoid reintroducing competing docker-compose.yml files or requiring host PostgreSQL installation.

## HTTP / Task contract revision

- Lessons 10, 11, 13, 15 are revised with their own SVGs. 10 introduces request/response and container HTTP; 11 receives a supplied app/DTO/memory/Decode scaffold, 15 revisits and explains its interfaces with two implementations.
- Core 11–24 TaskService/TaskRepository and req/res live in pkg; api.New accepts pkg.TaskService, app.New accepts pkg.TaskRepository. Object request/response for every public operation; outbound Insert differs from inbound Create. Small private storage helpers may return domain values/slices.
- HTTP v1 remains flat Task objects and a bare array for GET /tasks, explicitly mapped from ListTasksResponse.Tasks. Do not silently change the external protocol while refactoring service contracts.
- Existing memory/postgres assertions call the public object port. Preserve their scenarios; do not require private helper names in submitted projects. No new tests or test runs without a request.
- Intro 10 has no TaskService dependency because it serves static/query-only endpoints. Background Job 20+ is not yet migrated; do not claim all later domains are complete.
