.PHONY: build test lint fmt vet vulncheck check clean install-tools test-cover generate schema

MODULE  := github.com/misham/linear-cli
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
CLIENT_ID ?= $(LNR_CLIENT_ID)
LDFLAGS := -X $(MODULE)/internal/version.Version=$(VERSION) \
           -X $(MODULE)/internal/version.Commit=$(COMMIT) \
           -X $(MODULE)/internal/version.Date=$(DATE) \
           -X $(MODULE)/internal/auth.defaultClientID=$(CLIENT_ID)

# Build
build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o lnr ./cmd/lnr

# Test with race detector
test:
	go test -race -count=1 ./...

# Test with coverage
test-cover:
	go test -race -count=1 -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Lint (golangci-lint runs staticcheck, errcheck, govet, gosec, and more)
lint:
	golangci-lint run ./...

# Format (gofumpt is a strict superset of gofmt)
fmt:
	gofumpt -w -modpath $(MODULE) .

# Format check (CI — fails if files need formatting)
fmt-check:
	@test -z "$$(gofumpt -l -modpath $(MODULE) .)" || (echo "files need formatting:"; gofumpt -l -modpath $(MODULE) .; exit 1)

# Go vet
vet:
	go vet ./...

# Vulnerability check
vulncheck:
	govulncheck ./...

# Run all static checks
check: fmt-check vet lint vulncheck

# Download Linear GraphQL schema
schema:
	curl -fsSL -o schema.graphql \
		https://raw.githubusercontent.com/linear/linear/master/packages/sdk/src/schema.graphql

# Generate GraphQL client code
generate: schema
	go generate ./...
	@# Add omitempty to all JSON struct tags so nil fields are omitted instead of
	@# being sent as explicit nulls (which the Linear API rejects in filter inputs).
	perl -pi -e 's/`json:"([^"]+)"`/`json:"$$1,omitempty"`/g' internal/api/generated/operations.gen.go

# Install development tools (versions pinned via go.mod tool directives)
install-tools:
	go install tool

# Clean build artifacts
clean:
	rm -f lnr coverage.out
