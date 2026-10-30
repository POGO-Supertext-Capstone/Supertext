# ADR-004: Key Server

### Submitters
- Lukas Werner

## Status: Proposed

## Context

This ADR lays out the requirements of providing SSL/TLS certificates for the
embedded Raspberry Pi which provides both the network and application services
on the network edge. Modern mobile devices require encrypted HTTP connections
(HTTPS) in order to allow the browser to establish a connection without showing
a large screen warning, an explicitly stated anti-goal of the project.
To avoid having such an HTTP scare screen we need to get an SSL/TLS certificate
that will allow the user's device to trust that the connection is secure.
Getting a valid certificate presents a large challenge. This ADR provides a
description of the various options in creating such an HTTPS connection.

## Decision

To provide and access such an HTTPS certificate, we will be running a cloud
service which is responsible for creating and maintaining a good relationship
with Let's Encrypt (and potentially other [ACME compatible CAs](https://datatracker.ietf.org/doc/html/rfc8555)).
Then when the edge device (Raspberry Pi) needs to update its certificate it will
make an authenticated request to the Key Server and pull the private key needed
for it to be able to operate under the same domain as the Key Server.

### Alternatives Considered

An alternative that was considered was requiring an active network connection in
order to connect through using Cloudflare Tunnels or similar technologies, however
that introduces additional latency in an already potentially poor network
environment that would be unwarranted and unnecessary.

Another option that was considered was only providing HTTP service, however
given the trend of everything requiring an HTTPS connection, this was deemed
untenable in the long term. This leaves the last option, self-signed certificates,
which due to the design requirement of avoiding scare screens was also not viable.

## Consequences

### Positive
- Valid HTTPS connection during local edge operation
- Removal of HTTP scare screens on patron's personal devices

### Negative
- Requires synchronization with Key Server every 2–3 months due to shortened SSL/TLS certificate durations
- Requiring the server (Raspberry Pi) to terminate the SSL/TLS encryption may cause extra compute cycles spent on encryption that could otherwise be spent sending out slide change packets or handling extra users

## Change Log
- V1 (10/23/2026) Initial ADR [Lukas]
