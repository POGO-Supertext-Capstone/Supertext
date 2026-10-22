# RFC: Key Server and Edge Device SSL/TLS Encryption

| Field         | Value        |
| ------------- | ------------ |
| Author        | Lukas Werner |
| Date          | 2026-10-25   |
| Status        | Proposed     |
| Decision Date | TBD          |
| Deciders      | TBD          |
| Superseded By | N/A          |

## Problem and Context

This RFC lays out the requirements of providing SSL/TLS certificates for the
embedded Raspberry Pi which provides both the network and application services
on the network edge. Modern mobile devices require encrypted HTTP connections
(HTTPS) in order to allow the browser to establish a connection without showing
a large scare screen warning, an explicitly stated anti-goal of the project. To
avoid having such an HTTP scare screen we need to get an SSL/TLS certificate
that will allow the user's device to trust that the connection is secure.
Getting a valid certificate presents a large challenge: the Raspberry Pi runs on
a self-managed WiFi network that may have no internet access during a
performance, so it cannot count on obtaining a certificate when a patron
connects. This RFC describes the options for creating such an HTTPS connection.

## Options

### Option A: Self-Managed Key Server (ACME)

Run a cloud service, the Key Server, which is responsible for creating and
maintaining the relationship with Let's Encrypt (and potentially other
[ACME-compatible CAs](https://datatracker.ietf.org/doc/html/rfc8555)). When the
edge device (Raspberry Pi) needs to update its certificate, it makes an
authenticated request to the Key Server and pulls the certificate and private
key needed to operate under the same domain as the Key Server.

- **Benefits:** Valid HTTPS during local edge operation with no internet
  connection needed at show time. No HTTP scare screens on patrons' personal
  devices. Certificate issuance logic lives in one place rather than on every
  device.
- **Costs:** This server would be outside the Portland Opera's Control, leading
  to the fact that it would be one more cloud service to build, host, and
  secure. The edge device must synchronize with the Key Server every 40 days
  because of [shortened SSL/TLS certificate lifetimes](https://letsencrypt.org/2025/12/02/from-90-to-45).
  The private key leaves the Key Server, so the authenticated channel is a
  sensitive point. The Pi terminates TLS itself, which spends compute that
  could otherwise go to slide change packets or extra users.
- **Right choice when:** The edge device must work offline or on a poor network
  during the performance, and browsers must trust the connection without
  warnings.

### Option B: Cloudflare Tunnels (or Similar)

Require an active internet connection and route patron traffic through a tunnel
so the public endpoint presents a valid certificate.

- **Benefits:** No certificate handling on the device and no Key Server to
  operate.
- **Costs:** Adds latency in an already potentially poor network environment,
  and makes the show depend on external connectivity, which runs against the
  local-first goal of the project.
- **Right choice when:** The venue is guaranteed to have reliable internet and
  latency does not matter.

### Option C: Do Nothing / Use HTTP

Serve the application over plain HTTP.

- **Benefits:** No certificate infrastructure at all.
- **Costs:** Given the trend of everything requiring an HTTPS connection, this
  is untenable in the long term. Browsers already warn on or restrict plain HTTP
  and that gets stricter over time, so the cost arrives as scare screens or
  broken features on patron devices.
- **Right choice when:** The audience is a controlled set of devices that are
  configured to accept HTTP. That is not the case here.

### Option D: Self-Signed Certificates

The Pi generates its own certificate.

- **Benefits:** Works fully offline and costs nothing.
- **Costs:** No patron device trusts the certificate, so browsers show the scare
  screen, which violates the design requirement.
- **Right choice when:** Every client can be preconfigured to trust the
  certificate. Patrons use their own devices, so this does not apply.

## Analysis

| Concern                    | A: Key Server | B: Tunnel         | C: HTTP         | D: Self-signed   |
| -------------------------- | ------------- | ----------------- | --------------- | ---------------- |
| Trusted by user            | Yes           | Yes               | No encryption   | No               |
| Works offline              | Yes           | No                | Yes             | Yes              |
| Added latency              | None          | Extra hop         | None            | None             |
| Operations load            | Cloud service | Third-party agent | None            | None             |
| Meets no-scare-screen goal | Yes           | Yes               | No              | No               |

Options C and D fail a stated requirement (no scare screens), and Option B fails
the local-first design when connectivity is poor, which leaves Option A as the
only one that satisfies every requirement. Its main costs are operational: the
2–3 month renewal cycle, a service to run, and TLS termination load on the Pi.

Open questions and uncertainty:

- The TLS compute overhead on the Pi has not been measured. Modern ciphers on
  Raspberry Pi hardware are expected to be inexpensive relative to the
  WebSocket/SSE workload, but this should be confirmed with a load test.
- Distributing the private key to devices means a compromised device exposes the
  key for the shared domain. Per-device certificates or short-lived credentials
  may be worth evaluating.

## Recommendation and Verification Plan

Adopt Option A. Build the Key Server using CertMagic, which handles ACME
issuance and renewal, to obtain certificates from Let's Encrypt, and have the
edge device fetch its certificate over an authenticated request.

This wins the analysis because it is the only option that gives a trusted HTTPS
connection without depending on internet access during the show.

Acceptance criteria:

- A freshly provisioned Pi serves the app over HTTPS and a stock iOS and Android
  browser load it with no certificate warning. Checked before merging the first
  implementation.
- The Pi continues to serve valid HTTPS with the internet disconnected. Checked
  before merging.
- An unauthenticated request to the Key Server is rejected and returns no key
  material. Checked by an automated test in CI.
- The Pi renews its certificate before expiry (for example, at 30 days
  remaining), and a failed sync surfaces a visible warning. Checked with a
  simulated short-lifetime certificate in staging.
- A load test of slide-change broadcasts with TLS enabled shows acceptable
  latency at the target number of patrons. Threshold to be set before testing.

## Constraints

- **Licensing:** CertMagic is Apache-2.0 licensed, which is compatible with this
  project. Use of Let's Encrypt is subject to its Subscriber Agreement.
- **Domain and rate limits:** The Key Server requires a domain we control, and
  issuance is subject to Let's Encrypt rate limits, which affects how many
  devices can request certificates independently.
- **Data use:** The Key Server stores and serves private key material but no
  patron data. Its access credentials and the key store must be protected.
- **Partner confidentiality:** None apply to this decision.
- **Operations:** This adds a cloud service that must be hosted and maintained,
  which is a departure from the otherwise local-first deployment.

## Delegation and Validation

- **What you handed to AI tools:** The initial translation of ADR-004 into this
  RFC structure.
- **What checks made that reasonable:** This repo's CI runs `go fmt` and `go
  test` workflows, but neither covers documentation, so the check on this
  document is human review in the pull request. Additionally
  [Harper](https://writewithharper.com) was run for gramatical correctness.

## Decision

Proposed, until your team decides. Then accepted, rejected, or superseded, with
the date and who chaired the decision. If superseded, link the RFC that replaces
this one.
