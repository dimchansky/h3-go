#!/usr/bin/env bash
# Runs the full comparative benchmark suite (this library vs the official
# uber/h3-go cgo binding) and the process-level memory matrix, writing
# reproducible artifacts under docs/benchmarks/<goos>-<goarch>/:
#
#   metadata.txt   environment pin: commit, module versions, toolchains,
#                  hardware, flags, date
#   bench-raw.txt  raw `go test -bench` output (COUNT repetitions)
#   benchstat.txt  benchstat summary, one column per /impl
#   benchstat.csv  the same table, machine-readable
#   memory.tsv     process-level memory matrix from cmd/memprobe
#
# The equivalence tests run first and abort the run on any disagreement:
# numbers are only produced for operation pairings that were just shown to
# return semantically equivalent results.
#
# Tunables (environment variables):
#   COUNT      benchmark repetitions for benchstat      (default 10)
#   BENCHTIME  go test -benchtime per measurement       (default 1s)
#   MEMITERS   iterations per memprobe workload         (default 3)
#   OUTDIR     artifact directory                       (default docs/benchmarks/<goos>-<goarch>)
#   ENVLABEL   free-form environment note for metadata  (default: none)
set -euo pipefail

cd "$(dirname "$0")"

COUNT="${COUNT:-10}"
BENCHTIME="${BENCHTIME:-1s}"
MEMITERS="${MEMITERS:-3}"
GOOS="$(go env GOOS)"
GOARCH="$(go env GOARCH)"
OUTDIR="${OUTDIR:-../../docs/benchmarks/${GOOS}-${GOARCH}}"
# Pinned so summaries are reproducible; bump deliberately.
BENCHSTAT="golang.org/x/perf/cmd/benchstat@v0.0.0-20260709024250-82a0b07e230d"

# Capture provenance before opening any tracked artifact for writing. Shell
# redirection truncates its target before command substitutions execute, so
# checking inside the metadata block would report the benchmark's own output
# as a dirty worktree whenever OUTDIR is already committed.
REPO_COMMIT="$(git -C ../.. rev-parse HEAD)"
REPO_DIRTY=""
if ! git -C ../.. diff --quiet HEAD -- 2>/dev/null; then
    REPO_DIRTY=" (dirty)"
fi

mkdir -p "$OUTDIR"

cpu_model() {
    case "$GOOS" in
        darwin) sysctl -n machdep.cpu.brand_string ;;
        linux) sed -n 's/^model name[[:space:]]*: //p' /proc/cpuinfo | head -1 ;;
        *) echo unknown ;;
    esac
}

CC_BIN="$(go env CC)"

# The vendored-version probe below reads the binding out of the module
# cache; make sure it is there (fresh CI runners start with an empty cache).
go mod download github.com/uber/h3-go/v4 >/dev/null 2>&1 || true

{
    echo "date_utc: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
    echo "repo_commit: ${REPO_COMMIT}${REPO_DIRTY}"
    echo "go_version: $(go version)"
    echo "uber_h3_go: $(go list -m github.com/uber/h3-go/v4)"
    echo "dimchansky_h3_go: $(go list -m github.com/dimchansky/h3-go) (replaced by ../.., i.e. repo_commit above)"
    echo "vendored_h3_c_in_binding: $(cat "$(go list -m -f '{{.Dir}}' github.com/uber/h3-go/v4)/H3_VERSION" 2>/dev/null || echo unknown)"
    echo "pure_go_h3_target: H3 C v4.4.0 (VersionMajor/Minor/Patch of the root module)"
    echo "os: $(uname -a)"
    echo "cpu: $(cpu_model)"
    echo "ncpu_online: $(getconf _NPROCESSORS_ONLN)"
    echo "gomaxprocs: default (= ncpu_online; benchmarks are single-goroutine)"
    echo "cgo_enabled: 1 (required by the binding; the pure-Go library itself needs none)"
    echo "cc: $CC_BIN — $("$CC_BIN" --version 2>/dev/null | head -1)"
    echo "cgo_cflags: $(go env CGO_CFLAGS)"
    echo "goflags: $(go env GOFLAGS)"
    echo "bench_flags: -count=$COUNT -benchtime=$BENCHTIME -benchmem"
    echo "memprobe_iters: $MEMITERS"
    echo "benchstat: $BENCHSTAT"
    [ -n "${ENVLABEL:-}" ] && echo "environment: $ENVLABEL"
    true
} > "$OUTDIR/metadata.txt"

echo "== metadata written to $OUTDIR/metadata.txt"
cat "$OUTDIR/metadata.txt"

echo "== running equivalence tests (gate)"
CGO_ENABLED=1 go test ./...

echo "== running benchmarks (count=$COUNT, benchtime=$BENCHTIME; this takes a while)"
CGO_ENABLED=1 go test -run '^$' -bench . -benchmem \
    -count "$COUNT" -benchtime "$BENCHTIME" -timeout 180m \
    | tee "$OUTDIR/bench-raw.txt"

echo "== summarizing with benchstat"
go run "$BENCHSTAT" -col /impl "$OUTDIR/bench-raw.txt" > "$OUTDIR/benchstat.txt"
go run "$BENCHSTAT" -format csv -col /impl "$OUTDIR/bench-raw.txt" > "$OUTDIR/benchstat.csv"

echo "== running process-level memory matrix (one process per cell)"
CGO_ENABLED=1 go run ./cmd/memprobe -header > "$OUTDIR/memory.tsv"
CGO_ENABLED=1 go build -o /tmp/uberbench-memprobe ./cmd/memprobe
for workload in $(CGO_ENABLED=1 go run ./cmd/memprobe -list | awk '{print $1}'); do
    for impl in pure uber; do
        /tmp/uberbench-memprobe -impl "$impl" -workload "$workload" -iters "$MEMITERS" \
            >> "$OUTDIR/memory.tsv"
    done
done
rm -f /tmp/uberbench-memprobe

echo "== done; artifacts in $OUTDIR"
