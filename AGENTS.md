## Project Overview

This is a distributed Supertext application using WebSockets to synchronize the
individual patron screens with the live performance. The application is
intended to stay as local first as possible, using a self-managed WiFi connection
through a Raspberry Pi which also runs the server.

## Set up Commands

- Get dependencies: `go get`
- Run code: `go run .`

## Code Style
- `Gofmt` is the standard format of the codebase

## PR Instructions
- Title follow scoped commits or the feature ticket title
- Always run `go fmt` and `go test` before committing.
