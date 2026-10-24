# ADR-003: Frontend Conventions 

### Submitters
- Tom Sharp

## Status: Proposed

## Context

This ADR lays out how the frontend (or client side) for Supertext should be handled. 

This ADR focuses on the frontend experience for audience members, and not on the admin / showrunner pages. 

## Decision

Supertext will use a combination of HTMX and pure JavaScript to handle frontend interactions.

### HTMX
HTMX will be used to connect the client to the server. Once a client connects to the server via the endpoint `/vx/slides`, they will receive the full presentation in a presentable format.

An example for getting the slides using HTMX could be:

```HTML
<div id="slide-loader"
  hx-get="/vx/slides"
  hx-trigger="load"
  hx-swap="none">
</div>
```

These slides can then be read in by JS for later presentation.

It is the server's responsibility to decode the slides into before they are sent to the client. Slides will be sent as a JSON document with limited HTML tags to supprt text types. 

The current allowed tags are:
- `p`: for seperating lines 
- `b`: for bolded text
- `i` for italisized text

Although this is subject to change. An example of a `slides` object would be:

```json
{
  "slides": [
    {
      "slideNum": 1,
      "data": "<p>This is slide 1</p>"
    },
    {
      "slideNum": 2,
      "data": "<p>This is slide 2</p><p>with <b>multiple</b> lines</p>"
    }
  ]
}
```

### JavaScript
For sending supertext to the clients, [Server-Sent Events (SSE)](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events) will be used to trigger slide transitions.

The server will send out slide numbers via the endpoint `/vx/present`. When a slide is to be changed, the server should send the event: 

```
event: show-slide
data: #
```

Where the `data` section holds the slide number to present. 

### Alternatives

One alternative discussed was using pure HTMX for the frontend. This was eventually decided against as we wanted the flexibility to control how slide data is sent, stored, and presented via server sent events. 

The other alternative discussed was using strictly HTML and JavaScript. This was also decided against as we felt that HTMX provided some features that could be helpful / make our lives easier.

## Consequences

### Positive
- **Separation of concerns:** HTMX will handle the initial retrieval of slide data, while JavaScript will handle slide parsing, rendering, and live presentation updates.
- **Reduced server communication during presentations:** The client can retrieve the presentation once and use the locally stored slide data as the server sends slide-change events, rather than requesting the full slide content for every change.
- **Flexible presentation logic:** JavaScript gives the client greater control over how slides are stored, rendered, and transitioned between.
- **Scalability:** By loading presentation content once and distributing small slide-change events, the architecture can reduce repeated content requests and the amount of data sent during a presentation. This can improve efficiency as the number of audience members grows.

### Negative
- **Initial loading dependency:** The client must successfully retrieve the presentation before it can display slides. Loading failures or malformed slide data may prevent the presentation from working correctly.
- **Connection and recovery handling:** SSE connections can be interrupted. The frontend must account for reconnection behavior and ensure that clients return to the current slide after reconnecting.

### Neutral
- **Hybrid frontend architecture:** The application will use both HTMX and custom JavaScript rather than relying exclusively on either approach.
- **Server-driven presentation state:** The server remains responsible for determining which slide should be shown, while the client is responsible for displaying it.

## Change Log
- V1 (10/22/2026) Initial ADR [Tom]
