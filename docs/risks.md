# Project Risk Register

## Active Risks 

### R-1: Key Server Certificate Synchronization and Edge Validity
-	**Category:** Technical
-	**Statement:** Given that the Raspberry Pi 5 operates on an isolated local edge network without constant internet connectivity, there is a possibility of SSL/TLS certificate expiration or synchronization failure with the cloud key server, leading to browser HTTPS security warnings on patrons' devices during live performances.
-	**Likelihood:** Medium
-	**Impact:** High
-	**Reasoning:** Modern mobile browsers strictly enforce valid HTTPS Certificates. If the synchronization interval fails or the key server is unreachable before a show, patrons cannot connect without friction.  
-	**Owner:** Lukas Werner
-	**Response:** Mitigate through scheduled pre-show sync checks and cached private key retention.
-	**Trigger:** Certificate validity drops below 14 days remaining without a successful synchronization configuration from the key server. 

### R-2: Hardware Storage Degradation and Database Corruption
-	**Category:** Technical
-	**Statement:** Given that the application executes on a Raspberry Pi 5 with a 4 GB microSD card for state persistence via BBolt, there is a possibility of high disk write amplification from rehearsals and key-value bucket updates, leading to file system corruption or catastrophic storage failure during a live performance.
-	**Impact:** High
-	**Likelihood:** Medium
-	**Reasoning:** MicroSD cards have limited write cycles. A corrupted BBolt data store during a live performance halts script delivery across all connected devices.
-	**Owner:** Jamie Murphy
-	**Response:** Mitigate through WAL-style key-value transaction safety and automated pre-rehearsal backups.
-	**Trigger:** Database write latency exceeds 200ms, or `fsck` reports file system errors during system service startup. 

### R-3: PowerPoint Shape Hierarchy and Text Extraction Edge Cases
-	**Category:** Technical
-	**Statement:** Given that POGO performance scripts will be supplied as PowerPoint files, there is a possibility of complex layout structures, non-standard text frames, or grouped text elements within individual slides, leading to text extraction failures or missing lines during automated ingestion.
-	**Likelihood:** Low
-	**Impact:** Medium
-	**Reasoning:** While standardizing on .pptx and implementing run-level decoding covers standard text boxes, complex PowerPoint templates provided by the opera may contain grouped elements or multi-column layouts that require careful validation against the parser’s shape iteration logic. 
-	**Owner:** Tom Sharp
-	**Response:** Mitigate through targeted test assertions using `.pptx` production files in the test suite.
-	**Trigger:** An incoming production script upload drops text content or fails to map slide lines correctly during staging verification. 

### R-4: High Concurrency Network Latency Under Peak Audience Load
-	**Category:** Technical
-	**Statement:** Given that potentially hundreds of concurrent patron devices will connect simultaneously to request the initial presentation payload via HTMX and receive live slide transitions via server-sent events (SSE), there is a possibility of network congestion or connection dropping on the Raspberry Pi 5, leading to desynchronized supertext display and noticeable latency between scene transitions. 
-	**Likelihood:** Medium
-	**Impact:** High
-	**Reasoning:** Real-time opera translation requires near-instantaneous slide delivery; latency higher than 1-2 seconds ruins live synchronization.
-	**Owner:** Jamie Murphy
-	**Response:** Mitigate through load testing stress benchmarks on the SSE dispatch loop and client-side caching.
-	**Trigger:** Benchmark stress testing shows SSE event dispatch latency exceeding 500ms at 150 simulated concurrent client connections. 

### R-5: Deployment Knowledge Concentration on Embedded Pi Hardware
-	**Category:** Team
-	**Statement:** Given that local server deployment, network binding, and system process management are configured primarily on a single team member’s hardware environment, there is a possibility of deployment bottlenecks if that individual is unavailable, leading to broken CI/CD deployments or blocked team integration testing.
-	**Likelihood:** Low
-	**Impact:** Medium
-	**Reasoning:** If only one person understands how to compile the Go binary, manage BBolt data paths, and configure the Raspberry Pi 5, edge hardware testing stalls when schedules diverge.
-	**Owner:** Jacob Hampton
-	**Response:** Mitigate through documentation, standardized build/test scripts, and shared setup instructions. 
-	**Trigger:** Team members are unable to deploy and run the application executable locally on a clean clone without direct assistance. 
