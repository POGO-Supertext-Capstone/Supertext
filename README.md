# Supertext for the Portland Opera to Go

`Supertext` is a Go web server for providing Supertext to patrons during Portland Opera to Go shows, using devices
that attendees already own.

## Quickstart

If you haven't already, install [Go 1.27 or later](https://go.dev/dl/).

With Go installed, run `go run ./...` and visit [localhost:8081](http://localhost:8081) in a web browser.

## Building

Go will automatically download all dependencies.

```bash
$ mkdir build && go build -o build ./...
```

The Supertext program will be compiled and saved to the `build` directory.

## Testing

To run the unit tests:

```bash
$ go test -v ./...
```

## Contributing

Pull requests are welcome. Please read [CONTRIBUTING.md](./CONTRIBUTING.md) before making any changes.

Ensure that `go fmt` is run on all changes, and that all tests pass. Please make sure to adjust or add new tests
as necessary.

## Licensing

The Supertext program is licensed under the [GPL-3.0](./COPYING) license.
