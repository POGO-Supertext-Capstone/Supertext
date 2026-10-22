# ADR-0013: Fronend Conventions 

### Submitters
- Tom Sharp

## Status: Proposed

## Context

This ADR lays out how the frontend (or client side) for Supertext should be handled. 

This ADR focuses on the front end experience for audience members, and not on the admin / showrunner pages. 

## Decision

Supertext will use a combination of HTMX and JavaScript to handle frontend interactions.

### HTMX
Once a user connects to the server, they will recieve the full presentation in a markdown format. HTMX will be used to connect the client to the server and recieve the slides. 

An example for getting the slides using HTMX could be:

```HTML
<div id="slide-loader"
  hx-get="/vx/slides"
  hx-trigger="load"
  hx-swap="none">
</div>
```

These slides can then be read in by JS for later use.

### JavaScript
For handling sending lyrics to the clients, Server-Sent Events (SSE) will be used to trigger the changing of slides. JS will also be responsibe for decoding the markdown format that the slides are saved in.

The server will send out slide numbers via the endpoint `/vx/present`. When a slide is to be changed, the server should send the event: 

```
event: show-slide
data: #
```

Where the `data` section holds the slide number to present. 

### Alternatives

One alternative discussed was using pure HTMX for the frontend. This was eventually decided against as we wanted the flexablity to control how slide data is sent, stored, and presented via server sent events. 

## Consequences

### Positive
- **Separation of concerns:** HTMX will handle the initial retrieval of slide data, while JavaScript will handle slide parsing, rendering, and live presentation updates.
- **Reduced server communication during presentations:** The client can retrieve the presentation once and use the locally stored slide data as the server sends slide-change events, rather than requesting the full slide content for every change.
- **Flexible presentation logic:** JavaScript gives the client greater control over how slides are decoded, stored, rendered, and transitioned between.

### Negative
- **Initial loading dependency:** The client must successfully retrieve and parse the presentation before it can display slides. Loading failures or malformed slide data may prevent the presentation from working correctly.
- **Connection and recovery handling:** SSE connections can be interrupted. The frontend must account for reconnection behavior and ensure that clients return to the current slide after reconnecting.
- **Parsing and rendering responsibility:** The frontend must correctly interpret the presentation's Markdown format. Changes to the format may require corresponding updates to the JavaScript parser and renderer.

### Neutral
- **Hybrid frontend architecture:** The application will use both HTMX and custom JavaScript rather than relying exclusively on either approach.
- **Server-driven presentation state:** The server remains responsible for determining which slide should be shown, while the client is responsible for displaying it.

## Change Log
- V1 (10/22/2026) Initial ADR [Tom]