# ADR-001: Endpoint Conventions 

### Submitters
- Tom Sharp

## Status: Proposed

## Context

This ADR lays out the conventions for creating and maintaining endpoints for the Supertext project. All endpoints will be expected to follow these conventions to ensure readability.

## Decision

### Naming Conventions
Endpoints should be lower case and use kebab-case. File extensions should be omitted from endpoints. In general, try to follow general [REST API naming conventions and best practices](https://restfulapi.net/resource-naming/) 

### Versioned Endpoints
All endpoints should be versioned. (EX /v1/upload, /v2/upload). Versions should be updated when major client-facing changes are introduced, including changing the payload parameters, arguments, etc. 

Endpoints should be backwards compatible even if they aren't used, meaning the end user should still be able to access the software if they are access it through an old endpoint.

### Payloads
Payloads should be formatted as JSON or [multipart/form-data](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/MIME_types#multipartform-data) for all applicable HTTP methods. Keys should be formatted using camel case. 

Payload information should be nested in keywords for what is being returned. The current allowed keywords are:
- `data`: main information for the endpoint.
- `error`: contains error information for the endpoint

Successful responses must use the data key, while error responses must use the error key. These keys must not be combined in a single response unless explicitly required by the API specification.

Example success:
```json
{
  "data": {
    "id": "1234",
    "name": "La Toriador",
    "pages": 120,
    "content": ...
  }
}
```

Example failure:
```json
{
  "error": {
    "code": "NOT FOUND",
    "message": "The requested resource was not found."
  }
}
```

### Status / Error Codes
Status codes should follow [standard HTTP status codes](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status)

Error codes should be formatted as uppercase, and should match HTTP error code. Message field should be used to describe what exactly went wrong.

### Arguments
Path separators should be used for accessing an individual resources (e.g. /v1/edit/<presentation_name>) and to represent structural heirarchy.

Query separators should be used for searches, filters, etc. (e.g. /v1/select?sort=author)

## Consequences

### Positive
- **Consistency:** Standardized naming conventions, payload structures, and error codes make endpoints easier to read, understand, and maintain across the project.
- **Predictability:** Developers can anticipate how endpoints behave and how responses are structured, reducing the learning curve for new contributors and consumers of the API.
- **Maintainability:** Versioned endpoints allow major changes to be introduced without immediately breaking existing clients.

### Negative
- **Additional overhead:** Maintaining multiple API versions may require extra development, testing, and documentation as the API evolves.
- **Implementation complexity:** Enforcing consistent response formats and error handling across all endpoints may require additional validation, shared utilities, and testing.

### Neutral
- **Documentation requirements:** Endpoint changes, version updates, and exceptions to these conventions will need to be documented to keep the API specification accurate.
- **Convention enforcement:** Code reviews, automated tests, or linting tools may be needed to ensure new and modified endpoints follow these standards.

## Change Log
- V1 (10/16/2026) Initial ADR [Tom]
