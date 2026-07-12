.PHONY: test test-cli test-cli-process test-cli-diff build-cli build-cli-cross check-cli-inventory test-upstream-fixtures bench lint test-c2go golangci-lint install-lint install-smrcptr fmt fix-fmt coverage coverage-html coverage-c2go coverage-c2go-html coverage-all check-unsafe api-inventory check-api test-uberdiff test-uberbench bench-uber upstream-diff check-test-inventory check-docs gen-benchdocs check-benchdocs

# Enforces DR-007 (docs/public-api-architecture.md): the production library is
# safe Go only. Two independent layers:
#  1) Build-selection check: for every normal build mode, ask the toolchain
#     which packages/files it would compile (including tests) and fail if any
#     package in THIS module imports unsafe.
#  2) Tag allowlist (platform/GOOS-independent): any file importing unsafe
#     must carry a build constraint containing c2go.
check-unsafe:
	@echo "check-unsafe: layer 1 (build-selection, all normal build modes)..."
	@for cgo in 0 1; do \
		for tags in "" "race"; do \
			out=$$(CGO_ENABLED=$$cgo go list -tags "$$tags" \
				-f '{{.ImportPath}}: {{join .Imports " "}} {{join .TestImports " "}} {{join .XTestImports " "}}' \
				./... | grep -w unsafe || true); \
			if [ -n "$$out" ]; then \
				echo "FAIL: unsafe reachable with CGO_ENABLED=$$cgo tags='$$tags':"; \
				echo "$$out"; \
				exit 1; \
			fi; \
		done; \
	done
	@echo "check-unsafe: layer 2 (tag allowlist)..."
	@bad=""; \
	for f in $$(grep -rl '^[[:space:]]*"unsafe"$$' --include='*.go' . 2>/dev/null | grep -v '^\./testref' | grep -v '^\./\.gocache'); do \
		if ! head -1 "$$f" | grep -q 'go:build.*c2go'; then \
			bad="$$bad $$f"; \
		fi; \
	done; \
	if [ -n "$$bad" ]; then \
		echo "FAIL: unsafe imported outside cgo && c2go tagged files:$$bad"; \
		exit 1; \
	fi
	@echo "check-unsafe: OK"

# Documentation gate: every relative Markdown link and #anchor must resolve
# (tools/docscheck). Runs in CI on every push/PR, including docs-only changes.
check-docs:
	@go run ./tools/docscheck
	@go run ./tools/benchdocs -verify

# Regenerate/verify selected README benchmark excerpts from the committed
# benchstat CSV and metadata artifacts (offline).
gen-benchdocs:
	@go run ./tools/benchdocs -write

check-benchdocs:
	@go run ./tools/benchdocs -verify

# Regenerate the generated tables in docs/comparison-uber-h3-go.md from
# docs/comparison-uber-h3-go.csv (offline; uses committed inventories).
gen-ubercompare:
	@go run ./tools/ubercompare -write

# Drift gate: comparison matrix vs C-API inventory vs locked API surface vs
# the generated doc tables (offline; runs in CI).
check-ubercompare:
	@go run ./tools/ubercompare -verify

# Regenerate the C-API inventory (requires testref sources; see make -C testref h3-source).
api-inventory:
	@go run ./tools/apiinventory -h3ver $(H3VER) > docs/c-api-inventory.csv
	@echo "docs/c-api-inventory.csv regenerated (H3 $(H3VER))"

# Completeness gate: every H3 C public function must be ported AND publicly
# represented (an "H3 C API:" doc line or a documented omission).
# Requires testref sources (make -C testref h3-source downloads them).
check-api:
	@go run ./tools/apiinventory -h3ver $(H3VER) -verify

# Ecosystem-level completeness gate: named test cases, input-driven programs,
# CLI registrations, fuzzers, benchmarks, filters, helpers, fixtures, and
# build definitions must have reviewed dispositions. Requires testref sources.
check-test-inventory:
	@go run ./tools/testinventory -h3ver $(H3VER) -verify

build-cli:
	@CGO_ENABLED=0 go build ./cmd/h3

