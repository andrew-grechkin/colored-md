# Contributing

## Development Setup

To get started with development, you'll need Go installed on your system.
Optionally, you can install [`just`](https://github.com/casey/just) for simplified command execution.

### Building

Build binary:

```bash
just build
```

Global install (into `$GOBIN`):

```bash
just install
```

Ensure that directory is on your `PATH` to run the utilities from anywhere.

### Updating Dependencies

To update the project's Go dependencies and clean up `go.mod` and `go.sum`, use:

```bash
just update
```

## Code Style and Linting

```bash
just fix
```

```bash
just lint
```

## Testing

Currently, there are no automated tests provided in this repository.
