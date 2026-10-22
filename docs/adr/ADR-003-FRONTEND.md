# ADR-0013: Fronend Conventions 

### Submitters
- Tom Sharp

## Status: Proposed

## Context

This ADR lays out how the frontend (or client side) for Supertext should be handled. 

## Decision

Supertext will use ...

For handling sending lyrics to the clients, Server-Sent Events (SSE) will be used to trigger the changing of slides. Once a user connects to the server, they will recieve the full presentation in a markdown format, and the server side will send the slide to switch to in plain text to the clients. The clients will then switch to that slide.

### Alternatives
