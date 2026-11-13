## Project Overview

This is a distributed Supertext application using WebSockets to synchronize the
individual patron screens with the live performance. The application is
intended to stay as local first as possible, using a self-managed WiFi connection
through a Raspberry Pi which also runs the server.

## Set up Commands

- Get dependencies: `go get`
- Run code: `go run .`

## Code Style
- `Gofmt` is the standard format of the codebase

## PR Instructions
- Title follow scoped commits or the feature ticket title
- Always run `go fmt` and `go test` before committing.
- When changing major sections, always update relevant documentation in the
  `docs/` folder. Make sure to follow instructions on ADRs when changing documents.

## ADR Instructions

Each ADR can be at max 2 pages long, and must contain the following information:
- Title: a short noun phrase naming the decision, with its number. EX: “ADR-004: PostgreSQL for the main database”.
- Status: where the decision stands in its lifecycle (below).
- Context: the forces at play, whether technical, political, social, or particular to the project, stated as facts.
- Options considered: the alternatives you evaluated, each with the reason it lost.
- Decision: what you’ll do and why it wins, in full sentences and the active voice (“We will use…”).
- Consequences: everything that follows from the decision, the costs as well as the benefits.

Statuses must be one of the following:
- Proposed: written and under discussion; the team hasn’t agreed yet.
- Accepted: the team agreed, and the decision is in force.
- Rejected: the team considered the proposal and said no. Keep the file and the reason, so the same debate doesn’t start again.
- Deprecated: the decision no longer applies. Use deprecated only when nothing replaces it.
- Superseded: a later ADR replaced this one. Name the replacement and link it, and have the new ADR link back.

A changelog is also required at the end to keep a log of how the document has
changed over time.
