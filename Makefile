.PHONY: help hello test test-all build lint vet vulncheck clean

# Default target — print available commands.
help:
	@echo "go-lab — common tasks"
	@echo ""
	@echo "  make hello       Run lesson 01 (prints 'Hello World!')"
	@echo "  make test        Run every lesson's tests (skips intentional 'before/' panics)"
	@echo "  make test-all    Run every test including the intentional panics (expect failures)"
	@echo "  make build       Compile every lesson"
	@echo "  make vet         Run go vet across the whole module"
	@echo "  make lint        Run golangci-lint"
	@echo "  make vulncheck   Scan dependencies for known vulnerabilities"
	@echo "  make clean       Remove test caches and build artifacts"

hello:
	go run ./lessons/01-hello

# Skip lessons/10-panic-and-recover/*/before — those packages intentionally
# demonstrate crashes. See lessons/10-panic-and-recover/README.md.
test:
	@go list ./... | grep -v '10-panic-and-recover/.*/before$$' | xargs go test

test-all:
	go test ./...

build:
	go build ./...

vet:
	go vet ./...

lint:
	golangci-lint run

vulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

clean:
	go clean -testcache
	rm -f coverage.out coverage.html