test-cli:
	@CGO_ENABLED=0 go test ./internal/cli

test-cli-process:
	@CGO_ENABLED=0 go test ./internal/cli -run '^TestBinaryProcessContract$$'

check-cli-inventory:
	@go run ./tools/cliinventory -upstream testref/h3-$(H3VER) -verify

# Builds the pristine upstream C CLI under /tmp and differentially executes all
# 170 registered scenarios. Requires cmake and a C toolchain.
test-cli-diff:
	@cmake -E remove_directory /tmp/h3-cli-src-$(H3VER)
	@cmake -E remove_directory /tmp/h3-cli-$(H3VER)
	@cmake -E copy_directory testref/h3-$(H3VER) /tmp/h3-cli-src-$(H3VER)
	@cmake -E rm -f /tmp/h3-cli-src-$(H3VER)/src/h3lib/include/h3api.h
	@cmake -S /tmp/h3-cli-src-$(H3VER) -B /tmp/h3-cli-$(H3VER) -DBUILD_FILTERS=ON -DBUILD_TESTING=OFF -DENABLE_FORMAT=OFF
	@cmake --build /tmp/h3-cli-$(H3VER) --target h3_bin
	@H3_CLI_C_BINARY=/tmp/h3-cli-$(H3VER)/bin/h3 go test ./internal/cli -run '^TestDifferentialWithCCLI$$' -count=1

build-cli-cross:
	@for target in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64; do \
		os=$${target%/*}; arch=$${target#*/}; suffix=""; \
		if [ "$$os" = windows ]; then suffix=.exe; fi; \
		echo "building h3 $$target"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -o /tmp/h3-$$os-$$arch$$suffix ./cmd/h3 || exit 1; \
	done

# Pure-Go ports of the three large input-driven upstream executables.
test-upstream-fixtures:
	@H3_UPSTREAM_FIXTURE_ROOT=testref/h3-$(H3VER)/tests/inputfiles \
		CGO_ENABLED=0 go test -run '^TestUpstream.*Fixtures$$' .

# Symbol-level diff between two upstream H3 trees, mapped to the Go port.
# Usage: make upstream-diff FROM=4.3.0 TO=4.4.0
# (both trees must exist under testref/; fetch with
#  make -C testref H3_VERSION=<ver> h3-source)
FROM ?= 4.3.0
TO ?= $(H3VER)
upstream-diff:
	@go run ./tools/upstreamdiff -from testref/h3-$(FROM) -to testref/h3-$(TO)

# Usage: make test [TEST=TestName] [VERBOSE=1] [TIMEOUT=duration] [COVERAGE=1]
# Examples:
#   make test                              # Run all tests (default timeout)
#   make test TEST=TestPolygonToCells_ZeroSize  # Run specific test
#   make test VERBOSE=1                    # Run all tests in verbose mode
#   make test TEST=TestPolygonToCells_ZeroSize VERBOSE=1  # Run specific test verbosely
#   make test TIMEOUT=30s                  # Run all tests with 30s timeout
#   make test COVERAGE=1                   # Run tests with coverage report
#   make test COVERAGE=1 COVERPROFILE=coverage.out  # Save coverage to specific file
test:
	@if [ -n "$(TEST)" ]; then \
		echo "Running test: $(TEST)..."; \
	else \
		echo "Running all tests..."; \
	fi
	@VERBOSE_FLAG=""; \
	if [ -n "$(VERBOSE)" ]; then VERBOSE_FLAG="-v"; fi; \
	TEST_FLAG=""; \
	if [ -n "$(TEST)" ]; then TEST_FLAG="-run=$(TEST)"; fi; \
	TIMEOUT_FLAG=""; \
	if [ -n "$(TIMEOUT)" ]; then TIMEOUT_FLAG="-timeout=$(TIMEOUT)"; fi; \
	COVERAGE_FLAG=""; \
	if [ -n "$(COVERAGE)" ]; then \
		COVERPROFILE="$${COVERPROFILE:-coverage.out}"; \
		COVERAGE_FLAG="-cover -coverprofile=$$COVERPROFILE"; \
		echo "Coverage will be saved to: $$COVERPROFILE"; \
	fi; \
	CGO_ENABLED=0 go test $$VERBOSE_FLAG $$TEST_FLAG $$TIMEOUT_FLAG $$COVERAGE_FLAG ./... && \
	if [ -n "$(COVERAGE)" ]; then \
		echo ""; \
		echo "Coverage report generated. View with:"; \
		echo "  go tool cover -func=$$COVERPROFILE    # Function coverage"; \
		echo "  go tool cover -html=$$COVERPROFILE    # HTML report"; \
	fi

