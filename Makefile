.PHONY: test bench lint gen ref test-oracle test-all golangci-lint install-lint fmt

test:
	go test -v ./...

bench:
	go test -bench=. -benchmem ./...

# Note: golangci-lint tags use v1.x; latest stable at time of writing.
# v2.x does not exist as a module tag; pin to latest v1 instead.
GOLANGCI_LINT_VERSION ?= v2.4.0
GOBIN ?= $(shell go env GOPATH)/bin
GOLANGCI_LINT := $(GOBIN)/golangci-lint

# Installs golangci-lint via official script into GOPATH/bin
install-lint:
	@echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION) to $(GOLANGCI_LINT) via official script"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | \
		sh -s -- -b $(GOBIN) $(GOLANGCI_LINT_VERSION)

# Runs vet and golangci-lint (installs it if missing)
lint: $(GOLANGCI_LINT)
	go vet ./...
	$(GOLANGCI_LINT) run

$(GOLANGCI_LINT):
	@$(MAKE) install-lint

gen:
	@echo "Reserved for future table generation"

ref:
	@echo "Building H3 reference oracle for test validation..."
	$(MAKE) -C testref
	@echo "Oracle built successfully at testref/h3ref"
	@echo "Run 'make -C testref test' to verify installation"

# Run tests that require the external C oracle as well
test-oracle: ref
	@echo "Running Go tests with oracle tag..."
	go test -v -tags=oracle ./...

# Convenience: run both regular and oracle-tagged tests
test-all: test test-oracle
fmt:
	@echo "Checking gofmt formatting..."
	@files=$$(gofmt -s -l .); \
	if [ -n "$$files" ]; then \
		echo "gofmt found unformatted files:"; echo "$$files"; \
		exit 1; \
	else \
		echo "gofmt OK"; \
	fi
