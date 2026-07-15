# Preserved fuzz findings (not auto-executed)

Go fuzz corpus entries that document confirmed findings but are kept **out**
of `testdata/fuzz/` on purpose: entries there are replayed by every ordinary
`go test` run, and the inputs preserved here are pathologically slow, which
would tax every test run for no assertion value. Each entry stays until its
finding is resolved, then moves into a targeted regression test.

## FuzzUpstreamPolygonOperations/bcefd1d03714b6b6

A finite but pathologically slow (~15 s on an Apple M1 Max, go1.26.5)
polygon input discovered by a 30-second fuzz run of
`FuzzUpstreamPolygonOperations` — tracked in
[issue #3](https://github.com/dimchansky/h3-go/issues/3). Because Go's fuzz
engine kills workers that exceed its per-execution budget, this keeps the
target out of the Nightly `-fuzz` rotation (seed-corpus only) until the
pathology is understood. Parity with the H3 C implementation on this input
has **not yet been established** — that comparison is part of the issue's
checklist.

SHA-256 (also the basis of the Go corpus entry name):

```
bcefd1d03714b6b6658c1649fbe49951fa2ebf91a4cc132c067aa7253168b967
```

Replay (from the repository root):

```sh
mkdir -p testdata/fuzz/FuzzUpstreamPolygonOperations
cp testdata/fuzz-findings/FuzzUpstreamPolygonOperations/bcefd1d03714b6b6 \
   testdata/fuzz/FuzzUpstreamPolygonOperations/
CGO_ENABLED=0 go test -run 'FuzzUpstreamPolygonOperations/bcefd1d03714b6b6' -v .
# --- PASS (finite), but expect roughly 15 s of CPU for this single input.
# Clean up afterwards so ordinary test runs stay fast:
rm -r testdata/fuzz
```
