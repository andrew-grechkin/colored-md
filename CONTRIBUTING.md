# Contributing to colored-md

## Development Setup

To get started with development, you'll need Go (version 1.21 or higher) installed on your system. Optionally, you can install `just` for simplified command execution.

### Building the Project

If you have `just` installed, you can build the executable with:

```bash
just build
```

This command compiles the `main.go` file and places the `colored-md` executable in `$XDG_CACHE_HOME/go/bin/`.

Alternatively, without `just`, you can use `go build`:

```bash
go build -o "$(go env GOCACHE)/bin/colored-md" .
```

Ensure that `$XDG_CACHE_HOME/go/bin/` (or `$(go env GOCACHE)/bin/`) is included in your system's `PATH` environment variable to run `colored-md` from any directory.

### Updating Dependencies

To update the project's Go dependencies and clean up `go.mod` and `go.sum`, use:

```bash
just update
```

## Code Style and Linting

This project follows standard Go idioms and formatting. Please ensure your code is formatted with `go fmt`:

```bash
go fmt ./...
```

We also recommend running `go vet` to catch common errors:

```bash
go vet ./...
```

## Testing

Currently, there are no automated tests provided in this repository.