bench:
	CGO_ENABLED=0 go test -bench=. -benchmem ./...

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

# Run c2go parity tests (require cgo). Uses local GOCACHE to avoid sandboxed home cache writes.
# Usage: make test-c2go [TEST=TestName] [VERBOSE=1] [TIMEOUT=duration] [COVERAGE=1]
# Examples:
#   make test-c2go                              # Run all tests (10s timeout)
#   make test-c2go TEST=Test_getIcosahedronFaces  # Run specific test
#   make test-c2go VERBOSE=1                    # Run all tests in verbose mode
#   make test-c2go TEST=Test_getIcosahedronFaces VERBOSE=1  # Run specific test verbosely
#   make test-c2go TIMEOUT=30s                  # Run all tests with 30s timeout
#   make test-c2go COVERAGE=1                   # Run tests with coverage report
#   make test-c2go COVERAGE=1 COVERPROFILE=coverage-c2go.out  # Save to specific file
H3VER ?= 4.4.0
TEST ?=
VERBOSE ?=
TIMEOUT ?= 30s
COVERAGE ?=
COVERPROFILE ?=
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
	APPS_BASE="$(PWD)/testref/h3-$(H3VER)/src/apps/applib"; \
	if [ ! -d "$$INC_BASE" ]; then \
		echo "H3 C sources not found at $$INC_BASE. Run 'make -C testref h3-source' or set H3VER=..."; \
		exit 1; \
	fi; \
	VERBOSE_FLAG=""; \
	if [ -n "$(VERBOSE)" ]; then VERBOSE_FLAG="-v"; fi; \
	TEST_FLAG=""; \
	if [ -n "$(TEST)" ]; then TEST_FLAG="-run=$(TEST)"; fi; \
	TIMEOUT_FLAG="-timeout=$(TIMEOUT)"; \
	COVERAGE_FLAG=""; \
	if [ -n "$(COVERAGE)" ]; then \
		COVERPROFILE="$${COVERPROFILE:-coverage-c2go.out}"; \
		COVERAGE_FLAG="-cover -coverprofile=$$COVERPROFILE"; \
		echo "Coverage will be saved to: $$COVERPROFILE"; \
	fi; \
	TOOLCHAIN_ENV=""; \
	if [ -n "$$CC" ]; then TOOLCHAIN_ENV="CC=$$CC CXX=$$CXX SDKROOT=$$SDKROOT"; fi; \
	CFLAGS_ENV=""; LDFLAGS_ENV="-lm"; \
	if [ "$$(uname -s)" = "Darwin" ]; then \
		CFLAGS_ENV="-ffunction-sections -fdata-sections"; \
		LDFLAGS_ENV="-Wl,-dead_strip"; \
	fi; \
	env $$TOOLCHAIN_ENV \
	GOCACHE=$(PWD)/.gocache \
	CGO_ENABLED=1 \
	CGO_CPPFLAGS="-I$$INC_BASE/include -I$$INC_BASE/lib -I$$APPS_BASE/include -I$$APPS_BASE/lib" \
	CGO_CFLAGS="$$CFLAGS_ENV" \
	CGO_LDFLAGS="$$LDFLAGS_ENV" \
	go test $$VERBOSE_FLAG $$TEST_FLAG $$TIMEOUT_FLAG $$COVERAGE_FLAG -tags="c2go" ./... || { \
		echo; \
		echo "c2go tests failed. If the error mentions 'use of cgo not supported':"; \
		echo " - Ensure Go was installed with cgo support (official pkg/Homebrew)."; \
		echo " - Ensure a C toolchain is present (macOS: xcode-select --install)."; \
		exit 1; \
	}; \
	if [ -n "$(COVERAGE)" ] && [ -z "$$TEST_FLAG" -o $$? -eq 0 ]; then \
		echo ""; \
		echo "Coverage report generated. View with:"; \
		echo "  go tool cover -func=$$COVERPROFILE    # Function coverage"; \
		echo "  go tool cover -html=$$COVERPROFILE    # HTML report"; \
	fi
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

