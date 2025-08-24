.PHONY: test bench lint gen ref test-oracle test-all test-c2go golangci-lint install-lint install-smrcptr fmt fix-fmt

test:
	go test -v ./...

bench:
	go test -bench=. -benchmem ./...

# Note: golangci-lint tags use v1.x; latest stable at time of writing.
# v2.x does not exist as a module tag; pin to latest v1 instead.
GOLANGCI_LINT_VERSION ?= v2.4.0
GOBIN ?= $(shell go env GOPATH)/bin
GOLANGCI_LINT := $(GOBIN)/golangci-lint
SMRCPTR := $(GOBIN)/smrcptr

# Installs golangci-lint via official script into GOPATH/bin
install-lint:
	@echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION) to $(GOLANGCI_LINT) via official script"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | \
		sh -s -- -b $(GOBIN) $(GOLANGCI_LINT_VERSION)

# Installs smrcptr for consistent receiver type checking
install-smrcptr:
	@echo "Installing smrcptr to $(SMRCPTR)"
	@go install github.com/nikolaydubina/smrcptr@latest

# Runs fmt, vet, golangci-lint, and smrcptr (installs them if missing)
lint: fmt $(GOLANGCI_LINT) $(SMRCPTR)
	go vet ./...
	$(GOLANGCI_LINT) run
	@echo "Running smrcptr for consistent receiver type checking..."
	@$(SMRCPTR) ./...

$(GOLANGCI_LINT):
	@$(MAKE) install-lint

$(SMRCPTR):
	@$(MAKE) install-smrcptr

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

# Run c2go parity tests (require cgo). Uses local GOCACHE to avoid sandboxed home cache writes.
# Usage: make test-c2go [TEST=TestName] [VERBOSE=1] [TIMEOUT=duration]
# Examples:
#   make test-c2go                              # Run all tests (10s timeout)
#   make test-c2go TEST=Test_getIcosahedronFaces  # Run specific test
#   make test-c2go VERBOSE=1                    # Run all tests in verbose mode
#   make test-c2go TEST=Test_getIcosahedronFaces VERBOSE=1  # Run specific test verbosely
#   make test-c2go TIMEOUT=30s                  # Run all tests with 30s timeout
H3VER ?= 4.3.0
TEST ?=
VERBOSE ?=
TIMEOUT ?= 10s
test-c2go:
	@if [ -n "$(TEST)" ]; then \
		echo "Running c2go parity test: $(TEST) (requires cgo)..."; \
	else \
		echo "Running c2go parity tests (requires cgo)..."; \
	fi
	@# Prefer clang from Xcode CLT if available
	@CC=$$(command -v xcrun >/dev/null 2>&1 && xcrun --find clang || true); \
	CXX=$$(command -v xcrun >/dev/null 2>&1 && xcrun --find clang++ || true); \
	SDKROOT=$$(command -v xcrun >/dev/null 2>&1 && xcrun --sdk macosx --show-sdk-path || true); \
	INC_BASE="$(PWD)/testref/h3-$(H3VER)/src/h3lib"; \
	if [ ! -d "$$INC_BASE" ]; then \
		echo "H3 C sources not found at $$INC_BASE. Run 'make ref' or set H3VER=..."; \
		exit 1; \
	fi; \
	VERBOSE_FLAG=""; \
	if [ -n "$(VERBOSE)" ]; then VERBOSE_FLAG="-v"; fi; \
	TEST_FLAG=""; \
	if [ -n "$(TEST)" ]; then TEST_FLAG="-run=$(TEST)"; fi; \
	TIMEOUT_FLAG="-timeout=$(TIMEOUT)"; \
	GOCACHE=$(PWD)/.gocache \
	CGO_ENABLED=1 CC="$$CC" CXX="$$CXX" SDKROOT="$$SDKROOT" \
	CGO_CPPFLAGS="-I$$INC_BASE/include -I$$INC_BASE/lib" \
	CGO_CFLAGS="-ffunction-sections -fdata-sections" \
	CGO_LDFLAGS="-Wl,-dead_strip" \
	go test $$VERBOSE_FLAG $$TEST_FLAG $$TIMEOUT_FLAG -tags="c2go" ./internal/c2go || { \
		echo; \
		echo "c2go tests failed. If the error mentions 'use of cgo not supported':"; \
		echo " - Ensure Go was installed with cgo support (official pkg/Homebrew)."; \
		echo " - Ensure a C toolchain is present (macOS: xcode-select --install)."; \
		exit 1; \
	}
fmt:
	@echo "Checking gofmt formatting..."
	@files=$$(gofmt -s -l .); \
	if [ -n "$$files" ]; then \
		echo "gofmt found unformatted files:"; echo "$$files"; \
		exit 1; \
	else \
		echo "gofmt OK"; \
	fi

fix-fmt:
	@echo "Formatting Go files..."
	@gofmt -s -w .
	@echo "All Go files formatted"
