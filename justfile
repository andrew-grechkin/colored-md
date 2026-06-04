#!/usr/bin/env -S just --one --justfile

export tool := 'colored-md'

# Build the binary to cache directory
@build: fix
    go build -o "$XDG_CACHE_HOME/go/bin/"

# Install the binary globally
@install:
    go install "$tool"

# Format Go source code
@fix:
    go fmt
    go fix

# Run Go linter
@lint: cc
    go vet

# Report functions over cyclomatic complexity 15 (installs gocyclo if missing)
@cc:
    test -x "$XDG_CACHE_HOME/go/bin/gocyclo" || GOBIN="$XDG_CACHE_HOME/go/bin" go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
    "$XDG_CACHE_HOME/go/bin/gocyclo" -over 15 .

# Update Go dependencies
@update:
    go get -u
    go mod tidy