# Convenience target to run tests with coverage and display report
coverage:
	@echo "Running tests with coverage..."
	@$(MAKE) test COVERAGE=1 COVERPROFILE=coverage.out
	@echo ""
	@echo "Function coverage summary:"
	@go tool cover -func=coverage.out | tail -5

# Generate and open HTML coverage report
coverage-html: coverage
	@echo "Generating HTML coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Opening coverage.html in browser..."
	@if command -v open >/dev/null 2>&1; then \
		open coverage.html; \
	elif command -v xdg-open >/dev/null 2>&1; then \
		xdg-open coverage.html; \
	else \
		echo "Please open coverage.html manually"; \
	fi

# Convenience target to run c2go parity tests with coverage and display report
coverage-c2go:
	@echo "Running c2go parity tests with coverage..."
	@$(MAKE) test-c2go COVERAGE=1 COVERPROFILE=coverage-c2go.out TIMEOUT=30s
	@echo ""
	@echo "Function coverage summary:"
	@go tool cover -func=coverage-c2go.out | tail -5

# Generate and open HTML coverage report for c2go tests
coverage-c2go-html: coverage-c2go
	@echo "Generating HTML coverage report for c2go tests..."
	@go tool cover -html=coverage-c2go.out -o coverage-c2go.html
	@echo "Opening coverage-c2go.html in browser..."
	@if command -v open >/dev/null 2>&1; then \
		open coverage-c2go.html; \
	elif command -v xdg-open >/dev/null 2>&1; then \
		xdg-open coverage-c2go.html; \
	else \
		echo "Please open coverage-c2go.html manually"; \
	fi

# Combined coverage: run both test suites and merge coverage
coverage-all:
	@echo "Running all tests with coverage (regular + c2go)..."
	@echo "Step 1: Running regular tests..."
	@$(MAKE) test COVERAGE=1 COVERPROFILE=coverage-regular.out
	@echo ""
	@echo "Step 2: Running c2go parity tests..."
	@$(MAKE) test-c2go COVERAGE=1 COVERPROFILE=coverage-c2go-temp.out TIMEOUT=30s
	@echo ""
	@echo "Step 3: Merging coverage profiles..."
	@echo "mode: set" > coverage-all.out
	@tail -n +2 coverage-regular.out >> coverage-all.out 2>/dev/null || true
	@tail -n +2 coverage-c2go-temp.out >> coverage-all.out 2>/dev/null || true
	@rm -f coverage-regular.out coverage-c2go-temp.out
	@echo ""
	@echo "Combined coverage summary:"
	@go tool cover -func=coverage-all.out | tail -5
	@echo ""
	@echo "View combined coverage with:"
	@echo "  go tool cover -func=coverage-all.out    # Function coverage"
	@echo "  go tool cover -html=coverage-all.out    # HTML report"

# Differential tests against the official cgo binding (separate module;
# requires cgo + network). Not run by default CI.
test-uberdiff:
	cd interop/uberdiff && go test ./...

# Equivalence tests of the comparative benchmark module (separate module;
# requires cgo + network). Gates the benchmarks below.
test-uberbench:
	cd interop/uberbench && go test ./...

# Full comparative benchmark + memory suite vs the official cgo binding.
# Writes artifacts to docs/benchmarks/<goos>-<goarch>/ (see
# interop/uberbench/README.md; tune with COUNT/BENCHTIME/MEMITERS/OUTDIR).
bench-uber:
	cd interop/uberbench && ./run.sh
