# RFC: Administration page

| Field	         |  Value         |
| ---------------| ---------------|
| Author	       |  Jacob Hampton |
| Date	         | 2026-10-25     |
| Status	       | Proposed       |
| Decision date  | TBD            |
| Deciders       | TBD            |
| Superseded By  | NA             |

## Problem and Context: 
This RFC outlines the architectural and design requirements for the administration page and management interface of the Supertext application. Operating on a local Raspberry Pi 5 captive Wi-Fi network during live opera performances, the system requires a robust, dedicated portal for the showrunner. This portal must handle the initial performance script ingestion (.pptx, .txt), render visual verification views so showrunners can review slide sequencing prior to showtime, support in-place text corrections during pre-show rehearsals without requiring full file re-uploads, and integrate physical keyboard hotkeys for real-time prompter control.

Choosing the core architectural pattern for the administration page. Specifically, whether to build a decoupled single-page application, rely on constructed standalone endpoints, or implement a unified server-rendered module using Go and HTMX is an important foundational decision. It dictates our backend routing structure, authentication boundaries, deployment topology, and frontend maintenance overhead. If the team commits to a complex decoupled framework or unstructured custom routes mid-project, refactoring the frontend build toolchain and API contracts on an embedded edge device would cost significantly more than establishing the correct architectural boundary now. 

## Options:

### Option A: A unified server-rendered Go backend Module with HTMX and protected route handlers
Build the administration interface as a protected server-rendered module within the core Go binary, utilizing Go’s native packaging combined with HTMX for dynamic form submissions, zero-reload inline slide edits, and server-sent updates. Access is secured via local network isolation and lightweight operator authentication checks within the Go router. 

- **Benefits:** Keeps the entire application within a single Go binary with zero external service dependencies or complex Node.js toolchains. Delivers quick DOM updates via HTMX and server-rendered HTML fragments, preserving memory on the Raspberry Pi 5.
- **Costs:** Requires writing server-rendered HTML templates and managing DOM updates via HTML fragments rather than leveraging a heavy client-side component framework.
- **Right choice when:"** The application is deployed on an edge device that must operate entirely offline minimal resource overhead and zero external service dependencies.
  
### Option B: Decoupled single-page application (SPA) using react/TypeScript served through Go
Develop a standalone single-page application using React, Vite, and Tailwind CSS, communicating with the Go backend exclusively through a JSON REST API. Static build assets are embedded into the Go binary or served from disk. 

- **Benefits:** provides a rich component ecosystem and advanced client-side state management capabilities. 
- **Costs:** introduces a heavy Node.js build toolchain, increases binary build complexity, and requires managing client-side routing, state hydration, and separate API authentication tokens on an embedded device.
- **Right choice when:"** The administrative dashboard requires complex data visualization, multi-user real-time collaborative editing, or rich client-side state management that exceeds server-rendered capabilities

### Option C: Do nothing
Rely on basic, unauthenticated static HTML forms or manual command-line database scripts run via SSH on the Raspberry Pi 5 to manage shows, upload scripts, and edit text records. 

- **Benefits:** Zero initial development overhead for UI routing.
- **Costs:** Severe security vulnerabilities, high operational friction during live rehearsals, zero visual verification tools for showrunners, and extreme risk of human error during showtime.
- **Right choice when:"** The applicatoin is an internal developer utility with a single technical user and zero live performance time contraints

### Option D: Static file re-uploading for every edit
Require showrunners to re-upload the entire source presentation file every time a single type or line correction is needed during rehearsal. 

- **Benefits:** avoids building an in-place API or database mutation endpoints for individual slides
- **Costs:** Extremely slow and error-prone during live opera environments where quick fixes are needed immediately, risking script desynchronization.
- **Right choice when:"** Presentations are strictly immutable and rehearsal turnaround time is not relevant. 

## Analysis

Meeting the live rehearsal needs can be defined as an approach that satisfies the operational requirements during active pre-show rehearsals and performances on the local Raspberry Pi 5 network. The requirements are as follows:
- **Low-latency in-place editing:** allows showrunners to make immediate text corrections and adjustments through the management interface without requiring full file re-uploads or server restarts. 
- **Instant visual verification:** Renders structured cue sequences immediately upon ingestion or editing, giving the showrunner visual confirmation of formatting and layout before patrons connect.
- **Hight fault tolerance:** Prevent malformed inputs, accidental syntax errors, or unexpected edge cases from crashing the application server or corrupting active database state.
- **Speed and responsiveness:** Delivers near instant UI feedback and seamless keyboard-driven hotkey integration so operators can navigate scripts with zero perceptible lag.

| Concern                       | A: Unified Go and HTMX | B: Decoupled SPA	   | C: Do nothing	 | D: Static re-upload |
| ----------------------------- | ---------------------- | ------------------- | --------------- | ------------------- |
| Edge resource                 | Low (single Go binary) | High(Node/Vite)     | Low             | Low                 |
| Development complexity        | Low to moderate        | High (dual stack)   | Low             | Low                 |
| Real-time inline editing      | Yes (HTMX/SSE)         | Yes (REST/State)    | No              | No                  |
| Security and route protection | Native Go middleware   | Token/CORS overhead | None            | None                |
| Meets live rehearsal needs    | Yes                    | Yes                 | No              | No                  |

Options C and D fail critical operational and usability requirements during live opera rehearsals, and option B introduces unnecessary architectural complexity and resource overhead for an embedded edge device. Option A satisfies every performance, security, and deployment constraint while maintaining a lightweight footprint on the Raspberry Pi 5. 

Uncertainty:
- The exact performance threshold of concurrent server-rendered HTML fragment swaps under peak operator usage has not been benchmarked under live network load, though Go’s performance is well-documented

## Recommendation and verification plan
Adopt option A. Build the administration interface as a protected server-rendered module within the Go backend using HTMX for dynamic interactions and server-sent updates. 
### Acceptance criteria:
-	Administrative route handlers reject unauthenticated requests or requests originating outside the local management interface. Checked by automated unit tests in CI. 
-	Script file uploads successfully parse into structured records, persist within the database engine, and render accurately in the admin visual verification view. Checked by integration tests. 
-	On-device text edits correctly update database records and trigger cache invalidation or timestamp tracking for active broadcast streams. Checked before deployment. 
-	Repository CI workflows (lint, format, and run tests) pass successfully on every pull request.

### Constraints
-	Licensing and IP: The administrative portal must operate entirely offline without external SaaS dependencies (e.g. cloud authentication providers or external UI component libraries) to protect Portland Opera script confidentiality and intellectual property. All Go and HTMX dependencies use permissive open-source licenses. 
-	Hardware limitations: The portal must execute efficiently within the strict CPU and RAM bounds of the Raspberry Pi 5 captive portal without exhausting system resources.
-	Single-binary deployment: All administrative routing, HTML templating, and backend logic must compile into a single executable binary for effortless deployment onto a target hardware.

## Delegation and validation
-	What you handed AI tools: The initial translation of ADR-005 into this RFC structure. 
-	What checks made that reasonable: Our repository enforces automated validation on every push and pull request to the main branch through GitHub workflows (lint and format running go fmt, go vet, and go mod tidy; Run Tests running go build and go test -v).  However, neither of these covers documentation, so the check on this document is human review in the pull request and running it through Grammarly for grammatical correctness. 
-	Where you stopped and validated the result yourself, and why those points: TODO (author). Since the analysis above has not been backed by a prototype or benchmark yet, these claims must be personally verified.

## Decision:
Proposed
