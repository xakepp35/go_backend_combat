# go-routine-bot — stable project instructions

- This directory contains the agreed product documentation. Current deliverable is documentation, not a running bot.
- Name and directory are exactly `go-routine-bot`. Keep all project documentation here.
- Pilot: one mentor and initial participants; 2–3 additional learners may join. Roster is not confirmed. Telegram bot and Telegram Mini App.
- Tracks describe goals: backend start, transition to Go, systems design. Do not promise professional grades from course completion.
- For a learner transitioning to Go, the bot can be the primary learning project. Two 45–60 minute sessions plus a short review are a proposal, to be agreed with the participant.
- Use warm Russian addressed to «ты» in learning content; keep adult learners' engineering experience and available time in view.
- Human mentors accept required work and confirm skills. AI gives bounded, source-grounded hints and cannot write progress or accept work.
- Target implementation: one Go service, PostgreSQL/pgxpool v5, embedded Goose v3 at startup, Docker Compose, small HTML/JS Mini App.
- Public contracts/DTOs/types/enums/errors/interfaces belong in pkg; outbound ports in pkg/ports.go. Service operations use object `(ctx, req) -> (res, err)` contracts.
- Preserve the surrounding course, logo, interview bank and the user's existing Git changes. Current checkout has large course chapters and no course/lessons.json; use stable content IDs and versioned source references.
- Keep VISION.md, MVP_SPEC.md, PILOT.md, ROADMAP.md consistent. README.md is the entry point; .project/CHECKPOINT.md records actual progress and limits.
- No new tests or test execution without an explicit request. Acceptance scenarios in documentation are requirements, not evidence of passing behavior.
- Do not send personal invitations, create external accounts, deploy or spend money merely because the plan describes these future actions.
- Do not spawn subagents without explicit authorization.
