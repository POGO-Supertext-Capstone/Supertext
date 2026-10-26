# RFC: Database Engine Selection

| Field | Value |
| :---: | :---: |
| Author | Jamie Murphy <murpjami@oregonstate.edu> |
| Date | 2026-10-25 |
| Status | Proposed |
| Decision Date | N/A |
| Deciders | N/A |

## Problem and Context

The Supertext program needs to be able to run off a Raspberry Pi with no remote connection, and be able to communicate with possibly hundreds of devices simultaneously. Therefore,
storing data in an efficient manner that has no dependencies on the open internet and can be used quickly and concurrently is crucial.

This decision will explain the rationale behind the selected database engine, with consideration to the above project requirements.

## Options

### Option A: BBolt

BBolt is a Go-native storage database designed for application data. It's simple, fast, and reliable - it only supports simple key/value storage but can scale to massive amounts of data, and has proper thread safety checks
with normal considerations given to a concurrent program. Another plus is that BBolt doesn't require a standalone server, and thus makes distribution of the application easier.

The Supertext program doesn't necessarily require a relational database, as application data is not necessarily relating to other data. Shows are stored individually with limited to no
connection to each other, and user information isn't a desired feature. Realistically, the majority of the data stored would be supertext uploads.

With this consideration, a read-intensive key/value storage solution like BBolt fits the project needs well.

However, BBolt isn't perfect, and does have some notable downsides:

- **Thread Safety**: BBolt does not guarantee thread safety for individual transactions. Creating new transactions is guaranteed to be thread safe, but the results of those transactions may not be safe. As such,
  fetching data from the database may incur a high performance cost by running numerous transactions, or introduce complexity by external locking of the results.
- **Memory Management**: Even though BBolt is written in Go, it does have some notable considerations with memory, such as its usage of memory-mapped files or how BBolt unmaps information once a transaction is completed (which may affect loading data).

However, even with the considerations presented, BBolt may still be a solid option for the Supertext program.

### Option B: SQLite

SQLite is an embedded relational SQL database written in C. Notably, it does not require any standalone server either, yet supports the vast majority of SQL features. It's used in numerous high-performance and security-critical contexts
such as web browsers, and is well suited to storing application-specific data. It also supports concurrency with caveats about compilation and statement usage.

SQLite does have some significant downsides to its usage:

- **Complexity**: SQLite is a full-featured SQL "server," and as such introduces all the normal complexities of writing, testing, and executing SQL queries.
- **Security**: SQLite, similar to any other SQL server, is still vulnerable to SQL injections and other security concerns. There are tools such as prepared statements that can alleviate some of these concerns,
  but this adds more expectations onto developers.

SQLite is a solid option, but likely overkill for this project.

### Option C: Do Nothing

It is possible to store data onto the disk manually without a database engine.

This would be immediately more complicated, and rely on us developing proper concurrent file-management solutions within the application. As well, this solution does not scale whatsoever,
and would likely create a massive bottleneck.

## Analysis

### Maintainability

|          | Positive | Negative |
| :------: | :------: | :------: |
| *BBolt*  | Extremely simple API with no SQL knowledge necessary. Minimal feature-set and API means less likeliness of introducing serious flaws. | Not as trusted or tested as SQLite, not as well documented as SQL. | 
| *SQLite* | Industry-standard, well-documented, the most used application database in the world. Well-updated. | Complicated API with numerous considerations that must be taken for usage such as remembering to always used prepared SQL statements. |

### Complexity

|          | Positive | Negative |
| :------: | :------: | :------: |
| *BBolt*  | Provides a simple API, with a focus on key/value storage. The entire read/write functionality is simply just three operations [(Put, Get, Delete)](https://github.com/etcd-io/bbolt#using-keyvalue-pairs). | No SQL query support, and complicated transactions are impossible by the very nature of BBolt. [(Source)](https://github.com/etcd-io/bbolt#bbolt) | 
| *SQLite* | SQL is extremely standard, well-documented, and supported by nearly all ORMs to reduce complexity. Go itself has numerous libraries ([gosqlite](https://gosqlite.com/), [go-sqlite3](https://github.com/mattn/go-sqlite3)) that provide relatively simple access to SQLite. | SQL itself is complicated and without the use of an ORM would require developers to know and understand SQL. As well, SQL can introduce security vulnerabilites with improperly-written SQL statements, and requires numerous prevention methods to keep secure. [Source](https://cheatsheetseries.owasp.org/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.html)

## Recommendation and Verification Plan

I propose the Supertext project use the BBolt key/value storage system. The application does not benefit from the added features of SQLite due to the nature of the application's use cases; if anything, SQLite
would introduce more unnecessary complexity and risk introducing more bugs or security vulnerabilities.

Acceptance criteria:

- The application can serve a maximum number of clients in the hundreds with minimal performance loss.
- Developers are able to interact with the database easily.
- Uploads or edits of supertext data are efficient, and do not cause corruption or other issues.

These can be checked via manual testing of the application and the development procedures.

## Constraints

- **Licensing**: BBolt is licensed under the MIT License, which is compatible with this project.
- **Data use**: Data stored within the database is local-only, but may contain confidential IP of an opera production. The database data must be treated as confidential and be protected as such.
- **Confidentiality**: Not applicable to this project


## Delegation and Validation

- **What you handed to AI tools**: No part of the writing, research, or general development of this RFC were tainted by the willful usage of Generative AI tools.
- **What checks made that reasonable**: This repo's CI runs `go fmt` and `go test` workflows. This document is entirely human-written and therefore understood by someone with full knowledge of the repository and project.
- **Where you stopped and validated the result yourself, and why those points.**: The entire result was validated by a human, as no artificial intelligence tool was used in the first place.

## Decision

Proposed (N/A)
