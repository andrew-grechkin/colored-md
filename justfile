#!/usr/bin/env -S just --one --justfile

set export

export tool := 'colored-md'

export GOBIN := `echo "${GOBIN:-${GOPATH:-$HOME/go}/bin}"`

alias fmt := fix

# Default recipe
[private]
@default:
    just --list

# Build the binary
@build: fix
    go build

# Run complexity lint
cc:
    #!/usr/bin/env -S bash -Eeuo pipefail
    [[ -x "$GOBIN/gocyclo" ]] || go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
    "$GOBIN/gocyclo" -over 15 .

# Remove build artifacts (everything in .gitignore)
@clean:
    git clean -Xdf

# Format and modernize Go source code
@fix:
    go fmt
    go fix

# Install the binary to $GOBIN
@install:
    go install

# Run Go linter
@lint: cc
    go vet

# Update Go dependencies
@update:
    go get -u
    go mod tidy

# Upgrade Golang
upgrade: && update
    #!/usr/bin/env -S bash -Eeuo pipefail
    go get go@latest

# Scan dependencies for known CVEs
vulncheck:
    #!/usr/bin/env -S bash -Eeuo pipefail
    [[ -x "$GOBIN/govulncheck" ]] || go install golang.org/x/vuln/cmd/govulncheck@latest
    "$GOBIN/govulncheck" ./...
