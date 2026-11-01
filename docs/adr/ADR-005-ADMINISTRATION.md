# ADR-005: Administration Page

### Proposers: 
- Jacob Hampton <hamptjac@oregonstate.edu>

## Status: Proposed

## Context:
This ADR defines the architectural pattern and operational requirements for the administration page.

The supertext system requires a dedicated management interface for showrunners and operators running on the local Raspberry Pi 5 network. The admin portal must handle initial slideshow file ingestion, visual verification, live script adjustments without full re-uploads, and keyboard-driven show progression. 

## Decision:
The administration page will be built as a protected server-rendered module with the Go backend, utilizing HTMX and pure JS for dynamic UI updates and server-sent updates, specifically tailored for local browser access by the showrunner. The administration routes are protected via local network isolation and session/password check to restrict access from patron devices on the public Wi-Fi network. 

### Core architectural decisions:
- **Unified Go backend Route:** Admin route handlers will reside in the core Go binary running on the Raspberry Pi 5, ensuring zero extra external service dependencies.

- **File ingestion and visual verification pipeline:** The portal will accept initial slideshow file formats (e.g., .pptx, .txt), parse raw files into structured application records, and save those records within the database engine ([ADR-002])(https://github.com/POGO-Supertext-Capstone/Supertext/blob/main/docs/adr/ADR-002-DATABASE.md), rendering a visual verification view so the showrunner can review slide sequencing before showtime
  
- **In-place slide editing API:** The admin interface will support direct, on-device slide text editing via the UI so showrunners can make corrections during rehearsal or pre-show checks without needing to re-upload the entire source presentation file.
  
- **Keyboard-driven control integration:** the live prompter view will listen for keyboard inputs (arrow keys, spacebar, etc.) from the showrunner’s laptop connected to the local Pi network to advance or reverse cues. 

### Alternatives considered:
An alternative to this is static file re-uploading for every edit, but this was not chosen as it is too slow and error-prone during live opera environments where quick line or type fixes are needed immediately.

## Consequences:
### Positive:
- **Rapid iteration and verification:** showrunners can instantly verify uploaded slides and fix typos in real-time on the device.
- **Low overhead:** maintains our lightweight Go + HTMX footprint on the Raspberry Pi 5.
- **Modular documentation:** Aligns administrative database persistence, file parsing, and workflow requirements into a single cohesive architecture. 

### Negative:
- **State coordination:** Requires careful backend synchronization between BBolt database records, the active presentation cache, and active SSE broadcast channels when inline edits occur.

## Change Log
- V1 (10/24/2026) Initial ADR – Jacob Hampton <hamptjac@oregonstate.edu> 
