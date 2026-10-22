# ADR 002: Database Selection

### Proposers
- Jamie Murphy <murpjami@oregonstate.edu>

## Status: Proposed

## Context

This ADR standardises the database format and engine selection for the Supertext program. It explains
the rationale behind the selected database solution and the expected format.

Supertext will use a database to store uploaded, processed supertext from an admin. It will
also need to provide these processed supertexts to clients efficiently.

## Decision

Supertext will use the BBolt storage database for storing application data.

A primary [bucket](https://github.com/etcd-io/bbolt#using-buckets) will be provided named `supertext`. This bucket will contain a key, `shows`, with an array of show names.
This bucket will also provide account details for the admin user.

Opera-provided show supertext will be stored within a separate Bucket for each show. These buckets will be named the same as the
show name provided within the primary bucket. The actual processed supertext data will be stored under the key `text`.

Other show information, if necessary, will be stored under the show bucket. Each bucket will be stored within a subdirectory under
`$XDG_DATA_HOME` or equivalent, with the filename `supertext.db`.

### BBolt

BBolt, or Bolt, is a Go native key-value storage database. Notably, it is easier to use than a full SQL Relational database,
and is preferred for read-intensive workflows.

The Supertext program doesn't necessarily require much of a relational database. A key-value data store matches most, if not all,
of the desired functionality; clients should realistically not be saving much data to a database besides possible session data and the only main writes
are initial supertext uploads.

Therefore, our needs are primarily read-intensive with minimal relational requirements. A database like BBolt, built into
the application, will be the most efficient option.

### Alternatives

There are many relational SQL databaes such as PostgreSQL, SQLite, MySQL and others. These provide overhead in
the requirement of SQL parsing and possible SQL injection. As well, most databases are standalone servers
that run separately from the Supertext program. Realistically, this is not something for the Supertext server to
care about - we wouldn't have multiple application servers or complex multi-app networking.

The main exception to these rules is SQLite

#### SQLite

SQLite is an embedded relational database written in C. It is used in many high-performance and security-critical applications
such as web browsers, and is well suited to storing application data. It does support multi-threading with caveats about compliation and usages however.
Prepared statements can also be used to reduce any security concerns via SQL Injection. However, SQLite still has
the overhead of parsing and writing SQL. SQLite is overkill for this project and introduces unnecessary complexity.

## Consequences

### Positive

- **Maintainability:** Bolt is easier to work with and has a less complicated API.
- **Simplicity:** Bolt is much simpler and fits our use cases naturally.

### Negative

- **Threading complexity:** Individual transactions are not thread-safe, and so transactions must be initiated within
  each goroutine separately. However, creating transactions is thread-safe.

## Changelog

- Rev 1 (2026-10-22) - Jamie Murphy <murpjami@oregonstate.edu>
  Initial ADR
