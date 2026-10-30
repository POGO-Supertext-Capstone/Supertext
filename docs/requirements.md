# Requirements

## Problem statement
The Portland Opera's main venue has a supertext display. This display is used to display translated lyrics in real time via a slideshow so that attendees can understand the story. POGO (Portland Opera to Go) wants a similar feature for their traveling performances, where people can connect and view the supertext in real time, possibly in remote areas without internet access.

## Functional requirements
### FR-1: Browser Based Attendee Access

An attendee can access the active performance's supertext through a web browser on their personal device.

**Requirement Classification:** Must (Explicit requirement by project partner)

**Acceptance Criteria:**
1. Given an active supertext session, when an attendee navigates to the session's provided URL, then the browser displays the current supertext slide.
2. The attendee does not need to install an application or browser extension to view supertext.

### FR-2: Real-Time Supertext Synchronization

An attendee's browser automatically displays the same supertext slide as the show runner's display.

**Requirement Classification:** Must (Explicit requirement by project partner)

**Acceptance Criteria:**
1. Given an attendee is connected to an active session, when the show runner advances to another slide, then the attendee's browser displays the new slide without requiring a page refresh.
2. Given an attendee joins an active session after the performance has started, when the connection is established, then the attendee sees the currently active slide.

### FR-3: Local Network Connectivity

The system allows attendees to receive supertext through a local wireless network provided by a Raspberry Pi without an external internet connection.

**Requirement Classification:** Must (Required to support traveling performances without reliable internet)

**Acceptance Criteria:**
1. Given the attendee is connected to the Raspberry Pi hosting an active supertext session, when the attendee navigates to the provided session URL without external internet access, then the attendee can access the supertext.
2. When the show runner advances a slide while without external internet access, then the attendee's browser receives the updated slide.

### FR-4: Concurrent Attendee Connections

The system allows multiple attendees to access the same active supertext session simultaneously.

**Requirement Classification:** Must (Required for multiple attendees to use the system during a performance)

**Acceptance Criteria:**
1. Given an active supertext session, when two or more attendee devices connect to the session, then all connected devices display the same current slide.
2. When the show runner advances a slide, then all connected attendee devices receive the updated slide.

### FR-5: Accessibility Features

The system interacts with the attendee's phone's accessibility features for accessibility devices

**Requirement Classification:** Must (Explicit requirement by project partner)

**Acceptance Criteria:**
1. Given an attendee is connected to an active supertext session, when the slide changes on their phone, then the phone's accessibility features can recognize that change.
2. Given an attendee is using a braille device, when the slide changes the braille device should display the slide content.

### FR-6: SSL/TLS Certificates

The system should be able to provide a HTTPS connection to attendees

**Requirement Classification:** Must (Explicit requirement by project partner to avoid scare screens)

**Acceptance Criteria:**
1. Given the Raspberry Pi is connected to the internet, then the system should update the SSL/TLS certificates automatically.
2. Given the Raspberry Pi isn't able to update the SSL/TLS certificates before they expire, the show runners should be made aware so that they can connect it to the internet.

### FR-7: Presentation Enhancement

The system should be able to display presentations including different colored text, and animations

**Requirement Classification** Could (Not yet confirmed with project partner)

**Acceptance Criteria:**
1. Given a slide contains some form of animation, when the show runner switches to that slide, then the animation should play for the attendee.

## Non-functional requirements
- **NFR-1:** Compatibility: The attendee interface works in the latest stable versions of Chrome, Safari, and Firefox on mobile.
- **NFR-2:** Responsive Design: The supertext remains visible without horizontal scrolling on supported smartphone and tablet screen sizes.
- **NFR-3:** Performance: The attendee interface displays a slide update within 250ms (to be confirmed) after the show runner advances the slide.
- **NFR-4:** Accessibility: The attendee interface provides sufficient text/background contrast to meet WCAG 2.2 AA requirements.

## Constraints
- C-1: Attendees access supertext through a web browser rather than a dedicated mobile application.
- C-2: The system must be capable of operating without an external internet connection.
- C-3: Supertext server must run on a Raspberry Pi.

## Success metrics
- The server can handle 250 concurrent attendees receiving slides.
- Show runners can go through the entire process from upload to present without show-stopping errors
- Attendees can receive real time supertext without any display desynchronization or crashing

## Non-goals
- Developing a native iOS or Android application.
- Replacing or redesigning the Portland Opera House's existing supertext presentation system.
- Providing translated lyrics independently of the operator's existing supertext content.
- Allowing attendees to control or advance the supertext slideshow.

## Open questions
- What is the acceptable delay between the operator advancing a slide and attendee devices displaying it?
